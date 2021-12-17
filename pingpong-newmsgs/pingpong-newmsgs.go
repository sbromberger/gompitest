package main

import (
	"fmt"
	"os"
	"strconv"

	mpi "github.com/sbromberger/gompi"
	messages "github.com/sbromberger/gompitest/messages2"
)

func main() {
	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic("Invalid number of iterations")
	}

	str := "Hello this is a message"

	mpi.Start(true)
	fmt.Println("started")
	o := mpi.NewCommunicator(nil)
	myRank := o.Rank()

	node := messages.NewNode(myRank, o, 1)
	t0 := mpi.WorldTime()
	for i := 0; i < n; i++ {
		if myRank == 0 {
			msg := messages.Msg{Remote: 1, Tag: 0, Bytes: []byte(str)}
			node.Send(&msg)
			node.Recv()
			node.Send(&msg)
			node.Recv()

		} else {
			rmsg := node.Recv()
			rmsg.Remote = 0
			node.Send(&rmsg)
			node.Recv()
			node.Send(&rmsg)
		}
	}
	t1 := mpi.WorldTime()

	if myRank == 0 {
		fmt.Printf("elapsed %v s, average %v µs\n", (t1 - t0), (t1-t0)/float64(n)*1000000)
	}
	mpi.Stop()
}
