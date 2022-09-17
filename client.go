package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

// 创建客户端对象
func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       2,
	}
	//连接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}
	client.conn = conn

	//返回对象
	return client
}

// 处理server回应的消息，直接显示到标准输出
func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn)
}

func (client *Client) menu() bool {
	var flag int

	fmt.Println("1,公聊模式")
	fmt.Println("2,私聊模式")
	fmt.Println("3,更新用户名")
	fmt.Println("0,退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag < 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("请输入合法的数字")
		return false
	}

}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {

		}
		//
		switch client.flag {
		case 1:
			//fmt.Println("公聊模式选择")
			client.PublicChat()
		case 2:
			//fmt.Println("私聊模式选择")
			client.PrivateChat()
		case 3:
			//fmt.Println("更新用户名选择")
			client.UpdateName()
		}
	}

}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器地址")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口")
}

func main() {

	//命令行解析
	flag.Parse()

	client := NewClient("127.0.0.1", 8888)
	if client != nil {
		fmt.Println("连接服务器失败")
		return
	}

	go client.DealResponse()

	fmt.Println("连接服务器成功")
	//启动客户端业务
	//select {}
	client.Run()
}
