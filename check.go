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
)

func afterDash(s string) string {
	if i := strings.IndexByte(s, '-'); i >= 0 {
		return s[i+1:]
	}
	return ""
}
func RealPing(addr string) error {
	conn, err := net.DialTimeout(
		"tcp",
		addr+":80",
		time.Second,
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	req := "GET / HTTP/1.1\r\n" +
		"Host: example.com\r\n" +
		"User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36\r\n" +
		"Accept: */*\r\n" +
		"Connection: close\r\n" +
		"\r\n"

	if _, err = conn.Write([]byte(req)); err != nil {
		return err
	}
	resp, err := http.ReadResponse(bufio.NewReader(conn), nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	return nil
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
		Timeout:   10 * time.Second,
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
		if time.Since(start) >= 5*time.Second {
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

	// MB/s
	speed := float64(total) / elapsed / 1024 / 1024 / 8

	return speed, nil
}
func uploadIp(ip, port string) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			ServerName: C.Host,
		},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}
	ip = ip + ":" + port
	req, err := http.NewRequest(
		"POST",
		"https://"+C.Host+":443/ip",
		bytes.NewBufferString(ip),
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
