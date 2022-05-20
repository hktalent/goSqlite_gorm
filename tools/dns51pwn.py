#!/usr/bin/env python3
# -*- coding:utf-8 -*-

# DNSRR DNSQR
from scapy.all import *
from datetime import datetime
import urllib3,json,os
import requests as requestsSs
urllib3.disable_warnings()
requestsSs.packages.urllib3.disable_warnings()

requests = requestsSs.session()

# for centos

szCurlDir=os.path.dirname(os.path.abspath(__file__))

# interface = 'eth0'
interface = 'en0'
filter_bpf = 'udp '
g_szDomain = ["51pwn.com",'exploit-poc.com']


g_szEsHost = 'https://www.51pwn.com'
if 'DEBUG' in os.environ and 'True' == os.environ['DEBUG']:
    g_szEsHost = 'http://127.0.0.1:1443'

# @pysnooper.snoop()
def select_DNS(pkt):
    try:
        if DNS not in pkt:
            return
        if DNSRR in pkt and pkt[DNS].qd:
            # pkt.show()
            s=pkt[DNS].qd.qname
            s=s[0:len(s)-1]
            if isinstance(s,bytes):
                s = s.decode('utf8')
            for k in g_szDomain:
                if k in s:
                    return
            aIp = []
            for x in range(pkt[DNS].ancount):
                s1 = pkt[DNSRR][x].rdata
                if isinstance(s1,bytes):
                    s1 = s1.decode('utf8')
                if re.findall(r'([\d]{1,3}\.){3}([\d]{1,3})',s1,re.DOTALL):
                    aIp.append(s1)
            print(s,aIp)
            if 0 < len(aIp) and not re.findall(r'([\d]{1,3}\.){3}([\d]{1,3})',s,re.DOTALL) and '127.0.0.1' not in aIp:
                requests.post(g_szEsHost+'/ip2domain', data=json.dumps({"domain":s,"ip":aIp}), verify=False, headers={"Content-Type": "application/json;charset=UTF-8","User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.2 Safari/605.1.15"},timeout=(10,15),allow_redirects=False)
                szIp = "".join(aIp)
                print('send ok: ' + s + " " + szIp)
    except Exception as e:
        print(e)

sniff(iface=interface, filter=filter_bpf, store=0,  prn=select_DNS)
