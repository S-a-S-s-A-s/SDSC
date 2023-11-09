package main

import (
	"github.com/gin-gonic/gin"
	"os"
)

var Data map[string]interface{} = make(map[string]interface{})
var ServerName, ServerName1, ServerName2 string

func main() {
	go grpcServer()
	getServerName()
	r := gin.Default()
	InitRouter(r)
	r.Run(":9527")
}

func getServerName() {
	ServerName = os.Getenv("SERVERNAME")
	switch ServerName {
	case "sdsc1":
		ServerName1 = "sdsc2"
		ServerName2 = "sdsc3"
	case "sdsc2":
		ServerName1 = "sdsc1"
		ServerName2 = "sdsc3"
	case "sdsc3":
		ServerName1 = "sdsc1"
		ServerName2 = "sdsc2"
	default:
		ServerName1 = "localhost"
		ServerName2 = "loaclhost"
	}
}
