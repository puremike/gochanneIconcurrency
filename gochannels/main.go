package main

import "fmt"

func main() {

	userch := make(chan string, 2)
	userch <- "Bob"
	userch <- "Scop"

	// After a channel is consumed, the buffer becomes full. Adding to the channel will result in a deadlock

	fmt.Println(<-userch)
	// Buffered channels will only block if the buffer is full. If you do not provide a buffer, it will always be full. Well said.

	msg := make(chan string)
	go sendMessage(msg)
	// message := <-msg
	fmt.Println(<-msg)
}


func sendMessage(ch chan string) {
	ch <- "Did you receive my secret message?"; // Send a message to the channel
}