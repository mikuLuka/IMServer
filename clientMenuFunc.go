package main

import "fmt"

//1.公聊

func (client *Client) PublicChat() {

	var charMsg string
	fmt.Println("输入内容,exit退出")
	fmt.Scanln(&charMsg)

	for charMsg != "exit" {
		//消息不为空，发给服务端
		if len(charMsg) != 0 {
			sendMsg := charMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn Write err:", err)
				break
			}
		}
		charMsg = ""
		fmt.Println("请输入聊天内容，exit退出")
		fmt.Scanln(&charMsg)
	}
}

//2.私聊

// 查询在线用户
func (client *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write err:", err)
		return
	}

}

func (client *Client) PrivateChat() {
	var remoteName string
	var charMsg string

	client.SelectUsers()
	fmt.Println("请输入聊天对象，exit退出")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println("请输入消息内容，exit退出")
		fmt.Scanln(&charMsg)

		for charMsg != "exit" {
			if len(charMsg) != 0 {

				sendMsg := "to|" + remoteName + "|" + charMsg + "\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn Write err:", err)
					break
				}
			}
		}
	}
}

// 3.更新用户名
func (client *Client) UpdateName() bool {
	fmt.Println("请输入用户名")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}
	return true
}
