package main

import (
	"errors"
	"golang.org/x/crypto/chacha20"
	"io"
	"net"
)

var ErrShortWrite = errors.New("short write")

var EOF = errors.New("EOF")

type LimitedReader struct {
	R io.Reader // underlying reader
	N int64  // max bytes remaining
}

func (l *LimitedReader) Read(p []byte) (n int, err error) {
	if l.N <= 0 {
		return 0, EOF
	}
	if int64(len(p)) > l.N {
		p = p[0:l.N]
	}
	n, err = l.R.Read(p)
	l.N -= int64(n)
	return
}

func CopyCrypt(dst io.Writer, src io.Reader, key []byte, nonce []byte) (written int64, err error) {
	return copyBufferCrypt(dst, src, nil, key, nonce)
}

func copyBufferCrypt(dst io.Writer, src io.Reader, buf []byte, key []byte, nonce []byte) (written int64, err error) {
	if buf == nil {
		size := 32 * 1024
		if l, ok := src.(*LimitedReader); ok && int64(size) > l.N {
			if l.N < 1 {
				size = 1
			} else {
				size = int(l.N)
			}
		}
		buf = make([]byte, size)
	}
	cipher, err := chacha20.NewUnauthenticatedCipher(key, nonce[:])
	if err != nil {
		return written, err
	}
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			cipher.XORKeyStream(buf[0:nr], buf[0:nr])
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != EOF {
				err = er
			}
			break
		}
	}
	return written, err
}

func RelayCrypt(lhs net.Conn, rhs net.Conn, key []byte, nonce []byte) {
	type direction byte

	const (
		dirUplink direction = iota
		dirDownlink
	)

	type duplexConn interface {
		net.Conn
		CloseRead() error
		CloseWrite() error
	}

	upCh := make(chan struct{})

	cls := func(dir direction, interrupt bool) {
		lhsDConn, lhsOk := lhs.(duplexConn)
		rhsDConn, rhsOk := rhs.(duplexConn)
		if !interrupt && lhsOk && rhsOk {
			switch dir {
			case dirUplink:
				lhsDConn.CloseRead()
				rhsDConn.CloseWrite()
			case dirDownlink:
				lhsDConn.CloseWrite()
				rhsDConn.CloseRead()
			default:
				panic("unexpected direction")
			}
		} else {
			lhs.Close()
			rhs.Close()
		}
	}

	// Uplink
	go func() {
		var err error
		_, err = CopyCrypt(rhs, lhs, key, nonce)
		if err != nil {
			cls(dirUplink, true) // interrupt the conn if the error is not nil (not EOF)
		} else {
			cls(dirUplink, false) // half close uplink direction of the TCP conn if possible
		}
		upCh <- struct{}{}
	}()

	// Downlink
	var err error
	_, err = CopyCrypt(lhs, rhs, key, nonce)
	if err != nil {
		cls(dirDownlink, true)
	} else {
		cls(dirDownlink, false)
	}

	<-upCh // Wait for uplink done.
}

func Relay(lhs net.Conn, rhs net.Conn) {
	type direction byte

	const (
		dirUplink direction = iota
		dirDownlink
	)

	type duplexConn interface {
		net.Conn
		CloseRead() error
		CloseWrite() error
	}

	upCh := make(chan struct{})

	cls := func(dir direction, interrupt bool) {
		lhsDConn, lhsOk := lhs.(duplexConn)
		rhsDConn, rhsOk := rhs.(duplexConn)
		if !interrupt && lhsOk && rhsOk {
			switch dir {
			case dirUplink:
				lhsDConn.CloseRead()
				rhsDConn.CloseWrite()
			case dirDownlink:
				lhsDConn.CloseWrite()
				rhsDConn.CloseRead()
			default:
				panic("unexpected direction")
			}
		} else {
			lhs.Close()
			rhs.Close()
		}
	}

	// Uplink
	go func() {
		var err error
		_, err = io.Copy(rhs, lhs)
		if err != nil {
			cls(dirUplink, true) // interrupt the conn if the error is not nil (not EOF)
		} else {
			cls(dirUplink, false) // half close uplink direction of the TCP conn if possible
		}
		upCh <- struct{}{}
	}()

	// Downlink
	var err error
	_, err = io.Copy(lhs, rhs)
	if err != nil {
		cls(dirDownlink, true)
	} else {
		cls(dirDownlink, false)
	}

	<-upCh // Wait for uplink done.
}