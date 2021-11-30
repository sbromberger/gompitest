package main

import (
	"fmt"
	"time"

	mpi "github.com/sbromberger/gompi"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.WarnLevel)
	fmt.Println("starting")
	mpi.Start(true)
	o := mpi.NewCommunicator(nil)
	nranks := o.Size()
	myRank := o.Rank()

	msg := fmt.Sprintf("Hello from rank %d!", myRank)
	t0 := time.Now()
	for i := 0; i < nranks-1; i++ {
		if myRank == i {
			fmt.Printf("%d -> %d\n", i, i+1)
			o.SendString(msg, i+1, 9999)
		}
		if myRank == i+1 {
			o.RecvString(i, 9999)
		}
	}
	t1 := time.Since(t0)

	o.Barrier()
	mpi.Stop()
	if myRank == 0 {
		fmt.Printf("Time elapsed: %v\n", t1)
	}

}
