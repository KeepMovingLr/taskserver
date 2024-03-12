package server

import (
	"errors"
	"io"
	"log"
	"net"
	"os"
	"ray.li/entrytaskserver/dto"
	"ray.li/entrytaskserver/utils"
	"sync"
	"time"
)

var (
	IllegalMethod = errors.New("this method have not been defined")
	m             sync.RWMutex
)

func RegisterServerListener(port string, readWaitTime int, writeWaitTime int) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", port)
	ExitOnError(err)
	listener, err := net.ListenTCP("tcp4", tcpAddr)
	ExitOnError(err)
	testCount := 0
	for {
		testCount = testCount + 1
		conn, err := listener.AcceptTCP()
		conn.SetKeepAlive(true)
		log.Println("init:", testCount)
		if err != nil {
			log.Println("Failed accepting a connection request:", err)
			continue
		}
		conn.SetReadDeadline(time.Now().Add(time.Duration(readWaitTime) * time.Millisecond))
		go handleClient(conn, readWaitTime, writeWaitTime)
	}
}

func handleClient(conn net.Conn, readWaitTime int, writeWaitTime int) {
	var requestUser dto.UserDTO
	readErr := utils.ReadMsg(conn, &requestUser)
	switch {
	case readErr == io.EOF:
		log.Println("read message failed - this connection has been closed", readErr)
		return
	case readErr != nil:
		log.Println("read message failed", readErr)
		return
	}
	// readDeadline should be reset each time we have read some msg successfully
	conn.SetReadDeadline(time.Now().Add(time.Duration(readWaitTime) * time.Millisecond))
	if handleResult, err := handleBusiness(requestUser); err != nil {
		log.Println("there is an illegal request", requestUser.MethodName, err)
	} else {
		// writeDeadline should be set each time we want to write some msg
		conn.SetWriteDeadline(time.Now().Add(time.Duration(writeWaitTime) * time.Millisecond))
		if writeErr := utils.SendMsg(conn, handleResult); writeErr != nil {
			log.Println("write message failed", writeErr)
		}
	}
}

func handleBusiness(requestUser dto.UserDTO) (handleResult dto.UserDTO, err error) {
	m.Lock()
	method, ok := HandlerMap[requestUser.MethodName]
	m.Unlock()
	if ok {
		return HandlerDispatcher(requestUser, method), nil
	} else {
		return requestUser, IllegalMethod
	}
}

// fatal error, should exit the application
func ExitOnError(err error) {
	if err != nil {
		log.Fatal("Fatal error ", err.Error())
		os.Exit(1)
	}
}
