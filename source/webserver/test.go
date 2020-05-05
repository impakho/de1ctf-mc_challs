package main

import (
	"fmt"
	"golang.org/x/crypto/chacha20"
	"net"
	"time"
)

func TestWebProxy() {
	time.Sleep(1 * time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:80")
	if err != nil {
		return
	}
	cipher, _ := chacha20.NewUnauthenticatedCipher(KEY[:], NONCE[:])
	body := []byte("172.20.10.166:80|GET / HTTP/1.1\r\nHost: 127.0.0.1\r\n\r\n")
	buff := make([]byte, len(body))
	cipher.XORKeyStream(buff, body)
	conn.Write([]byte(fmt.Sprintf("POST /webproxy HTTP/1.1\r\nContent-Length: %d\r\nHost: 127.0.0.1\r\n\r\n%s", len(buff), buff)))
	for {
		buff := make([]byte, 2048)
		n, err := conn.Read(buff)
		if err != nil {
			continue
		}
		buff = buff[:n]
		fmt.Println(string(buff))
	}
}

func TestRealWorldProxy() {
	time.Sleep(1 * time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		return
	}
	cipher, _ := chacha20.NewUnauthenticatedCipher(KEY[:], NONCE[:])
	body := []byte("T,Hello World\n")
	buff := make([]byte, len(body))
	cipher.XORKeyStream(buff, body)
	conn.Write(buff)
	cipher, _ = chacha20.NewUnauthenticatedCipher(KEY[:], NONCE[:])
	for {
		buff := make([]byte, 2048)
		n, err := conn.Read(buff)
		if err != nil {
			continue
		}
		buff = buff[:n]
		s := make([]byte, len(buff))
		cipher.XORKeyStream(s, buff)
		fmt.Println(string(s))
	}
}
