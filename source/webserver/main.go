/**
* @Author: impakho
* @Date: 2020/04/12
* @Github: https://github.com/impakho
*/

package main

import (
	"bytes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/chacha20"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
	"v2link/common"
	"webserver/crypt"
)

const COPYRIGHT_AUTHOR = "impakho"
const COPYRIGHT_DATE = "2020/04/12"
const COPYRIGHT_GITHUB = "https://github.com/impakho"

const DEBUG = false
var REALWORD_ADDR string
var KEY []byte
var NONCE []byte
var proxyListener net.Listener
var ticktock = []byte{164, 163, 4, 185, 30, 241, 150, 198, 10, 38, 77, 233, 175, 253, 177, 255, 6, 238, 229, 207, 107, 46, 12, 2, 23, 106, 151, 183, 149, 172, 184, 17, 26, 143, 19, 131, 229, 175, 103, 201, 106, 38, 153, 43, 28, 173, 63, 65, 223, 170, 54, 54, 8, 162, 4, 157}
var ticking map[string]int64
var tickingLock sync.Mutex

func main() {
	if DEBUG {
		REALWORD_ADDR = "127.0.0.1:4080"
	}else{
		if len(os.Args) < 2 {
			log.Fatal("./webserver [realworld_addr]")
		}
		REALWORD_ADDR = os.Args[1]
	}
	KEY = Sha256([]byte("de1ctf-mc2020"))
	NONCE = Sha256([]byte("de1ta-team"))[:24]
	server := &http.Server{
		Addr: "0.0.0.0:80",
		WriteTimeout: 10 * time.Second,
		ReadTimeout: 10 * time.Second,
		Handler: muxRegister(),
	}
	crypt.Init()
	ticking = make(map[string]int64, 0)
	fmt.Println("webserver started.")
	go StartRealWorldProxy()
	if DEBUG {
		go TestWebProxy()
		go TestRealWorldProxy()
	}
	go MinecraftTickTock()
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func muxRegister() * mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/webproxy", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(405)
			return
		}
		ip, _, err := common.ParseAddr(r.RemoteAddr)
		if err != nil {
			w.Write([]byte("internal error"))
			return
		}
		if !DEBUG {
			tickingLock.Lock()
			if _, ok := ticking[ip]; !ok {
				tickingLock.Unlock()
				w.WriteHeader(403)
				return
			}
			tickingLock.Unlock()
		}
		s, err := ioutil.ReadAll(r.Body)
		if err != nil || len(s) < 10 {
			w.WriteHeader(500)
			return
		}
		cipher, err := chacha20.NewUnauthenticatedCipher(KEY[:], NONCE[:])
		if err != nil {
			w.WriteHeader(500)
			return
		}
		buff := make([]byte, len(s))
		cipher.XORKeyStream(buff, s)
		if !bytes.Contains(buff, []byte("|")) {
			w.WriteHeader(500)
			return
		}
		sep := bytes.Split(buff, []byte("|"))
		ip, port, err := ParseAddr(string(sep[0]))
		if err != nil {
			w.WriteHeader(500)
			return
		}
		if !IsPrivateIP(ip) {
			w.WriteHeader(403)
			return
		}
		rconn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
		hj, _ := w.(http.Hijacker)
		conn, _, _ := hj.Hijack()
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		rconn.SetReadDeadline(time.Now().Add(10 * time.Second))
		rconn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		rconn.Write(sep[1])
		go Relay(conn, rconn)
	})
	r.HandleFunc("/ticktock", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(405)
			return
		}
		params := r.URL.Query()
		text, ok := params["text"]
		if !ok || len(text) <= 0 || len(text[0]) <= 0 {
			w.Write([]byte("invalid param: text"))
			return
		}
		c, err := crypt.NewCipher(KEY[:16])
		if err != nil {
			w.Write([]byte("internal error"))
			return
		}
		s := cipher.NewCFBEncrypter(c, NONCE[:16])
		plain := []byte(text[0])
		buff := make([]byte, len(plain))
		s.XORKeyStream(buff, plain)
		w.Write([]byte(fmt.Sprintf("TickTock: %s", base64.StdEncoding.EncodeToString(buff))))
		if bytes.Equal(buff, ticktock) {
			w.Write([]byte("\nMinecraft is ticking..."))
			go func() {
				ip, _, err := common.ParseAddr(r.RemoteAddr)
				if err != nil {
					return
				}
				tickingLock.Lock()
				if _, ok := ticking[ip]; ok {
					tickingLock.Unlock()
					return
				}
				tickingLock.Unlock()
				go func() {
					f, err := os.OpenFile("log/weblog.txt", os.O_WRONLY|os.O_APPEND, 0222)
					if err == nil {
						f.Write([]byte(fmt.Sprintf("%s %s\n", time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05"), ip)))
						f.Close()
					}
				}()
				tickingLock.Lock()
				ticking[ip] = time.Now().Unix()
				tickingLock.Unlock()
			}()
		}
	})
	r.PathPrefix("/").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(405)
			return
		}
		if r.URL.Path == "/index.html" {
			r.URL.Path = "/"
		}
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "./www/index.html")
			return
		}
		http.FileServer(http.Dir("./www/")).ServeHTTP(w, r)
	}))
	return r
}

func StartRealWorldProxy() {
	var err error
	proxyListener, err = net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("realworld proxy started.")
	AcceptTCPConn()
}

func AcceptTCPConn() {
	for {
		conn, err := proxyListener.Accept()
		if err != nil {
			continue
		}

		go HandleTCPConn(conn)
	}
}

func HandleTCPConn(conn net.Conn) {
	ip, _, err := common.ParseAddr(conn.RemoteAddr().String())
	if err != nil {
		conn.Close()
		return
	}
	if !DEBUG {
		tickingLock.Lock()
		if _, ok := ticking[ip]; !ok {
			tickingLock.Unlock()
			conn.Close()
			return
		}
		tickingLock.Unlock()
	}

	rconn, err := net.Dial("tcp", REALWORD_ADDR)
	if err != nil {
		conn.Close()
		return
	}
	go RelayCrypt(conn, rconn, KEY[:], NONCE[:])
}

func MinecraftTickTock() {
	tickCnt := 0
	for {
		if tickCnt % 24000 == 0 {
			tickingLock.Lock()
			ticking = make(map[string]int64, 0)
			tickingLock.Unlock()
		}
		time.Sleep(50 * time.Millisecond)
		tickCnt++
	}
}