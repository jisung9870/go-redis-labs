package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
)

var (
	configPath string
	serverPort int
)

func main() {

	flag.StringVar(&configPath, "config", "./config.yaml", "Specify config file path")
	flag.IntVar(&serverPort, "port", 8080, "Specify server port")
	flag.Parse()

	cfg, err := LoadFile(configPath)
	if err != nil {
		panic(err)
	}

	queue, err := RedisConnect(cfg.RedisConfig)
	if err != nil {
		panic(err)
	}

	defer queue.Close()

	mg := NewManager()
	mg.SetQueue(queue)

	ln, err := net.Listen("tcp", ":"+strconv.Itoa(serverPort))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Listening server - 0.0.0.0:%d \n", serverPort)
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("error accepting connection:", err.Error())
			continue
		} else {
			fmt.Printf("Connection Client Access[%s] - %s\n", conn.LocalAddr().String(), conn.RemoteAddr())
		}

		go ReadHandler(conn, mg)
	}
}
