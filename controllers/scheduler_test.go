package controllers


import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"
	"testing"
)

var wg sync.WaitGroup

var urls = []string{
	"http://119.23.67.53:8080/hello",
	"http://119.23.67.53:8080/hello1",
	"http://119.23.67.53:8080/hello2",
}

// 运行命令：go test scheduler_test.go -v -bench=.

// 这里进行压力测试
func BenchmarkScheduler(b *testing.B){
	rand.Intn(2)
	var a = 0
	for i:=0;i<1000000;i++{
		wg.Add(1)
		go func(i int){
			defer wg.Done()
			resp,err := http.Get(urls[rand.Intn(3)])
			if err!=nil{
				b.Error(err)
				a++
				return
			}
			defer resp.Body.Close()
			data,err := ioutil.ReadAll(resp.Body)
			if err!=nil{
				b.Error(err)
				a++
				return
			}
			b.Log(string(data))
		}(i)
	}
	wg.Wait()
	b.Log("错误率为：",a/10000)
}