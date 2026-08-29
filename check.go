package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func RealPing(addr string) (string, error) {
	conn, err := net.DialTimeout(
		"tcp",
		addr+":80",
		500*time.Millisecond,
	)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	req := "GET / HTTP/1.1\r\n" +
		"Host: example.com\r\n" +
		"User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36\r\n" +
		"Accept: */*\r\n" +
		"Connection: close\r\n" +
		"\r\n"

	if _, err = conn.Write([]byte(req)); err != nil {
		return "", err
	}
	resp, err := http.ReadResponse(bufio.NewReader(conn), nil)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errors.New(resp.Status)
	}
	h := resp.Header.Get("cf-ray")
	hs := strings.Split(h, "-")
	if len(hs) < 2 {
		return "UN", nil
	}
	return hs[1], nil
}

func ping(ip string) int64 {
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return 1000
	}
	defer conn.Close()

	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   1,
			Seq:  1,
			Data: []byte("ping"),
		},
	}

	data, _ := msg.Marshal(nil)

	start := time.Now()

	_, err = conn.WriteTo(data, &net.IPAddr{
		IP: net.ParseIP(ip),
	})
	if err != nil {
		return 1000
	}

	conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))

	buf := make([]byte, 1500)
	_, _, err = conn.ReadFrom(buf)
	if err != nil {
		return 1000
	}

	return time.Since(start).Milliseconds()
}

func SpeedTestDownload(addr string) (float64, error) {
	url := "http://" + addr + ":80/download"

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			ServerName: C.Host,
		},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   5 * time.Second,
	}

	req, err := http.NewRequest(
		"GET",
		url,
		nil,
	)
	if err != nil {
		return 0, err
	}

	// HTTP Host
	req.Host = C.Host

	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Cookie", "auth="+C.Auth)
	req.Header.Set("Connection", "close")

	start := time.Now()

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, errors.New(resp.Status)
	}

	buf := make([]byte, 1024*1024)

	var total int64

	for {
		if time.Since(start) >= 8*time.Second {
			break
		}

		n, err := resp.Body.Read(buf)

		if n > 0 {
			total += int64(n)
		}

		if err != nil {
			break
		}
	}

	elapsed := time.Since(start).Seconds()

	if elapsed <= 0 {
		return 0, nil
	}

	speed := float64(total) / elapsed / 1024 / 1024 / 8

	return speed, nil
}
func uploadIp(ip, port, ray, ms string) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			ServerName: C.Host,
		},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}
	b := bytes.NewBufferString(ip)
	b.WriteString(":")
	b.WriteString(port)
	b.WriteString("#")
	b.WriteString(ray)
	b.WriteString("--")
	b.WriteString(ms)

	req, err := http.NewRequest(
		"POST",
		"https://"+C.Host+":443/ip",
		b,
	)
	if err != nil {
		return err
	}
	req.Host = C.Host
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Cookie", "auth="+C.Auth)
	req.Header.Set("Connection", "close")

	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}
