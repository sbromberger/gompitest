package main

import (
	"fmt"

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
		t0 := mpi.WorldTime()
		node.Outbox <- msg
		t1 := mpi.WorldTime()
		_ = <-node.Inbox
		t2 := mpi.WorldTime()

		fmt.Printf("sent in %v µs, round trip %v µs\n", (t1-t0) *100000, (t2-t0) * 100000)
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
