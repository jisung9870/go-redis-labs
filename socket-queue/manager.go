package main

import (
	"errors"
	"net"
	"time"

	"github.com/twinj/uuid"
)

type Manager struct {
	ConnMap           map[string]net.Conn
	ReadDeadLineTime  time.Duration
	WriteDeadLineTime time.Duration
	BufferSize        int
	RedisQueue        *RedisClient
}

func NewManager() *Manager {
	return &Manager{
		ConnMap:           map[string]net.Conn{},
		ReadDeadLineTime:  3 * time.Second,
		WriteDeadLineTime: 3 * time.Second,
		BufferSize:        256,
	}
}

func (mg *Manager) Add(conn net.Conn) string {
	uuid := uuid.NewV4().String()

	mg.ConnMap[uuid] = conn
	return uuid
}

func (mg *Manager) Remove(uuid string) {
	delete(mg.ConnMap, uuid)
}

func (mg *Manager) Close(uuid string) {
	conn := mg.ConnMap[uuid]
	mg.Remove(uuid)
	conn.Close()
}

func (mg *Manager) SetQueue(client *RedisClient) {
	mg.RedisQueue = client
}

func (mg *Manager) GetQueue() *RedisClient {
	return mg.RedisQueue
}

func (mg *Manager) Read(uuid string) (string, error) {
	conn := mg.ConnMap[uuid]
	conn.SetReadDeadline(time.Now().Add(mg.ReadDeadLineTime))

	buffer := make([]byte, mg.BufferSize)
	length, err := conn.Read(buffer)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			err = errors.New("no data ready to be read")
		}
		return "", err
	}
	message := buffer[:length]
	return string(message), nil
}

func (mg *Manager) Write(uuid string, message string) error {
	conn := mg.ConnMap[uuid]
	conn.SetWriteDeadline(time.Now().Add(mg.WriteDeadLineTime))

	_, err := conn.Write([]byte(message))
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			err = errors.New("write timeover")
		}
		return err
	}
	return nil
}
