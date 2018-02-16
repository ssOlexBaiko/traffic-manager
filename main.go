package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"math/rand"
	"os/signal"

	"github.com/sirupsen/logrus"
)

type Car struct {
	ID     string
	Timing time.Time
}

type Road struct {
	Sleep time.Duration
	Cars  int
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

	circle := make(chan Car, 8)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		select {
		case _, ok := <-c:
			if ok {
				cancel()
				return
			}
		}
	}()
	wg := &sync.WaitGroup{}
	wg.Add(len(inputRoads) + len(outputRoads))
	for _, road := range inputRoads {
		go inputWorker(ctx, wg, road, circle)
	}

	for _, road := range outputRoads {
		go outputWorker(ctx, wg, road, circle)
	}

	// TODO: Fix stop order
	// TODO: Add traffic_circle

	wg.Wait()
}

func inputWorker(ctx context.Context, wg *sync.WaitGroup, road Road, circle chan Car) {
	logrus.Info("Car generator starting")
	ticker := time.NewTicker(road.Sleep / time.Duration(road.Cars))
	defer func() {
		logrus.Info("Car generator was stoped")
		ticker.Stop()
		wg.Done()
	}()

	for i := 0; ; i++ {
		select {
		case <-ctx.Done():
			logrus.Info("Traffic receive ctx.Done signal. Waiting for roads")
			return
		case <-ticker.C:
			carID := fmt.Sprintf("Car #%d from road #%d\n", i, road.Cars)
			newCar := Car{carID, time.Now()}
			circle <- newCar
			logrus.Info("Generate car: ", i)
		}
	}
}

func outputWorker(ctx context.Context, wg *sync.WaitGroup, road Road, circle <-chan Car) {
	defer func() {
		logrus.Info("Car releaser was stoped")
		wg.Done()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			for i := 0; i < road.Cars; i++ {
				car := <-circle
				fmt.Println("We took out off road ", car.ID)
			}
			time.Sleep(road.Sleep)
		}
	}
}
