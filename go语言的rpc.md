# go语言的rpc

## 服务端：

1. 注册 rpc 服务对象。给对象绑定方法（ 1. 定义类， 2. 绑定类方法 ）

   ```go
   rpc.RegisterName("服务名"，回调对象)
   ```

2. 创建监听器 

   ```go
   listener, err := net.Listen()
   ```

3. 建立连接

   ```go
   conn, err := listener.Accept()
   ```

4. 将连接 绑定 rpc 服务。

   ```go
   rpc.ServeConn(conn)
   ```


## 客户端：

1. 用 rpc 连接服务器。

   ```go
   conn, err := rpc.Dial()
   ```

2. 调用远程函数。

   ```go
   conn.Call("服务名.方法名", 传入参数, 传出参数)
   ```

   

## demo

1. 服务端

   ```go
   package main
   
   import (
   	"fmt"
   	"net"
   	"net/rpc"
   )
   
   // 定义类对象
   type World struct {
   }
   
   // 绑定类方法
   func (this *World) HelloWorld(name string, resp *string) error {
   	*resp = name + " 你好!"
   	return nil
   }
   
   func main() {
   	// 1. 注册RPC服务, 绑定对象方法
   	err := rpc.RegisterName("hello", new(World))
   	if err != nil {
   		fmt.Println("注册 rpc 服务失败!", err)
   		return
   	}
   
   	// 2. 设置监听
   	listener, err := net.Listen("tcp", "127.0.0.1:8800")
   	if err != nil {
   		fmt.Println("net.Listen err:", err)
   		return
   	}
   	defer listener.Close()
   
   	fmt.Println("开始监听 ...")
   	// 3. 建立链接
   	conn, err := listener.Accept()
   	if err != nil {
   		fmt.Println("Accept() err:", err)
   		return
   	}
   	defer conn.Close()
   	fmt.Println("链接成功...")
   
   	// 4. 绑定服务
   	rpc.ServeConn(conn)
   }
   ```

2. 客户端

   ```go
   package main
   
   import (
   	"fmt"
   	"net/rpc"
   )
   
   func main() {
   	// 1. 用 rpc 链接服务器 --Dial()
   	conn, err := rpc.Dial("tcp", "127.0.0.1:8800")
   	if err != nil {
   		fmt.Println("Dial err:", err)
   		return
   	}
   	defer conn.Close()
   
   	// 2. 调用远程函数
   	var reply string // 接受返回值 --- 传出参数
   
   	err = conn.Call("hello.HelloWorld", "李白", &reply)
   	if err != nil {
   		fmt.Println("Call:", err)
   		return
   	}
   
   	fmt.Println(reply)
   }
   
   ```

   