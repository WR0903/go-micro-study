# 用consul做grpc服务发现

1、安装 consul 源码包：

```shell
$ go get -u -v github.com/hashicorp/consul
```

2、proto编写

```protobuf
syntax = "proto3";

package pb;

option go_package="pb/hello";

message Person {
    string name = 1;
    int32 age = 2;
}

// 添加 rpc服务
service hello {
    rpc sayHello (Person) returns (Person);
}
```

然后执行`protoc --go_out=plugins=grpc:./ *.proto`

3、启动consul

```shell
consul  agent --dev
```

4、服务端编写

```go
package main

import (
	"context"
	"fmt"
	"net"

	hello "gorpcConsul/pb/hello"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
)

// 定义类
type Children struct {
}

// 绑定类方法, 实现借口
func (this *Children) SayHello(ctx context.Context, p *hello.Person) (*hello.Person, error) {
	p.Name = "hello  " + p.Name
	return p, nil
}

func main() {
	// 把grpc服务,注册到consul上.
	// 1. 初始化 consul 配置
	consulConfig := api.DefaultConfig()

	// 2. 创建 consul 对象
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		fmt.Println("api.NewClient err:", err)
		return
	}
	// 3. 告诉consul, 即将注册的服务的配置信息
	reg := api.AgentServiceRegistration{
		ID:      "study",
		Tags:    []string{"grcp", "consul"},
		Name:    "grpc And Consul",
		Address: "127.0.0.1",
		Port:    8800,
		Check: &api.AgentServiceCheck{
			CheckID:  "consul grpc test",
			TCP:      "127.0.0.1:8800",
			Timeout:  "1s",
			Interval: "5s",
		},
	}

	// 4. 注册 grpc 服务到 consul 上
	consulClient.Agent().ServiceRegister(&reg)

	//////////////////////以下为 grpc 服务远程调用////////////////////////

	// 1.初始化 grpc 对象,
	grpcServer := grpc.NewServer()

	// 2.注册服务
	hello.RegisterHelloServer(grpcServer, new(Children))

	// 3.设置监听, 指定 IP/port
	listener, err := net.Listen("tcp", "127.0.0.1:8800")
	if err != nil {
		fmt.Println("Listen err:", err)
		return
	}
	defer listener.Close()

	fmt.Println("服务启动... ")

	// 4. 启动服务
	grpcServer.Serve(listener)

}

```

5、客户端编写

```go
package main

import (
	"context"
	"fmt"
	hello "gorpcConsul/pb/hello"
	"strconv"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
)

func main() {
	// 初始化 consul 配置
	consulConfig := api.DefaultConfig()

	// 创建consul对象 -- (可以重新指定 consul 属性: IP/Port , 也可以使用默认)
	consulClient, err := api.NewClient(consulConfig)

	// 服务发现. 从consuL上, 获取健康的服务
	services, _, err := consulClient.Health().Service("grpc And Consul", "grcp", true, nil)

	// 简单的负载均衡.

	addr := services[0].Service.Address + ":" + strconv.Itoa(services[0].Service.Port)

	//////////////////////以下为 grpc 服务远程调用///////////////////////////
	// 1. 链接服务
	//grpcConn, _ := grpc.Dial("127.0.0.1:8800", grpc.WithInsecure())

	// 使用 服务发现consul 上的 IP/port 来与服务建立链接
	grpcConn, _ := grpc.Dial(addr, grpc.WithInsecure())

	// 2. 初始化 grpc 客户端
	grpcClient := hello.NewHelloClient(grpcConn)

	var person hello.Person
	person.Name = "Andy"
	person.Age = 18

	// 3. 调用远程函数
	p, err := grpcClient.SayHello(context.TODO(), &person)

	fmt.Println(p, err)
}

```

然后就有一个简单的验证