package main

import (
	"context"
	//"math/rand"
	//"sync"
	//
	//"github.com/sirupsen/logrus"
	//"fmt"
	//"time"
	"github.com/sirupsen/logrus"
	"math/rand"
	"sync"
	"time"
)

var GlobalCars = make(chan int, 8)

func inputsWorker(wg *sync.WaitGroup) {
	defer func(wg *sync.WaitGroup) {
		logrus.Info("Car generator was stoped")
		wg.Done()
	}(wg)
	logrus.Info("Car generator starting")
	for i := 0; i < 16; i++ {
		GlobalCars <- rand.Intn(5)
		logrus.Info("Generate car: ", i)
		time.Sleep(time.Second)
	}
	close(GlobalCars)
}

func outputsWorker(wg *sync.WaitGroup) {
	defer func(wg *sync.WaitGroup) {
		logrus.Info("Car releaser was stoped")
		wg.Done()
	}(wg)
	for i := range GlobalCars {
		logrus.Info("Release cars:", i)
	}

}

func trafficWorker(ctx context.Context, wg *sync.WaitGroup) {
	defer func() {
		logrus.Info("Traffic was stoped")
		wg.Done()
	}()
	childWG := &sync.WaitGroup{}
	childWG.Add(2)
	go inputsWorker(childWG)
	go outputsWorker(childWG)
	childWG.Wait()
	<-ctx.Done()
	logrus.Info("Traffic receive ctx.Done signal. Waiting for roads")
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go trafficWorker(ctx, wg)
	cancel()
	wg.Wait()
}
