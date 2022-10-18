package main

import (
	"bufio"
	crypto_rand "crypto/rand"
	"encoding/binary"
	"flag"
	"math/rand"
	"net"
	"os"
)

var (
	random *rand.Rand
	writer *bufio.Writer
	count  int64
)

func init() {
	flag.Int64Var(&count, "c", 0, "count of IPs to generate (0 = infinite)")
	writer = bufio.NewWriterSize(os.Stdout, 4096)
}

func init() {
	b := make([]byte, 8)
	if _, err := crypto_rand.Read(b); err != nil {
		panic("Cryptorandom seed failed: " + err.Error())
	}
	seed := int64(binary.LittleEndian.Uint64(b))
	random = rand.New(rand.NewSource(seed))
}

func notGlobal(intip uint32) bool {
	return (0xe0000000 <= intip && intip <= 0xefffffff) || // 224.0.0.0 - 239.255.255.255
		(0xf0000000 <= intip && intip <= 0xfffffffe) || // 240.0.0.0 - 255.255.255.254
		(0x0A000000 <= intip && intip <= 0x0AFFFFFF) || // 10.0.0.0 - 10.255.255.255
		(0x7F000000 <= intip && intip <= 0x7FFFFFFF) || // 127.0.0.0 - 127.255.255.255
		(0x64400000 <= intip && intip <= 0x647FFFFF) || // 100.64.0.0 - 100.127.255.255
		(0xAC100000 <= intip && intip <= 0xAC1FFFFF) || // 172.16.0.0 - 172.31.255.255
		(0xC6120000 <= intip && intip <= 0xC613FFFF) || // 198.18.0.0 - 198.19.255.255
		(0xA9FE0000 <= intip && intip <= 0xA9FEFFFF) || // 169.254.0.0 - 169.254.255.255
		(0xC0A80000 <= intip && intip <= 0xC0A8FFFF) || // 192.168.0.0 - 192.168.255.255
		(0xC0000000 <= intip && intip <= 0xC00000FF) || // 192.0.0.0 - 192.0.0.255
		(0xC0000200 <= intip && intip <= 0xC00002FF) || // 192.0.2.0 - 192.0.2.255
		(0xc0586300 <= intip && intip <= 0xc05863ff) || // 192.88.99.0 - 192.88.99.255
		(0xC6336400 <= intip && intip <= 0xC63364FF) || // 198.51.100.0 - 198.51.100.255
		(0xCB007100 <= intip && intip <= 0xCB0071FF) || // 203.0.113.0 - 203.0.113.255
		(0xe9fc0000 <= intip && intip <= 0xe9fc00ff) // 233.252.0.0 - 233.252.0.255
}

func GenerateIP() net.IP {
	var intip uint32
	for {
		intip = 0x01000000 + random.Uint32()%0xfeffffff
		if notGlobal(intip) {
			continue
		}
		return Uint32ToIP(intip)
	}
}

func Uint32ToIP(intip uint32) net.IP {
	ip := make(net.IP, net.IPv4len)
	binary.BigEndian.PutUint32(ip, intip)
	return ip
}

func PrintIP() {
	ip := GenerateIP()
	writer.WriteString(ip.String())
	writer.WriteByte(10)
}

func main() {
	var i int64

	flag.Parse()

	if count > 0 {
		for i = 0; i < count; i++ {
			PrintIP()
		}
		writer.Flush()
		return
	}

	for {
		PrintIP()
	}
}
