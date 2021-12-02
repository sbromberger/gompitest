package main

import (
	"fmt"
	"time"

	mpi "github.com/sbromberger/gompi"
	messages "github.com/sbromberger/gompitest/messages2"
)

func main() {
	str := "Hello this is a message"

	mpi.Start(true)
	fmt.Println("started")
	o := mpi.NewCommunicator(nil)
	myRank := o.Rank()

	node := messages.NewNode(myRank, &o, 1)
	if myRank == 0 {
		msg := messages.Msg{Remote: 1, Tag: 0, Bytes: []byte(str)}
		node.Send(msg)
		node.Recv()
		fmt.Println("rank 0")
		t0 := time.Now()
		node.Send(msg)
		t1 := time.Now()
		node.Recv()
		t2 := time.Now()

		fmt.Printf("sent in %v, round trip %v\n", t1.Sub(t0), t2.Sub(t0))
	} else {
		rmsg := node.Recv()
		rmsg.Remote = 0
		node.Send(rmsg)
		node.Recv()
		node.Send(rmsg)
	}

	mpi.Stop()
}
