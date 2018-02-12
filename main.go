package main

import (
	"context"
	"sync"
	"time"
	"os"

	"os/signal"
	"math/rand"

	"github.com/sirupsen/logrus"
)

type Car struct {
	ID           string
	Timing		 time.Time
}

type Road struct {
	Sleep	time.Duration
	Cars	int
}

type Roads []Road

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	inputRoads := Roads{
		Road{time.Second * 2, rand.Intn(6)},
		Road{time.Second * 2, 1},
		Road{time.Second * 3, 1},
		Road{time.Second * 4, 10},
	}
	outputRoads := Roads{
		Road{time.Second * 2, rand.Intn(6)},
		Road{time.Second, 1},
		Road{time.Hour, 1},
		Road{time.Second, 10},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	go func() {
		for {
			select {
			case _, ok := <- c:
				if ok {
					<-ctx.Done()
					return
				}
			}
		}
	} ()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go trafficWorker(inputRoads, outputRoads, ctx, wg)
	wg.Wait()
}

func trafficWorker(inputRoads, outputRoads Roads, ctx context.Context, wg *sync.WaitGroup) {
	defer func() {
		logrus.Info("Traffic was stoped")
		wg.Done()
	}()
	childWG := &sync.WaitGroup{}
	childWG.Add(2)
	circle := inputsWorker(ctx, childWG, inputRoads)
	go outputsWorker(ctx, childWG, outputRoads, circle)

	childWG.Wait()
}

//
//
func inputsWorker(ctx context.Context, wg *sync.WaitGroup, roads Roads) <-chan Car {
	logrus.Info("Car generator starting")
	circle := make(chan Car, 8)
	go func() {
		defer func(wg *sync.WaitGroup) {
			logrus.Info("Car generator was stoped")
			close(circle)
			wg.Done()
		}(wg)



		for i := 0; ; i++ {
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
	return circle
}
//
//func outputsWorker(wg *sync.WaitGroup, inChan <-chan int) {
//	defer func(wg *sync.WaitGroup) {
//		logrus.Info("Car releaser was stoped")
//		wg.Done()
//	}(wg)
//
//	for i := range inChan {
//		logrus.Info("Release cars:", i)
//	}
//}
//
//func trafficWorker(ctx context.Context, wg *sync.WaitGroup) {
//	defer func() {
//		logrus.Info("Traffic was stoped")
//		wg.Done()
//	}()
//
//	childWG := &sync.WaitGroup{}
//	childWG.Add(2)
//	inChan := inputsWorker(ctx, childWG)
//	go outputsWorker(childWG, inChan)
//
//	childWG.Wait()
//}
//
//func main() {
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
//	defer cancel()
//
//	wg := &sync.WaitGroup{}
//	wg.Add(1)
//	go trafficWorker(ctx, wg)
//	wg.Wait()
//}
