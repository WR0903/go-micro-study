# proto

1. message 成员编号， 可以不从1开始, 但是不能重复. -- 不能使用 19000 - 19999 
2. 可以使用 message 嵌套。
3. 定义数组、切片 使用 repeated 关键字
4. 可以使用枚举 enum
5. 可以使用联合体。 oneof 关键字。成员编号，不能重复。

```go
// 默认是 proto2
syntax = "proto3";

// 指定所在包包名
package pb;

// 定义枚举类型
enum Week {
    Monday = 0;   // 枚举值,必须从 0 开始.
    Turesday = 1;
}

// 定义消息体
message Student {
    int32 age = 1;  // 可以不从1开始, 但是不能重复. -- 不能使用 19000 - 19999
    string name = 2;
    People p = 3;
    repeated int32 score = 4;  // 数组
    // 枚举
    Week w = 5;
    // 联合体
    oneof data {
        string teacher = 6;
        string class = 7;
    }
}

// 消息体可以嵌套
message People {
    int32 weight = 1;
}
```

编译 protobuf

```shell
`protoc --go_out=./ *proto`
```

添加rpc服务

```protobuf
// 默认是 proto2
syntax = "proto3";

// 指定所在包包名
package pb;

option go_package="pd/a"; // 这个是指定编译出的go包

message People {
	string name = 1;
}

message Student {
	int32 age = 2;
}

service hello {
	rpc HelloWorld(People) returns (Student);
}
```

编译：

```shell
protoc --go_out=plugins=grpc:./ *proto
```



