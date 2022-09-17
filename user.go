package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

// 创建用户接口

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,

		server: server,
	}
	//启动监听的协程
	go user.ListenMessage()
	return user
}

// 1.上线
func (u *User) Online() {
	//用户上线，将用户加入到onlineMap中
	u.server.mapLock.Lock()
	u.server.onlineMap[u.Name] = u
	u.server.mapLock.Unlock()

	//广播当前用户上线
	u.server.BroadCast(u, "已上线")
}

// 2.下线
func (u *User) Offline() {

	u.server.mapLock.Lock()
	delete(u.server.onlineMap, u.Name)
	u.server.mapLock.Unlock()

	//广播当前用户下线
	u.server.BroadCast(u, "已下线")
}

// 指定的用户发消息
func (u *User) SendMsg(msg string) {
	u.conn.Write([]byte(msg))

}

// 3.用户处理消息业务
// 1.who查询用户2.重命名
func (u *User) DoMessage(msg string) {
	if msg == "who" {
		u.server.mapLock.Lock()
		for _, user := range u.server.onlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "在线上。。。\n"
			u.SendMsg(onlineMsg)
		}
		u.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename" {
		newName := strings.Split(msg, "|")[1]
		//判断是否存在
		_, ok := u.server.onlineMap[newName]
		if ok {
			u.SendMsg("用户名已被使用\n")
		} else {
			u.server.mapLock.Lock()
			delete(u.server.onlineMap, u.Name)
			u.server.onlineMap[newName] = u

			u.server.mapLock.Unlock()

			u.Name = newName
			u.SendMsg("已更新用户名：" + u.Name + "\n")
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		//消息格式： to|
		//1.获取用户名
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			u.SendMsg("消息格式不正确，请使用 \"to|张三|消息\"格式.\n")
			return
		}
		//2.根据用户名。得到User对象
		remoteUser, ok := u.server.onlineMap[remoteName]
		if !ok {
			u.SendMsg("该用户名不存在\n")
		}
		//3.获取消息内容
		content := strings.Split(msg, "|")[2]
		if content == "" {
			u.SendMsg("无消息内容，请重发\n")
			return
		}
		remoteUser.SendMsg(u.Name + "对你说：" + content + "\n")

	} else {
		u.server.BroadCast(u, msg)
	}
}

// 监听当前user Channel的方法，一旦有消息，就直接发送给对端客户端
func (u *User) ListenMessage() {
	for {
		msg := <-u.C
		u.conn.Write([]byte(msg + "\n"))
	}

}
