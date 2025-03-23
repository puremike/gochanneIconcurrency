package main

import (
	"fmt"
	"sync"
	"time"
)

type Server struct {
	msgch  chan Message
	quitch chan struct{}
}

type Message struct {
	From    string
	Payload string
}

// WITH GRACEFUL SHUTDOWN

func main() {

	s := Server{
		msgch:  make(chan Message),
		quitch: make(chan struct{}),
	}

	var wg = &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		s.startAndListen()
		wg.Done()
	}()

	sendMessage(s.msgch)

	time.Sleep(2 * time.Second)
	fmt.Println("Triggering graceful shutdown...")

	graceFulShutdown(s.quitch)

	wg.Wait()
	close(s.msgch)

	fmt.Println("All goroutines completed. Exiting...")
}

func (s *Server) startAndListen() {
	for {
		select {
		case msg, ok := <-s.msgch:
			if !ok {
				fmt.Println("Message channel closed. Exiting Listener .....")
			}
			fmt.Printf("Received message from: %s\npayload: %s\n", msg.From, msg.Payload)

		case <-s.quitch:
			fmt.Println("Received shutdown signal. Closing message channel ....")
			return
		}
	}
}

func sendMessage(ch chan Message) {
	msg := Message{
		From:    "Joe",
		Payload: "Hello Joe",
	}
	ch <- msg
}

func graceFulShutdown(ch chan struct{}) {
	close(ch)
}

// // WITHOUT GRRACEFUL SHUTDOWN
// type Server struct {
// 	msgch  chan Message
// 	quitch chan struct{}
// }

// type Message struct {
// 	From    string
// 	Payload string
// }

// func main() {

// 	s := Server{
// 		msgch:  make(chan Message),
// 		quitch: make(chan struct{}),
// 	}

// 	var wg = &sync.WaitGroup{}
// 	wg.Add(1)

// 	go func() {
// 		s.startAndListen()
// 		wg.Done()
// 	}()

// 	sendMessage(s.msgch)

// 	close(s.msgch)
// 	wg.Wait()

// 	fmt.Println("All goroutines completed. Exiting...")
// }

// func (s *Server) startAndListen() {
// 	for msg := range s.msgch {
// 		fmt.Printf("Received message from: %s\npayload: %s\n", msg.From, msg.Payload)
// 	}

// 	fmt.Println("\nListener stopped: Channel closed")
// }

// func sendMessage(ch chan Message) {
// 	msg := Message{
// 		From:    "Joe",
// 		Payload: "Hello Joe",
// 	}
// 	ch <- msg
// }

// NOTE: If multiple goroutines are sending, we must ensure they are done before closing --> close(ch). Otherwise, we close --> close(ch) the channel and for the goroutine to complete wait --- wg.Wait()

// func main() {
// 	now := time.Now()
// 	userId := 10

// 	ch := make(chan string, 128) // you must pass a buffer, else there will be deadlock

// 	var wg = &sync.WaitGroup{}

// 	// A wrapper function to wrap all the goroutines
// 	fetchUser := func(fn func(id int, ch chan string), id int) {
// 		wg.Add(1)
// 		go func() {
// 			fn(id, ch)
// 			wg.Done()
// 		}()
// 	}

// 	fetchUser(fetchUserData, userId)
// 	fetchUser(fetchUserRecommendation, userId)
// 	fetchUser(fetchUserLikes, userId)

// 	wg.Wait()
// 	close(ch)

// 	for i := range ch {
// 		fmt.Println("Went through:", i)
// 	}

// 	// userData := fetchUserData(userId)
// 	// userRecs := fetchUserRecommendation(userId)
// 	// userLikes := fetchUserLikes(userId)

// 	// fmt.Println(userData)
// 	// fmt.Println(userRecs)
// 	// fmt.Println(userLikes)

// 	fmt.Println(time.Since(now))
// }

// func fetchUserData(id int, ch chan string) {
// 	time.Sleep(80 * time.Millisecond)
// 	ch <- "user data"
// }

// func fetchUserRecommendation(id int, ch chan string) {
// 	time.Sleep(120 * time.Millisecond)
// 	ch <- "user recommendations"
// }

// func fetchUserLikes(id int, ch chan string) {
// 	time.Sleep(50 * time.Millisecond)
// 	ch <- "user likes"
// }
