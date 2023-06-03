package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func(ctx context.Context, cancel context.CancelFunc) {
		<-signalChan
		stopGoroutines(cancel)
	}(ctx, cancel)
	fmt.Println("-----from main function-----")
	checkCtx(ctx)
	//wg.Wait()
	//time.Sleep(10 * time.Second)

}
func checkCtx(ctx context.Context) {
	fmt.Println("----->>from checkCtx function<<-----")
	parallelProcessingNodes := 5
	wg := sync.WaitGroup{}
	semaphore := make(chan struct{}, parallelProcessingNodes)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		fmt.Printf("----->>inside the for loop with i=%d\n", i)
		semaphore <- struct{}{}
		go func(i int) {
			fmt.Printf("----->>starting a goroutine for i=%d\n", i)
			defer func() { <-semaphore }()
			replaceDeleteWait(ctx, i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func replaceDeleteWait(ctx context.Context, i int) {
	select {
	case <-ctx.Done():
		fmt.Println("got a ctx done")
		return
	default:
		fmt.Printf("Proceeding for WaitForNewProcessToBeReady on %d go routine\n", i)
		WaitForNewProcessToBeReady(ctx, i)
	}
}

func WaitForNewProcessToBeReady(ctx context.Context, goi int) {
	fmt.Printf(">>>>>> running the gorutine: %d\n", goi)
	for i := 0; i < 10; i++ {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(1 * time.Second)
			fmt.Printf(">>>>>> now in %d iteration for goroutine: %d\n", i, goi)
		}

	}
}

func stopGoroutines(cancel context.CancelFunc) {
	fmt.Printf("stopping all goroutines as the contaier is being evicted\n")
	fmt.Printf("shutting down %d goroutines\n", runtime.NumGoroutine())
	cancel()
}
