package main

import (
	"math"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	stressWorkers []chan bool // Goroutine을 제어하는 채널 리스트
	mutex         sync.Mutex  // 동시성 제어를 위한 Mutex
)

func main() {
	r := gin.Default()

	// API 엔드포인트 정의
	r.GET("/home", index)
	r.POST("/home/stress", startStress)             // CPU 부하 시작
	r.POST("/home/stop-stress", stopStress)         // CPU 부하 중지
	r.GET("/home/cpu-usage", func(c *gin.Context) { // 현재 CPU 코어 사용량 확인
		c.JSON(http.StatusOK, gin.H{
			"num_cpu":       runtime.NumCPU(),
			"num_goroutine": runtime.NumGoroutine(),
		})
	})

	r.Run(":8080") // 서버 실행
}

// 인덱스 페이지
func index(c *gin.Context) {
	html := `
		<!DOCTYPE html>
		<html lang="ko">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Welcome</title>
		</head>
		<body>
			<h1>Welcome to Helloworld</h1>
		</body>
		</html>`

	c.Data(200, "text/html; charset=utf-8", []byte(html))
}

// CPU 부하를 발생시키는 작업
func stressCPU(stop chan bool) {
	for {
		select {
		case <-stop:
			return // 채널에서 신호를 받으면 종료
		default:
			_ = math.Sqrt(float64(time.Now().UnixNano())) // CPU 연산을 지속적으로 수행
		}
	}
}

// /stress 요청 시 CPU 부하 발생
func startStress(c *gin.Context) {
	mutex.Lock()
	defer mutex.Unlock()

	stop := make(chan bool) // Goroutine을 중지할 채널
	stressWorkers = append(stressWorkers, stop)

	go stressCPU(stop) // 백그라운드에서 CPU 부하 발생

	c.JSON(http.StatusOK, gin.H{"message": "CPU stress started", "workers": len(stressWorkers)})
}

// /stop-stress 요청 시 모든 부하 중지
func stopStress(c *gin.Context) {
	mutex.Lock()
	defer mutex.Unlock()

	for _, stop := range stressWorkers {
		stop <- true // 모든 Goroutine을 종료
	}
	stressWorkers = nil // 리스트 초기화

	c.JSON(http.StatusOK, gin.H{"message": "CPU stress stopped"})
}
