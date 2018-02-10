package main

import (
	"context"
	//"math/rand"
	//"sync"
	//
	//"github.com/sirupsen/logrus"
	//"fmt"
	//"time"
	"math/rand"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

func inputsWorker(ctx context.Context, wg *sync.WaitGroup) <-chan int {
	var GlobalCars = make(chan int, 8)

	logrus.Info("Car generator starting")
	go func() {
		defer func(wg *sync.WaitGroup) {
			logrus.Info("Car generator was stoped")
			close(GlobalCars)
			wg.Done()
		}(wg)

		for i := 0; i < 16; i++ {
			select {
			case <-ctx.Done():
				logrus.Info("Traffic receive ctx.Done signal. Waiting for roads")
				return
			case <-time.After(time.Second):
				GlobalCars <- rand.Intn(5)
				logrus.Info("Generate car: ", i)
			}
		}
	}()
	return GlobalCars
}

func outputsWorker(wg *sync.WaitGroup, inChan <-chan int) {
	defer func(wg *sync.WaitGroup) {
		logrus.Info("Car releaser was stoped")
		wg.Done()
	}(wg)

	for i := range inChan {
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
	inChan := inputsWorker(ctx, childWG)
	go outputsWorker(childWG, inChan)

	childWG.Wait()
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go trafficWorker(ctx, wg)
	wg.Wait()
}
