package main

import (
	context "context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"math"
	"net"
	"soil/consul"
	test "soil/proto"
)

type server struct {}

var (
	name = flag.String("name", "", "任务id")

	port = flag.Int("port", 0, "数据集id")
)

func (s server)Predict(c context.Context,req *test.Request) ( *test.Response, error) {
	for _,v := range req.SourcesConfig {
		fmt.Printf("%d  %d  %s\n",v.SourceId,v.ModelType,v.AttrsConfig)
	}
	var resp test.Response
	resp.ErrCode =200
	resp.ErrMsg = "succeed"
	return &resp,nil
}

/*func (s server)Predict2(c context.Context,req *test.Request) ( *test.Response, error) {
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
}*/

func run() error {

	sock,err:= net.Listen("tcp",fmt.Sprintf("%s:%d","10.31.3.111",*port))
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

	test.RegisterPredictServiceServer(s,test_server)

	if err:= s.Serve(sock);err!=nil{
		fmt.Println("---->err:",err)
		return err
	}
	return nil
}



func main () {
	flag.Parse()
	consul.InitConsul()
	go consul.RegisterServer(*name,*name,*port)
	run()
}