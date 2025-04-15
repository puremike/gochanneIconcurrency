package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type Server struct {
	msgch    chan Message
	quitch   chan struct{}
	errch    chan error
	wgWorker sync.WaitGroup
}

type Message struct {
	From    string
	Payload string
}

// WITH GRACEFUL SHUTDOWN

func main() {

	s := Server{
		// msgch:  make(chan Message), //unbuffered channel
		msgch:  make((chan Message), 10), // buffered channel
		quitch: make(chan struct{}),
		errch:  make(chan error),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := s.startAndListen(ctx); err != nil {
			fmt.Printf("Server failed : %v\n", err)
		}
	}()

	// Worker pool (3 workers)
	s.StartWorker(ctx, 3)

	// send message
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			if err := sendMessageWithRetry(s.msgch, fmt.Sprintf("Hello Gophers %d", i), 3); err != nil {
				s.errch <- fmt.Errorf("send failed :%v", err)
				return
			}

			time.Sleep(500 * time.Millisecond)
		}

	}()

	// Simulate external shutdown after 3 seconds
	time.AfterFunc(time.Second*3, func() {
		fmt.Println("Initiating graceful shutdown...")
		graceFulShutdown(s.quitch)
		cancel() // propagate cancellation
	})

	// error handling
	go func() {
		for err := range s.errch {
			fmt.Printf("Error : %v\n", err)
		}
	}()

	wg.Wait()
	s.wgWorker.Wait()
	close(s.errch)

	fmt.Println("All goroutines completed. Exiting...")
}

func (s *Server) startAndListen(ctx context.Context) error {
	for {
		select {
		case msg, ok := <-s.msgch:
			if !ok {
				fmt.Println("[Listener] Message channel closed .....")
			}
			fmt.Printf("[Listener] Received message from: %s\nPayload: %s\n", msg.From, msg.Payload)

		case <-s.quitch:
			fmt.Println("[Listener] Shutdown Signal Received ......")
			return nil

		case <-ctx.Done():
			fmt.Println("[Listener] Context Cancelled .....")
			return ctx.Err()

		case err := <-s.errch:
			fmt.Printf("[Listener] Server error : %v", err)
		}
	}
}

// Worker pool implementation
func (s *Server) StartWorker(ctx context.Context, n int) {
	for i := 0; i < n; i++ {
		s.wgWorker.Add(1)
		go func(id int) {
			defer s.wgWorker.Done()
			for {
				select {
				case msg, ok := <-s.msgch:
					if !ok {
						fmt.Printf("[Worker] Worker %d: channel closed\n", id)
						return
					}
					fmt.Printf("[Worker] Worker %d processing: %s\n", id, msg.Payload)
					time.Sleep(1 * time.Second) // Simulate work

				case <-ctx.Done():
					fmt.Printf("[Worker] Worker %d: context cancelled\n", id)
					return
				}
			}

		}(i)
	}

}

// Send with retry logic
func sendMessageWithRetry(ch chan Message, payload string, maxRetries int) error {
	msg := Message{
		From:    "Mike",
		Payload: payload,
	}
	for i := 0; i < maxRetries; i++ {
		select {
		case ch <- msg:
			return nil
		case <-time.After(100 * time.Millisecond):
			fmt.Printf("Retry %d for: %s\n", i+1, payload)
		}
	}
	return errors.New("send timeout")
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
