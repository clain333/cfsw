package main

import (
	"strconv"
	"strings"
)

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
	buf := make([]byte, 0, 15)

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

func Pow(base, exp int) int {
	result := 1

	for i := 0; i < exp; i++ {
		result *= base
	}

	return result
}
func GetIPsByCIDRs(cidr string) (*ipStr, int, error) {
	c := strings.Split(cidr, "/")
	start := c[0]
	indexStr := c[1]
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return nil, 0, err
	}
	ipNum := Pow(2, 32-index)
	ipNum = ipNum / 256
	ipstr := strings.Split(start, ".")
	a1, err := strconv.Atoi(ipstr[0])
	b1, err := strconv.Atoi(ipstr[1])
	c1, err := strconv.Atoi(ipstr[2])
	d1, err := strconv.Atoi(ipstr[3])
	i := ipStr{
		a: uint32(a1),
		b: uint32(b1),
		c: uint32(c1),
		d: uint32(d1),
	}
	return &i, ipNum, nil
}
