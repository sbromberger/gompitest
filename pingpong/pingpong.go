package main

import (
	"fmt"
	"time"

	mpi "github.com/sbromberger/gompi"
	"github.com/sbromberger/gompitest/messages"
)

func main() {
	str := "Hello this is a message"

	mpi.Start(true)
	fmt.Println("started")
	o := mpi.NewCommunicator(nil)
	myRank := o.Rank()

	node := messages.NewNode(myRank, &o, 1)
	node.Launch()
	if myRank == 0 {
		msg := messages.Msg{Remote: 1, Tag: 0, Bytes: []byte(str)}
		node.Outbox <- msg
		_ = <-node.Inbox
		fmt.Println("rank 0")
		t0 := time.Now()
		node.Outbox <- msg
		t1 := time.Now()
		_ = <-node.Inbox
		t2 := time.Now()

		fmt.Printf("sent in %v, round trip %v\n", t1.Sub(t0), t2.Sub(t0))
	} else {
		rmsg := <-node.Inbox
		rmsg.Remote = 0
		node.Outbox <- rmsg
		_ = <-node.Inbox
		node.Outbox <- rmsg
	}

	node.Terminate()
	mpi.Stop()
}
