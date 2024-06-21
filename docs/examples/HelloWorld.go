package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/iotaledger/goshimmer/client"
)

func main() {
	nodeAPIURL := "http://49.12.35.135:8080"
	//nodeAPIURL := "http://localhost:8080"

	client := client.NewGoShimmerAPI(nodeAPIURL, client.WithHTTPClient(http.Client{Timeout: 90 * time.Second}))
	messageID, _ := client.Data([]byte("Hello GoShimmer World"))
	// print messageID, also can check message in explorer via browser
	fmt.Println(messageID)
	message, err := client.GetMessage(messageID)
	if err != nil {
		// return error
	}

	// this might be nil if the message wasn't available

	// will print "Hello GoShimmer World"
	fmt.Println(string(message.Payload))
}
