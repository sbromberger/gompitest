package main

import (
	"fmt"
	"os"
	"strconv"

	mpi "github.com/sbromberger/gompi"
	"github.com/sbromberger/gompitest/messages"
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

	node := messages.NewNode(myRank, &o, 1)
	node.Launch()
	t0 := mpi.WorldTime()
	for i := 0; i < n; i++ {

		if myRank == 0 {
			msg := messages.Msg{Remote: 1, Tag: 0, Bytes: []byte(str)}
			node.Outbox <- msg
			<-node.Inbox
			node.Outbox <- msg
			<-node.Inbox

		} else {
			rmsg := <-node.Inbox
			rmsg.Remote = 0
			node.Outbox <- rmsg
			<-node.Inbox
			node.Outbox <- rmsg
		}
	}
	t1 := mpi.WorldTime()

	if myRank == 0 {
		fmt.Printf("total elapsed %v s, average %v µs\n", t1-t0, (t1-t0)*1000000/float64(n))
	}
	node.Terminate()
	mpi.Stop()
}
