package main

import (
	"fmt"
	"os"
	"strconv"

	mpi "github.com/sbromberger/gompi"
	messages "github.com/sbromberger/gompitest/messages-recv-chan"
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
	node.Launch()
	t0 := mpi.WorldTime()
	for i := 0; i < n; i++ {
		if myRank == 0 {
			bs := []byte(str)
			msg := messages.Msg{Remote: 1, Tag: 0, Bytes: &bs}
			// fmt.Println("node 0 sending 1")
			node.Send(&msg)
			// fmt.Println("node 0 sent 1")
			// fmt.Println("node 0 receiving 1")
			<-node.Inbox
			// fmt.Println("node 0 received 1")
		} else {
			// fmt.Println("node 1 receiving 1")
			rmsg := <-node.Inbox
			// fmt.Println("node 1 received 1")
			rmsg.Remote = 0
			// fmt.Println("node 1 sending 1")
			node.Send(rmsg)
			// fmt.Println("node 1 sent 1")
		}
	}
	t1 := mpi.WorldTime()

	if myRank == 0 {
		fmt.Printf("elapsed %v s, average %v µs\n", (t1 - t0), (t1-t0)/float64(n)*1000000)
	}

	o.Barrier()
	node.Terminate()
	mpi.Stop()
}
