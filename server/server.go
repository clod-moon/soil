package main

import (
	context "context"
	"fmt"
	"google.golang.org/grpc"
	"math"
	"net"
	test "soil/proto"
)

type server struct {}


func (s server)Predict1(c context.Context,req *test.Request) ( *test.Response, error) {
	for _,v := range req.SourcesConfig {
		fmt.Printf("%d  %d  %s\n",v.SourceId,v.ModelType,v.AttrsConfig)
	}
	var resp test.Response
	resp.ErrCode =200
	resp.ErrMsg = "succeed"
	return &resp,nil
}

func (s server)Predict2(c context.Context,req *test.Request) ( *test.Response, error) {
	for _,v := range req.SourcesConfig {
		fmt.Printf("%d  %d  %s\n",v.SourceId,v.ModelType,v.AttrsConfig)
	}
	var resp test.Response
	resp.ErrCode =200
	resp.ErrMsg = "succeed"
	return &resp,nil
}

func (s server)Predict3(c context.Context,req *test.Request) ( *test.Response, error) {
	for _,v := range req.SourcesConfig {
		fmt.Printf("%d  %d  %s\n",v.SourceId,v.ModelType,v.AttrsConfig)
	}
	var resp test.Response
	resp.ErrCode =200
	resp.ErrMsg = "succeed"
	return &resp,nil
}

func run() error {

	sock,err:= net.Listen("tcp","127.0.0.1:8090")
	if err!=nil {
		fmt.Println("---->err:",err)
		return err
	}

	var options =[]grpc.ServerOption{
		grpc.MaxRecvMsgSize(math.MaxInt32),
		grpc.MaxSendMsgSize(math.MaxInt32),
	}

	s := grpc.NewServer(options...)

	test_server := &server{}

	test.RegisterPredict1ServiceServer(s,test_server)
	test.RegisterPredict2ServiceServer(s,test_server)
	test.RegisterPredict3ServiceServer(s,test_server)

	if err:= s.Serve(sock);err!=nil{
		fmt.Println("---->err:",err)
		return err
	}
	return nil
}


func main () {
	run()
}