package test

import (
	"log"
	"net"
)

func EchoServer(l net.Listener) chan error {
	errChan := make(chan error, 1)
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				errChan <- err
				return
			}

			data, err := ReadBlock(conn)
			if err != nil {
				log.Printf("test echo listener read error: %v", err)
				errChan <- err
				return
			}
			if err := WriteBlock(conn, data); err != nil {
				log.Printf("test echo listener write error: %v", err)
				errChan <- err
				return
			}

			if err := conn.Close(); err != nil {
				log.Printf("test echo listener close error: %v", err)
				errChan <- err
				return
			}
		}
	}()
	return errChan
}
