package main

import (
	"context"
	"go-todolist/task-srv/handler"
	pb "go-todolist/task-srv/proto/task"
	"go-todolist/task-srv/repository"
	"log"
	"time"

	"github.com/micro/go-micro/v2"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 这里是我内网的mongo地址，请根据你得实际情况配置，推荐使用dockers部署
const MONGO_URI = "mongodb://172.18.0.58:27017"

func main() {
	// 在日志中打印文件路径，便于调试代码
	log.SetFlags(log.Llongfile)

	conn, err := connectMongo(MONGO_URI, time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Disconnect(context.Background())

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.task"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	taskHandler := &handler.TaskHandler{
		TaskRepository: &repository.TaskRepositoryImpl{
			Conn: conn,
		},
	}
	if err := pb.RegisterTaskServiceHandler(service.Server(), taskHandler); err != nil {
		log.Fatal(errors.WithMessage(err, "register server"))
	}

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(errors.WithMessage(err, "run server"))
	}
}

// 连接到MongoDB
func connectMongo(uri string, timeout time.Duration) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, errors.WithMessage(err, "create mongo connection session")
	}
	return client, nil
}
