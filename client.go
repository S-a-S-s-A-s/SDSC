package main

import (
	sdsc "SDSC/grpc"
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"
	"log"
)

func connect(serverName string) *grpc.ClientConn {
	conn, err := grpc.Dial(serverName+":8080", grpc.WithInsecure())
	if err != nil {
		log.Println("grpc Dial error", err)
		return nil
	} else {
		return conn
	}
}

func grpcGetData(serverName string, key string) (interface{}, error) {
	conn := connect(serverName)
	if conn == nil {
		return nil, errors.New("grpc Dial error")
	}
	defer conn.Close()
	client := sdsc.NewSDSCClient(conn)
	rep, err := client.GetData(context.TODO(), &sdsc.Req{Key: key})
	//log.Println(rep)
	if err == nil {
		datas := rep.Value.Value
		buf := new(bytes.Buffer)
		buf.Write(datas)
		dec := gob.NewDecoder(buf)
		var info *MyData
		if err := dec.Decode(&info); err != nil {
			log.Println(err)
		}
		return info.Value, nil
	} else {
		//fmt.Println(err)
		return nil, err
	}
}

func grpcDeleteData(serverName string, key string) error {
	conn := connect(serverName)
	if conn == nil {
		return errors.New("grpc Dial error")
	}
	defer conn.Close()
	client := sdsc.NewSDSCClient(conn)
	_, err := client.DeleteData(context.TODO(), &sdsc.Req{Key: key})
	return err
}

func grpcUpdateData(serverName string, key string, value interface{}) error {
	myData := &MyData{
		Value: value,
	}
	buf := new(bytes.Buffer)
	//gob编码
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(myData); err != nil {
		log.Println(err)
	}
	conn := connect(serverName)
	if conn == nil {
		return errors.New("grpc Dial error")
	}
	defer conn.Close()
	client := sdsc.NewSDSCClient(conn)
	_, err := client.UpdateData(context.TODO(), &sdsc.ReqUpdate{Key: key, Value: &anypb.Any{
		TypeUrl: "type.googleapis.com/server.MyData",
		Value:   buf.Bytes(), // 将 MyData 序列化为字节流
	}})
	return err
}
