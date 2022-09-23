//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer szenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"time"
)

// Create channel for receiving tweets
var chanTweets = make(chan *Tweet, 1)

func producer(stream Stream) {
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			// Close the channel if there is no data again
			close(chanTweets)
			return
		}

		chanTweets <- tweet
	}
}

func consumer() {
	// Loop over channel
	for t := range chanTweets {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}
}

func main() {
	start := time.Now()
	stream := GetMockStream()

	// Producer
	// Make the producer as a go routine
	go producer(stream)

	// Consumer
	consumer()

	fmt.Printf("Process took %s\n", time.Since(start))
}
