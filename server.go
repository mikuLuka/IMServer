package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

//1.创建server对象，2.启动server服务，3.处理连接业务handler，协程处理

type Server struct {
	IP   string
	Port int
	//在线用户的列表
	onlineMap map[string]*User
	mapLock   sync.RWMutex

	//消息广播的channel
	Message chan string
}

// 创建一个server 对象,相当于构造函数
func NewServer(ip string, port int) *Server {
	server := &Server{
		IP:        ip,
		Port:      port,
		onlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

//监听message广播消息channel的groutine，有消息就发送给全部在线的user

func (s *Server) ListenMessager() {
	for {
		msg := <-s.Message
		s.mapLock.Lock()
		for _, cli := range s.onlineMap {
			cli.C <- msg
		}
		s.mapLock.Unlock()
	}
}

// 广播消息
func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	s.Message <- sendMsg

}

func (s *Server) Handler(conn net.Conn) {
	//fmt.Println("连接建立成功")
	user := NewUser(conn, s)
	/*
		//用户上线，将用户加入到onlineMap中
		s.mapLock.Lock()
		s.onlineMap[user.Name] = user
		s.mapLock.Unlock()

		//广播当前用户上线
		s.BroadCast(user, "已上线")*/
	//已把上面封装在user的上线方法里
	user.Online()
	//监听用户活跃状态
	isLive := make(chan bool)

	//接收客户端消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}
			//提取用户消息
			msg := string(buf[:n-1])

			//将用户消息广播
			user.DoMessage(msg)

			//用户的状态，活跃
			isLive <- true
		}
	}()

	//当前handler阻塞
	for {
		select {
		case <-isLive:
			//当前用户活跃，重置定时器
		case <-time.After(time.Second * 20):
			//已经超时
			//
			user.SendMsg("超时20s,你已被踢")

			//销毁资源
			close(user.C)
			//关闭连接
			conn.Close()
			//退出handler
			return
		}
	}

}

// 启动服务器的接口
func (s *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("启动服务器失败,err:", err)
		return
	}
	//close socket
	defer listener.Close()

	go s.ListenMessager()

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("监听接收失败,err", err)
			continue
		}
		//do handler
		go s.Handler(conn)
	}
}
