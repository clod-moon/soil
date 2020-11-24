package main

import (
	context "context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"soil/consul"
	test "soil/proto"
	"time"
)

/*func main() {
	conn,err:= grpc.Dial("127.0.0.1:8090",grpc.WithInsecure())
	if err!=nil {
		fmt.Println("----->err:",err)
		return
	}
	defer conn.Close()

	type PredictClient interface {
		Predict(ctx context.Context,req *test.Request,opts ...grpc.CallOption)(resp *test.Response, err error)
	}

	type client struct {
		oc *grpc.ClientConn
	}
}*/

const (
	address     = "localhost:8090"
	defaultName = "world"
)

type CallBackFuncs map[string]callback
type ModelMange []string

func (m *ModelMange)AddRegister(name string) {
	*m = append(*m,name)
}

func (m ModelMange)Register() {
	for _,v := range m {
		server_info := consul.Discover(v)
		conn, err := grpc.Dial(fmt.Sprintf("%s:%d",server_info[0].Addr,server_info[0].Port), grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		c1 := test.NewPredictServiceClient(conn)

		callBackFuncs.Register(v,c1.Predict)
	}
}



func (c *CallBackFuncs) Register(funcname string,callback callback){
	(*c)[funcname] = callback
}

func (c CallBackFuncs)GetCallback(funcname string) (callback,error) {
	function,ok := c[funcname]
	if !ok {
		return nil,errors.New("错误的函数名")
	}
	return function,nil
}

var (
	callBackFuncs =make(CallBackFuncs)
	modelMange ModelMange
)

type callback func(ctx context.Context, in *test.Request, opts ...grpc.CallOption) (*test.Response, error)


func init() {

	consul.InitConsul()

	modelMange.AddRegister("predict1")
	modelMange.AddRegister("predict2")

	modelMange.Register()
}


func main() {
	// Set up a connection to the server.


	// Contact the server and print out its response.

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := test.Request{}
	for i:=0;i<5;i++ {
		sc := &test.SourceConfig{
			SourceId:             int64(i),
			ModelType:            int64(i),
			AttrsConfig:          fmt.Sprintf(`{"属性":%d}`,i),
		}
		req.SourcesConfig = append(req.SourcesConfig,sc)
	}

	for i:=0;i<2;i++ {
		f, err := callBackFuncs.GetCallback(fmt.Sprintf("predict%d",i+1))
		if err != nil {
			log.Fatalf("could not found func: %v", err)
		}
		r,err :=f(ctx, &req)
		if err != nil {
			log.Fatalf("could not prdict: %v", err)
		}
		log.Printf("Greeting: %s", r.String())
	}

	/*	for i:=0;i<3;i++ {
		f ,err:= callBackFuncs.GetCallback(fmt.Sprintf("predict%d",i+1))
		if err!=nil {
			fmt.Println(err)
			return
		}
		callbackfunc :=f.(callback)
		callbackfunc(i)
	}*/
}