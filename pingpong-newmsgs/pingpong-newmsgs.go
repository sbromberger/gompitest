package main

import (
	"fmt"

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
		t0 := mpi.WorldTime()
		node.Send(msg)
		t1 := mpi.WorldTime()
		node.Recv()
		t2 := mpi.WorldTime()

		fmt.Printf("sent in %v µs, round trip %v µs\n", (t1-t0)*1e6, (t2-t0)*1e6)
	} else {
		rmsg := node.Recv()
		rmsg.Remote = 0
		node.Send(rmsg)
		node.Recv()
		node.Send(rmsg)
	}

	mpi.Stop()
}
