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
	myRank := o.Rank()

	msg := fmt.Sprintf("Hello from rank %d!", myRank)
	t0 := time.Now()
	if myRank == 0 {
		// log.Warnf("%d: sending", myRank)
		o.SendString(msg, 1, 0)
		_, _ = o.RecvString(1, 1)
	} else {
		// log.Warnf("%d: receiving", myRank)
		_, _ = o.RecvString(0, 0)
		o.SendString(msg, 0, 1)
	}
	t1 := time.Since(t0)

	o.Barrier()
	mpi.Stop()
	fmt.Printf("%d: Time elapsed: %v\n", myRank, t1)

}
