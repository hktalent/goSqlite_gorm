package socketserver

import (
	"bufio"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"strconv"
	"strings"
)

type SocketServer struct {
	ln net.Listener
	// The websocket connection.
	conn   *websocket.Conn
	s2wMsg chan map[string]interface{}
	w2sMsg chan map[string]interface{}
	stop   chan bool
}

func (c *SocketServer) Send2IpMsg(msg map[string]interface{}) {

}
func (c *SocketServer) Stop() {
	c.stop <- true
}

func (c *SocketServer) HandlerMsg(s map[string]interface{}) {

	for {
		select {
		case s2w := <-c.s2wMsg:
			if nil != c.conn {
				c.conn.WriteJSON(s2w)
			}
		case w2s := <-c.w2sMsg:
			go c.Send2IpMsg(w2s)
		case <-c.stop:
			break
		}
	}
}
func (ss *SocketServer) GetConn() {
	if nil != ss.ln {
		conn, _ := ss.ln.Accept()
		rmtIp := conn.RemoteAddr().String()
		ss.conn.ReadMessage()
		for {
			// get message, output
			message, err := bufio.NewReader(conn).ReadString('\n')
			if nil == err {
				s1 := string(message)
				s1 = strings.TrimSpace(s1)
				if "" != s1 {
					ss.s2wMsg <- map[string]interface{}{"ip": rmtIp, "msg": string(message)}
					log.Println("Message Received:", string(message))
				}
			}
		}
	}
}
func NewSocketServer(port int) *SocketServer {
	var ssv *SocketServer
	ssv = &SocketServer{}
	ssv.s2wMsg = make(chan map[string]interface{}, 500)
	ln, _ := net.Listen("tcp", ":"+strconv.Itoa(port))
	ssv.ln = ln
	return ssv
}
