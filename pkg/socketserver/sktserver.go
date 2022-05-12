package socketserver

import (
	"bufio"
	"encoding/json"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type MsgData struct {
	Ip  string `json:"ip"`
	Msg string `json:"msg"`
}

type IpConnMap struct {
	sync.RWMutex
	m map[string]*net.Conn
}

func (icm *IpConnMap) Set(k string, v *net.Conn) {
	icm.Lock()
	icm.m[k] = v
	icm.Unlock()
}

type SocketServer struct {
	ln net.Listener
	// The websocket connection.
	conn         *websocket.Conn
	s2wMsg       chan MsgData
	w2sMsg       chan MsgData
	stop         chan bool
	ClienThreads int
	clientMap    IpConnMap
}

func (ss *SocketServer) Stop() {
	close(ss.stop)
}

// check socket conn is closed
// 连接关闭时，返回true，否则返回 false
func (ss *SocketServer) IsClosed4Socket(conn net.Conn) bool {
	_, err := conn.Read(make([]byte, 0))
	// 读到末尾，才证明连接可用
	if err != io.EOF {
		return true
	}
	return false
}

// 轮询处理器
// 1、接收stoo信号
// 2、 读取websocket数据 写到对应 bash反弹ip的流中
func (ss *SocketServer) WsMsg2ReverseShell() {
	for {
		select {
		case <-ss.stop:
			break
		case msg := <-ss.s2wMsg:
			if nil != ss.conn {
				ss.conn.WriteJSON(msg)
			}
		default:
			if nil != ss.conn {
				_, message, err := ss.conn.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					continue
				}
				var msgData MsgData
				err = json.Unmarshal(message, &msgData)
				log.Printf("websocket send to %s bash shellrecv: %s", msgData.Ip, message)
				// 检测连接是否已经关闭
				if nil == err {
					conn1, ok := ss.clientMap.m[msgData.Ip]
					if ok {
						bData := []byte(msgData.Msg)
						var conn net.Conn
						conn = *conn1
						n, err1 := conn.Write(bData)
						if nil != err1 || n != len(bData) {
							//sync.Map
							//ss.clientMap.m[]
							conn.Close()
						}
					}
				}
			}
		}
	}
}

// 开启多个线程等待客户端 reverse shell请求到达
// 然后 基于ip和连接建立起关联关系
func (ss *SocketServer) Socket2Ws() {
	if nil != ss.ln {
		conn, _ := ss.ln.Accept()
		rmtIp := strings.Split(conn.RemoteAddr().String(), ":")[0]
		ss.clientMap.Set(rmtIp, &conn)
		for {
			// get message, output
			message, err := bufio.NewReader(conn).ReadString('\n')
			if nil == err {
				s1 := string(message)
				s1 = strings.TrimSpace(s1)
				if "" != s1 {
					ss.s2wMsg <- MsgData{Ip: rmtIp, Msg: string(message)}
					log.Println("Message Received:", string(message))
				}
			}
		}
	}
}

// 默认允许16个客户端并发
// 开启多个线程等待客户端 reverse shell请求到达
// 并将读区到到数据发送给选择器
func (ss *SocketServer) handler4ReverseSocketConn() {
	xThread := make(chan struct{}, ss.ClienThreads)
	for {
		func() {
			xThread <- struct{}{}
			defer func() {
				<-xThread
			}()
			ss.Socket2Ws()
		}()
	}
}

var upgrader = websocket.Upgrader{}

func (ss *SocketServer) BindRequest(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if nil == err {
		ss.conn = conn
		//defer conn.Close()
	}
}

// https://golangdocs.com/golang-gorilla-websockets
func NewSocketServer(port int) *SocketServer {
	var ssv *SocketServer
	ssv = &SocketServer{}
	ssv.s2wMsg = make(chan MsgData, 500)
	ln, _ := net.Listen("tcp", ":"+strconv.Itoa(port))
	ssv.ln = ln
	ssv.ClienThreads = 16
	go ssv.handler4ReverseSocketConn()
	go ssv.WsMsg2ReverseShell()
	return ssv
}

// 启动server后
// 调用BindRequest bind
func NewWebReverseShellServer() *SocketServer {
	return NewSocketServer(202004)
}
