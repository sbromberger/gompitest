package main

import (
	"fmt"

	mpi "github.com/sbromberger/gompi"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.WarnLevel)
	fmt.Println("starting")
	mpi.Start(true)
	o := mpi.NewCommunicator(nil)
	myRank := o.Rank()
	nRanks := o.Size()

	fmt.Printf("Rank %d of %d\n", myRank, nRanks)
	// o.Barrier()
	bytes := []byte("this is a test")

	if myRank == 0 {
		fmt.Printf("0: sending bytes\n")
		o.SendBytes(bytes, 1, 1)
	} else {
		fmt.Printf("1: mrecv\n")
		msg, _ := o.MrecvBytes(0, 1)
		fmt.Printf("1:    received %s\n", msg)
	}
	o.Barrier()
	mpi.Stop()
}
