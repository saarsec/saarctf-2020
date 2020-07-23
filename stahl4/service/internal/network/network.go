package network

import (
	"io"
	"stahl4/service/internal/file"
	"errors"
	"log"
	"net"
	"os"
	"time"
)

func init() {
	f, err := file.Open("server.log")
	if err != nil {
		panic(err)
	}
	go limitLogSize(f)
	multi := io.MultiWriter(f, os.Stdout)
	log.SetOutput(multi)
}

func limitLogSize(f *os.File) {
	for {
		time.Sleep(30 * time.Second)
		file.CheckFileSize(f, 1000, 50)
	}
}

func writeConnectionTimeout(timeout time.Duration, connection net.Conn, message []byte) (int, error) {
	// Caller has to lock the connection owner
	if connection == nil {
		return 0, errors.New("no connection to write to")
	}
	_ = connection.SetWriteDeadline(time.Now().Add(timeout))
	n, err := connection.Write(message)
	return n, err
}

func readConnectionTimeout(timeout time.Duration, connection net.Conn, buffer []byte) (int, error) {
	// Caller has to lock the connection owner
	if connection == nil {
		return 0, errors.New("no connection to read from")
	}
	_ = connection.SetReadDeadline(time.Now().Add(timeout))
	n, err := connection.Read(buffer)
	return n, err
}
