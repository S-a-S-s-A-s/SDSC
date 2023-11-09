package main

import (
	pb "SDSC/grpc"
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
)

type RpcServer struct {
	*pb.UnimplementedSDSCServer
}
type MyData struct {
	Value interface{}
}

func (s *RpcServer) GetData(c context.Context, req *pb.Req) (*pb.Res, error) {
	if _, flag := Data[req.Key]; flag {
		myData := &MyData{
			Value: Data[req.Key],
		}
		buf := new(bytes.Buffer)
		//gob编码
		enc := gob.NewEncoder(buf)
		if err := enc.Encode(myData); err != nil {
			log.Println(err)
		}
		myDataMessage := &pb.Res{
			Value: &anypb.Any{
				TypeUrl: "type.googleapis.com/server.MyData",
				Value:   buf.Bytes(), // 将 MyData 序列化为字节流
			},
		}
		return myDataMessage, nil
	} else {
		return nil, errors.New("not found!")
	}
}

func (s *RpcServer) UpdateData(c context.Context, req *pb.ReqUpdate) (*emptypb.Empty, error) {
	log.Println(req.Key)
	// 是否存在数据
	if _, flag := Data[req.Key]; flag {
		datas := req.Value.Value
		buf := new(bytes.Buffer)
		buf.Write(datas)
		dec := gob.NewDecoder(buf)
		var info *MyData
		if err := dec.Decode(&info); err != nil {
			log.Println(err)
		}
		//存在数据则更新
		Data[req.Key] = info.Value
		return &emptypb.Empty{}, nil
	} else {
		return &emptypb.Empty{}, errors.New("no data")
	}
}

func (s *RpcServer) DeleteData(c context.Context, req *pb.Req) (*emptypb.Empty, error) {
	if _, flag := Data[req.Key]; flag {
		delete(Data, req.Key)
		return &emptypb.Empty{}, nil
	} else {
		return &emptypb.Empty{}, errors.New("no data")
	}
}

func grpcServer() {
	grpcSrv := grpc.NewServer()
	pb.RegisterSDSCServer(grpcSrv, new(RpcServer))

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("net Listen error", err)
	}

	grpcSrv.Serve(listener)
}
