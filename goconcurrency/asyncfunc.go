package main

import (
	"fmt"
	"sync"
	"time"
)

func asyncFuncWg() {
	now := time.Now()
	userId := 10

	ch := make(chan string, 128)
	var wg sync.WaitGroup

	startTask := func(fn func(int, chan string), id int) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fn(id, ch)
		}()
	}

	startTask(fetchUserData1, userId)
	startTask(fetchUserRecommendation1, userId)
	startTask(fetchUserLikes1, userId)

	wg.Wait()
	close(ch)

	for i := range ch {
		fmt.Println("Went through:", i)
	}

	fmt.Println(time.Since(now))
}

func fetchUserData1(id int, ch chan string) {
	time.Sleep(80 * time.Millisecond)
	ch <- "user data"
}

func fetchUserRecommendation1(id int, ch chan string) {
	time.Sleep(120 * time.Millisecond)
	ch <- "user recommendations"
}

func fetchUserLikes1(id int, ch chan string) {
	time.Sleep(50 * time.Millisecond)
	ch <- "user likes"
}
