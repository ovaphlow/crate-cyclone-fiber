package main

import (
	"log"
	"os"
	"ovaphlow/cratecyclone/utility"
	"runtime"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	utility.InitSlog()

	err := godotenv.Load()
	if err != nil {
		utility.Slogger.Error("加载环境变量失败")
		log.Fatal(err.Error())
	}
	service := os.Getenv("SERVICE")
	http_port := os.Getenv("HTTP_PORT")
	grpc_port := os.Getenv("GRPC_PORT")

	service_list := strings.Split(service, ",")
	for _, item := range service_list {
		if item == "grpc" {
			GRPCServe(grpc_port)
		}
	}
	for _, item := range service_list {
		if item == "http" {
			HTTPServe(http_port)
		}
	}
}
