package main

import (
	"strconv"
	"strings"
	"sync"
)

var bytePool = sync.Pool{
	New: func() any {
		b := make([]byte, 0, 15)
		return b
	},
}

type ipStr struct {
	a uint32
	b uint32
	c uint32
	d uint32
}

func (ip *ipStr) IPv4ToUint32() uint32 {
	return (ip.a << 24) | (ip.b << 16) | (ip.c << 8) | ip.d
}

func Uint32ToIPv4(ip uint32) string {
	buf := bytePool.Get().([]byte)
	defer bytePool.Put(buf)
	buf = strconv.AppendUint(buf, uint64(ip>>24), 10)
	buf = append(buf, '.')
	buf = strconv.AppendUint(buf, uint64((ip>>16)&0xff), 10)
	buf = append(buf, '.')
	buf = strconv.AppendUint(buf, uint64((ip>>8)&0xff), 10)
	buf = append(buf, '.')
	buf = strconv.AppendUint(buf, uint64(ip&0xff), 10)
	return string(buf)
}

func (ip *ipStr) toString() string {
	ip.IPv4ToUint32()
	return Uint32ToIPv4(ip.IPv4ToUint32())
}

func GetIPsByCIDRs(cidr string) (*ipStr, int) {
	ipstr := strings.Split(cidr, "/")
	i, _ := strconv.Atoi(ipstr[1])
	i = 32 - i
	i = Pow(2, i)
	ii := i / 256
	ipstrs := strings.Split(ipstr[0], ".")
	a1, _ := strconv.Atoi(ipstrs[0])
	b1, _ := strconv.Atoi(ipstrs[1])
	c1, _ := strconv.Atoi(ipstrs[2])
	d1, _ := strconv.Atoi(ipstrs[3])
	iii := ipStr{
		a: uint32(a1),
		b: uint32(b1),
		c: uint32(c1),
		d: uint32(d1),
	}
	return &iii, ii
}

func Pow(base, exp int) int {
	result := 1
	for exp > 0 {
		if exp&1 == 1 {
			result *= base
		}
		base *= base
		exp >>= 1
	}
	return result
}
