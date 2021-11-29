package main

import (
	"fmt"

	mpi "github.com/sbromberger/gompi"
	"github.com/sbromberger/gompitest/messages"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.WarnLevel)
	fmt.Println("starting")
	mpi.Start(true)
	o := mpi.NewCommunicator(nil)
	nprocs := o.Size()
	myRank := o.Rank()

	myNode := messages.NewNode(myRank, &o, 1)
	o.Barrier()

	fmt.Printf("%d: launching\n", myRank)
	myNode.Launch()

	o.Barrier()

	fmt.Printf("%d: sending\n", myRank)
	dest := (myRank + 1) % nprocs
	mymsg := []byte(fmt.Sprintf("hello from rank %d to rank %d!", myRank, dest))
	fmt.Printf("  %d: sending message %s to %d\n", myRank, mymsg, dest)
	myNode.Outbox <- messages.Msg{Bytes: mymsg, Remote: dest, Tag: myRank}

	fmt.Printf("%d: receiving\n", myRank)
	m := <-myNode.Inbox

	fmt.Printf("  %d: message = %s\n", myRank, m)
	fmt.Printf("%d: terminating\n", myRank)
	myNode.Terminate()
	fmt.Printf("%d: exiting\n", myRank)
	o.Barrier()
	mpi.Stop()
}
