package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
)

type ipdata struct {
	ip  string
	ray string
	ms  int
	mb  float64
}

var DatumReal int64 = 0
var DatumDown float64 = 0
var singChan = make(chan struct{}, 2000)
var ipChan = make(chan string, 500)
var ip2Chan = make(chan ipdata, 20)
var num = 0
var uploadIP = 0
var wg sync.WaitGroup

func main() {
	log.Println("开始")
	err := LoadYaml("config.yaml")
	if err != nil {
		panic(err)
	}
	for {
		now := time.Now()
		_, err = RealPing(C.Host)
		if err != nil {
			continue
		}
		DatumReal = time.Since(now).Milliseconds()
		DatumDown, err = SpeedTestDownload(C.Host)
		if err != nil || DatumDown == 0 {
			continue
		}
		log.Printf("基准为真连接延迟：%dms\t\t下载速度为:%fmb/s\n", DatumReal, DatumDown)
		break
	}
	wg.Add(1)
	go CheckAgain()
	go CheckAgain2()
	time.Sleep(500 * time.Millisecond)
	ip, index, err := GetIPsByCIDRs(C.Cidr)
	go bar(index)
	if err != nil {
		panic(err)
	}
	for {
		if num == index {
			break
		}
		ip.c += 1
		if ip.c == 256 {
			ip.c = 0
			ip.b += 1
			if ip.b == 256 {
				ip.b = 0
				ip.a += 1
			}
		}

		for i := range 256 {
			singChan <- struct{}{}
			ip.d = uint32(i)
			go checktcp(ip.toString())

		}
		num++
	}
	close(singChan)
	for {
		one := len(ipChan)
		time.Sleep(5 * time.Second)
		two := len(ipChan)
		if one == two {
			close(ipChan)
			break
		}
	}
	wg.Wait()
	time.Sleep(time.Second)
	fmt.Println()
	log.Println("结束")
}

func checktcp(ip string) {
	defer func() {
		<-singChan
	}()
	now := time.Now()
	_, err := RealPing(ip)
	if err != nil {
		return
	}
	m := time.Since(now).Milliseconds()
	if m > DatumReal {
		return
	}
	ipChan <- ip
}

func CheckAgain() {

	for ip := range ipChan {
		now := time.Now()
		ray, err := RealPing(ip)
		if err != nil {
			continue
		}
		t := time.Since(now).Milliseconds()
		if t > DatumReal {
			continue
		}
		f1, err := SpeedTestDownload(ip)
		if err != nil {
			continue
		}

		if f1 < DatumDown {
			continue
		}
		ip2Chan <- ipdata{
			ip:  ip,
			ray: ray,
			ms:  int(t),
			mb:  f1,
		}
	}

	close(ip2Chan)
}
func CheckAgain2() {
	for ip := range ip2Chan {
		for {
			c := strconv.Itoa(ip.ms)
			d := strconv.FormatFloat(ip.mb, 'f', -1, 64)
			err := uploadIp(ip.ip, "443", ip.ray, c, d)
			if err == nil {
				uploadIP++
				break
			}
		}
	}
	wg.Done()
}

func bar(index int) {
	for {

		fmt.Printf(
			"\r筛选ip\t\t总长度%d\t\t当前%d\t\t完成：%.2f%%\t\t\t延迟，下载速度测试剩下%d个\t\t\t成功上传IP:%d个",
			index,
			num,
			float64(num)/float64(index)*100,
			len(ipChan),
			uploadIP,
		)
		time.Sleep(time.Second)
	}
}
