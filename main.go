package main

import (
	sdsc "SDSC/grpc"
	"github.com/gin-gonic/gin"
	"os"
)

var Data map[string]interface{} = make(map[string]interface{})
var ServerName, ServerName1, ServerName2 string
var conn1, conn2 sdsc.SDSCClient

func main() {
	go grpcServer()
	getServerName()
	conn1 = connect(ServerName1)
	conn2 = connect(ServerName2)
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
