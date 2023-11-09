package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

type rpcRes struct {
	value interface{}
	err   error
}

func InitRouter(r *gin.Engine) {
	r.POST("/", updateData)
	r.GET("/:key", findData)
	r.DELETE("/:key", deleteData)
}

func deleteData(c *gin.Context) {
	var key = c.Param("key")
	if _, flag := Data[key]; flag {
		delete(Data, key)
		c.JSON(200, 1)
	} else {
		//本机不存在，询问其他主机
		err := deleteDataFromOthers(key)
		if err != nil {
			c.JSON(200, 0)
			return
		}
		c.JSON(200, 1)
	}
}

// 更新数据
func updateData(c *gin.Context) {
	var jsonData = make(map[string]interface{})
	err := c.BindJSON(&jsonData)
	//fmt.Println(&jsonData)
	if err != nil {
		c.JSON(400, "get json error")
	}
	for key, value := range jsonData {
		//询问其他主机是否存在该数据
		err := updateDataFromOthers(key, value)
		if err != nil {
			//输出错误
			log.Println(err)
			//没有则存在本地
			Data[key] = value
		}
	}
	c.JSON(200, "update success")
}

// 查询数据
func findData(c *gin.Context) {
	var key = c.Param("key")
	if value, flag := Data[key]; flag {
		c.JSON(200, map[string]interface{}{key: value})
	} else if value, err := findDataFromOthers(key); err == nil {
		c.JSON(200, map[string]interface{}{key: value})
	} else {
		c.JSON(404, gin.H{})
	}
}

// 从其他主机获得数据
func findDataFromOthers(key string) (interface{}, error) {
	var res rpcRes
	c1 := make(chan rpcRes)
	c2 := make(chan rpcRes)
	go func() {
		data, err := grpcGetData(conn1, key)
		c1 <- rpcRes{value: data, err: err}
	}()
	go func() {
		data, err := grpcGetData(conn2, key)
		c2 <- rpcRes{value: data, err: err}
	}()
	for i := 0; i < 2; i++ {
		select {
		case res = <-c1:
		case res = <-c2:
		}
		//没错误代表拿到数据，退出
		if res.err == nil {
			break
		}
	}
	return res.value, res.err
}

// 从其他主机删除数据
func deleteDataFromOthers(key string) error {
	c1 := make(chan error)
	c2 := make(chan error)
	var err error
	go func() {
		c1 <- grpcDeleteData(conn1, key)
	}()
	go func() {
		c2 <- grpcDeleteData(conn2, key)
	}()
	for i := 0; i < 2; i++ {
		select {
		case err = <-c1:
		case err = <-c2:
		}
		//没错误代表拿到数据，退出
		if err == nil {
			break
		}
	}
	return err
}

// 从其他主机更新数据
func updateDataFromOthers(key string, value interface{}) error {
	c1 := make(chan error)
	c2 := make(chan error)
	var err error
	go func() {
		c1 <- grpcUpdateData(conn1, key, value)
	}()
	go func() {
		c2 <- grpcUpdateData(conn2, key, value)
	}()
	for i := 0; i < 2; i++ {
		select {
		case err = <-c1:
		case err = <-c2:
		}
		//没错误代表存在，退出
		if err == nil {
			break
		}
	}
	return err
}
