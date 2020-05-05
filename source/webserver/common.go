package main

import (
	"crypto/sha256"
	"errors"
	"net"
	"strconv"
	"strings"
)

func Sha256(bytes []byte) []byte {
	sum := sha256.Sum256(bytes)
	return sum[:]
}

func ParseAddr(addr string) (ip string, port uint16, err error) {
	sep := strings.Split(addr, ":")
	tip := net.ParseIP(sep[0])
	tport, err := strconv.Atoi(sep[1])
	if tip == nil || err != nil || tport < 1 || tport > 65535 {
		return ip, port, errors.New("invalid addr format")
	}
	return tip.String(), uint16(tport), nil
}

func IsPrivateIP(ip string) bool {
	ipNet := net.ParseIP(ip)
	if ipNet == nil {
		return false
	}
	if ip4 := ipNet.To4(); ip4 != nil {
		switch {
		case ip4[0] == 10:
			return true
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return true
		case ip4[0] == 192 && ip4[1] == 168:
			return true
		default:
			return false
		}
	}
	return false
}