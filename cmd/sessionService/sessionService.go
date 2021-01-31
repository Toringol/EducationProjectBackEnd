package main

import (
	"log"
	"net"

	"github.com/Toringol/EducationProjectBackEnd/app/sessionService"
	"github.com/Toringol/EducationProjectBackEnd/app/sessionService/session"
	"github.com/Toringol/EducationProjectBackEnd/config"
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("%s", err.Error())
	}

	lis, err := net.Listen("tcp", viper.GetString("sessionServiceListenPort"))
	if err != nil {
		log.Fatalf("Can`t listen port", err)
	}

	redisConn, err := redis.Dial("tcp", viper.GetString("redisDBPort"))
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()

	session.RegisterSessionCheckerServer(server, sessionService.NewSessionManager(redisConn))

	server.Serve(lis)
}
