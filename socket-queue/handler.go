package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func ReadHandler(conn net.Conn, mg *Manager) {
	uuid := mg.Add(conn)
	defer mg.Close(uuid)

	for {
		message, err := mg.Read(uuid)
		if err != nil {
			log.Printf("read handler error: %v\n", err)
			return
		}
		log.Printf("received message from server: %s\n", message)

		queue := mg.GetQueue()
		cur := time.Now()
		if err = queue.LPush("socket", message); err != nil {
			log.Printf("failed push to Redis: %v\n", err)
			return
		}
		fmt.Printf("push duration %s\n", time.Since(cur).String())
		log.Printf("success push to Redis\n")
		message2, err := queue.RPop("socket")
		if err != nil {
			log.Printf("failed pop from Redis: %v\n", err)
			return
		}
		log.Printf("success pop from redis %s\n", message2)
	}
}
