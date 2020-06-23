package controllers

import (
	"helloweb/logger"
	"log"
	"net/http"
	"sync"
	"time"
)

// Go来处理每分钟达百万的数据请求
// 设置工作者的数量
// 框架设置为：工作和工作者
// 在调度器中，每一个工作者都是处于等待状态，（监听状态）
// 有任务就加入到工作池

// MaxWorkers 最大的工作者数量
var MaxWorkers = 10

// DispatchNumControl 用于控制并发处理的协程数
var DispatchNumControl = make(chan bool, 30)

// SchedulerWg 调度协程组
var SchedulerWg sync.WaitGroup

// Payload 负载
type Payload struct {
	W  http.ResponseWriter
	R  *http.Request
	Fn func(http.ResponseWriter, *http.Request) (interface{}, error)
}

// Job 工作
type Job struct {
	Payload Payload
	Finish  chan bool
}

// JobChannel 任务队列，用于接收任务
var JobChannel chan Job

func init() {
	JobChannel = make(chan Job, MaxWorkers)
	log.SetFlags(log.Ldate | log.Lshortfile)
}

// Worker 工作者
type Worker struct {
	WorkerPool chan chan Job // 工作池--每个元素是一个工作者的私有任务
	JobChannel chan Job      // 每个工作者单元包含一个任务 用于获取任务
	Quit       chan bool     // 退出信号
	No         int           // 编号
}

// NewWorker 创建工作者
func NewWorker(WorkerPool chan chan Job, no int) *Worker {
	return &Worker{
		WorkerPool: WorkerPool,
		JobChannel: make(chan Job),
		Quit:       make(chan bool),
		No:         no,
	}
}

// Start 任务开始，用于监听
func (w *Worker) Start() {
	go func() {
		for {
			// 这里会发生阻塞的
			w.WorkerPool <- w.JobChannel
			logger.Z.Info("w.WorkerPool <- JobChannel")
			select {
			case job := <-w.JobChannel:
				// 接收到任务
				logger.Z.Info("job :=<-w.JobChannel 接收到任务")
				// 处理任务
				result, err := job.Payload.Fn(job.Payload.W, job.Payload.R)
				logger.Z.Info(result.(string))
				if err != nil {
					if err = ErrorResp(job.Payload.W, job.Payload.R, err.Error()); err != nil {
						logger.Z.Error(job.Payload.R.URL.String() + ",运行报错")
					}
				} else {
					if err = SuccessResp(job.Payload.W, job.Payload.R, result); err != nil {
						logger.Z.Error(job.Payload.R.URL.String() + ",运行报错")
					}
				}
				logger.Z.Info("释放协程")
				<-DispatchNumControl
				job.Finish <- true
			case <-w.Quit:
				// 接收到退出信号
				return
			}
		}
	}()
}

// Stop 监听停止
func (w *Worker) Stop() {
	go func() {
		w.Quit <- true
	}()
}

// Dispatcher 调度器
type Dispatcher struct {
	// WorkerPool 工作池
	WorkerPool chan chan Job
	// 工作者数量
	MaxWorker int
}

// NewDispatcher 创建调度中心
func NewDispatcher() *Dispatcher {
	return &Dispatcher{WorkerPool: make(chan chan Job, MaxWorkers), MaxWorker: MaxWorkers}
}

// Run 初始化工作池
func (d *Dispatcher) Run() {
	for i := 1; i <= MaxWorkers; i++ {
		// 初始化，这里传输d.WorkerPool，是为了控制工作者是否工作的
		worker := NewWorker(d.WorkerPool, i)
		// 开启监听
		worker.Start()
	}
	go d.dispatcher()
}

// dispatcher 执行调度
func (d *Dispatcher) dispatcher() {
	for {
		select {
		case job := <-JobChannel:
			logger.Z.Info("监听到有任务")
			go func(job Job) {
				// 等待空闲worker(任务多的时候会阻塞这里)
				logger.Z.Info("等待空闲worker (任务多的时候会阻塞这里)")
				jobChannel := <-d.WorkerPool
				logger.Z.Info("jobChannel := <-d.WorkerPool")
				// 将任务放到上述worker的私有任务channel中
				jobChannel <- job
				logger.Z.Info("将工作放到工作池")
			}(job)
		}
	}
}

// Limit 限制协程的数量
func Limit(work Job) bool {
	select {
	case <-time.After(time.Second * 1):
		logger.Z.Info("我很忙")
		work.Finish <- true
		return false
	case DispatchNumControl <- true:
		logger.Z.Info("JobChannel <- work")
		JobChannel <- work
		return true
	}
}