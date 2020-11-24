package consul

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	consulapi "github.com/hashicorp/consul/api"
)

var (
	count int64
	client *consulapi.Client
)

type ServerInfo struct {
	Addr string
	Port int
}

// consul 服务端会自己发送请求，来进行健康检查
func consulCheck(w http.ResponseWriter, r *http.Request) {

	s := "consulCheck" + fmt.Sprint(count) + "remote:" + r.RemoteAddr + " " + r.URL.String()
	fmt.Println(s)
	fmt.Fprintln(w, s)
	count++
}

func RegisterServer(id,name string,prot int) {
	fmt.Println("---->name:",name)
	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = id      // 服务节点的名称
	registration.Name = name      // 服务名称
	registration.Port = prot              // 服务端口
	registration.Tags = []string{} // tag，可以为空
	registration.Address = "10.31.3.111"      // 服务 IP

	checkPort := prot+1
	registration.Check = &consulapi.AgentServiceCheck{ // 健康检查
		HTTP:                           fmt.Sprintf("http://%s:%d%s", registration.Address, checkPort, "/check"),
		Timeout:                        "3s",
		Interval:                       "5s",  // 健康检查间隔
		DeregisterCriticalServiceAfter: "30s", //check失败后30秒删除本服务，注销时间，相当于过期时间
		// GRPC:     fmt.Sprintf("%v:%v/%v", IP, r.Port, r.Service),// grpc 支持，执行健康检查的地址，service 会传到 Health.Check 函数中
	}

	err := client.Agent().ServiceRegister(registration)
	if err != nil {
		log.Fatal("register server error : ", err)
	}

	http.HandleFunc("/check", consulCheck)
	http.ListenAndServe(fmt.Sprintf(":%d", checkPort), nil)
}

func Discover(name string) []ServerInfo{
	var lastIndex uint64
	services, metainfo, err := client.Health().Service(name, "", true, &consulapi.QueryOptions{
		WaitIndex: lastIndex, // 同步点，这个调用将一直阻塞，直到有新的更新
	})
	if err != nil {
		logrus.Warn("error retrieving instances from Consul: %v", err)
	}
	lastIndex = metainfo.LastIndex

	var servers []ServerInfo
	//addrs := map[string]struct{}{}
	for _, service := range services {
		servers = append(servers,ServerInfo{
			Addr: service.Service.Address,
			Port: service.Service.Port,
		})
		//fmt.Println("service.Service.Address:", service.Service.Address, "service.Service.Port:", service.Service.Port)
		//addrs[net.JoinHostPort(service.Service.Address, strconv.Itoa(service.Service.Port))] = struct{}{}

	}
	return servers
}

func InitConsul() {
	config := consulapi.DefaultConfig()
	config.Address = "10.17.1.126:8500"
	cli, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatal("consul client error : ", err)
	}
	client = cli
}

