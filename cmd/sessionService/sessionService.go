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
		log.Fatal("Can`t listen port", err)
	}

	redisConn, err := redis.DialURL(viper.GetString("redisDB"))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Success connect to redis server...")

	server := grpc.NewServer()

	session.RegisterSessionCheckerServer(server, sessionService.NewSessionManager(redisConn))

	server.Serve(lis)
}
