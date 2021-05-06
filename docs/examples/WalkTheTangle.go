package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/iotaledger/goshimmer/client"
	"github.com/iotaledger/goshimmer/packages/mana"
	"github.com/mr-tron/base58"
)

var futureCone map[string][]string
var futureTips map[string]bool
var manaList map[string]float64
var approvingNodes []string
var nodeAPIURL string = "http://49.12.35.135:8080"
var manaApproving float64 = 0
var totalActiveCmana float64 = 0

//nodeAPIURL := "http://localhost:8080"
var clientNode *client.GoShimmerAPI = client.NewGoShimmerAPI(nodeAPIURL, client.WithHTTPClient(http.Client{Timeout: 90 * time.Second}))

func main() {
	iterations := 1000
	futureCone = make(map[string][]string)
	futureTips = make(map[string]bool)
	manaList = make(map[string]float64)

	cMana, err := clientNode.GetOnlineConsensusMana()
	if err != nil {
		// return error
	}
	println(cMana.Online)

	for _, m := range cMana.Online {
		totalActiveCmana += m.Mana
	}
	fmt.Printf("Total amount of active consensus mana: %g \n", totalActiveCmana)
	fmt.Println("Issuing new message")
	messageID, _ := clientNode.Data([]byte("Hello GoShimmer World"))
	fmt.Println("Message ID", messageID)
	//time.Sleep(5000 * time.Millisecond)
	//messageID := "3nroCkpt3He5hxF5X94Mo2kNDxkgC3cHScUs6zLxHtXd"
	//messageID := "EQfZNRdXqNSZPJ7j7CgLHJLQLb7wNRX4JzThrNvvt2iM"
	start := time.Now()
	futureTips[messageID] = true
	WalkOneStepAndGet(messageID)
	//fmt.Println("Future Tips:", futureTips)
	//fmt.Println("Approving Nodes:", approvingNodes)
	//fmt.Println("Number of approving nodes:", len(approvingNodes))
	for k := 0; k < iterations; k++ {
		for tip, flag := range futureTips {
			if flag == true {
				WalkOneStepAndGet(tip)
			}
		}
		fmt.Printf("Number of approving nodes: %d \n", len(approvingNodes))
		fmt.Printf("Number of approving msgss: %d \n", len(futureTips))
		fmt.Printf("Percentage of active Consensus Mana approving: %g \n", manaApproving/totalActiveCmana)
		t := time.Now()
		elapsed := t.Sub(start)
		fmt.Println("Time passed since issuing (milliseconds):", elapsed.Milliseconds())
	}
}

func WalkOneStepAndGet(messageID string) {

	//fmt.Println(messageID)
	messageData, err := clientNode.GetMessage(messageID)
	if err != nil {
		return
	}
	nodeID, _ := mana.IDFromPubKey(messageData.IssuerPublicKey)
	nodeIDString := base58.Encode(nodeID.Bytes())
	numberOld := len(approvingNodes)
	approvingNodes = unique(append(approvingNodes, nodeIDString))
	numberNew := len(approvingNodes)
	if numberOld < numberNew {
		manaReply, err := clientNode.GetManaFullNodeID(nodeIDString)
		manaList[nodeIDString] = manaReply.Consensus
		manaApproving += manaList[nodeIDString]
		if err != nil {
			// return error
		}
	}
	StrongApprovers := unique(messageData.StrongApprovers)
	if len(StrongApprovers) == 0 {
		return
	}
	futureCone[messageID] = StrongApprovers
	for _, ApproverID := range StrongApprovers {
		_, ok := futureTips[ApproverID]
		if ok == false {
			futureTips[ApproverID] = true
		}
	}
	futureTips[messageID] = false

}
func unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
