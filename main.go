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

	service_list := strings.Split(service, ",")
	for _, item := range service_list {
		if item == "grpc" {
			GRPCServe(":50051")
		}
	}
	for _, item := range service_list {
		if item == "http" {
			HTTPServe(":8421")
		}
	}
}
