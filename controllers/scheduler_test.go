package controllers


import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"
	"testing"
	"log"
)

var wg sync.WaitGroup
var lock sync.Mutex

var urls = []string{
	"http://localhost:8080/hello",
	"http://localhost:8080/hello1",
	"http://localhost:8080/hello2",	
}

// 运行命令：go test scheduler_test.go -v -bench=.

func init(){
	log.SetFlags(log.Ldate|log.Lshortfile)
}

// 这里进行压力测试
func BenchmarkScheduler(b *testing.B){
	var a = 0
	var c = 0
	for i:=0;i<10000;i++{
		wg.Add(1)
		go func(i int){
			defer wg.Done()
			resp,err := http.Get(urls[rand.Intn(3)])
			if err!=nil{
				log.Println(err)
				a++
				return
			}
			defer resp.Body.Close()
			data,err := ioutil.ReadAll(resp.Body)
			if err!=nil{
				log.Println(err)
				a++
				return
			}
			lock.Lock()
			c++
			lock.Unlock()
			log.Println(string(data),c)
		}(i)
	}
	wg.Wait()
	log.Println("成功运行的概率为：",c)
	if a!=0{
		b.Error("错误率为：",a)
	}
}