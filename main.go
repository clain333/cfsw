package main

import (
	"container/heap"
	"fmt"
	"log"
	"strconv"
	"time"
)

var singChan = make(chan struct{})
var heapChan = make(chan Item)
var uploadIP = 0

func main() {
	log.Println("开始")
	err := LoadYaml("config.yaml")
	if err != nil {
		panic(err)
	}
	singChan = make(chan struct{}, C.TNum)
	heapChan = make(chan Item, C.SpeedNum)
	h := &MaxHeap{}
	*h = append(*h, &Item{
		IP: "0.0.0.0",
		Ms: 1000,
	})
	heap.Init(h)
	ip, index := GetIPsByCIDRs(C.Cidr)
	num := 0
	go inChan(h)
	for {
		if index == num {
			break
		}
		if ip.c == 256 {
			ip.b++
			ip.c = 0
			if ip.b == 256 {
				ip.a++
				ip.b = 0
				if ip.a == 256 {
					log.Panic("ip段超出")
				}
			}
		}
		for i := range 256 {
			singChan <- struct{}{}
			ip.d = uint32(i)
			go check(ip.toString())
		}
		fmt.Printf("\r当时的ip为%v\t\t当前堆内最大延迟为%v\t\t完成了%.2f%%\t\t堆里有%d", ip, (*h)[0].Ms, float64(num)/float64(index)*100, len(*h))
		num++
		ip.c++
	}

	time.Sleep(1 * time.Second)
	close(singChan)
	close(heapChan)
	for is := range *h {
		fmt.Printf("\r当时的ip为%s\t\t完成度:%.2f%%\t\t成功上传了%d个", (*h)[is].IP, float64(is)/float64(C.SpeedNum)*100, uploadIP)
		CheckAgain((*h)[is].IP)
	}
	fmt.Println()
	log.Println("结束")
}

func check(ip string) {
	defer func() {
		<-singChan
	}()
	for range 2 {
		m := ping(ip)
		if m == 1000 {
			return
		}
	}
	t := time.Now()
	ray, err := RealPing(ip)
	if err != nil {
		return
	}
	ti := time.Since(t).Milliseconds()
	if ti > 1000 {
		return
	}
	if ray != C.Ray {
		return
	}
	heapChan <- Item{IP: ip, Ms: ti}
}
func inChan(h *MaxHeap) {
	for i := range heapChan {
		if h.Len() < C.SpeedNum {
			heap.Push(h, &i)
		} else if (*h)[0].Ms > i.Ms {
			heap.Pop(h)
			heap.Push(h, &i)
		}
	}
}
func CheckAgain(ip string) {
	ti := time.Now()
	ray, err := RealPing(ip)
	if err != nil {
		return
	}
	tit := time.Since(ti).Milliseconds()
	if tit > 1000 {
		return
	}

	f, err := SpeedTestDownload(ip)
	if err != nil {
		return
	}
	if f < C.SpeedLimit {
		return
	}
	for range 5 {
		err = uploadIp(ip, "443", ray, strconv.FormatInt(tit, 10))
		if err == nil {
			uploadIP++
			return
		}
	}

}
