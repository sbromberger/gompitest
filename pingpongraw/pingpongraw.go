package main

import (
	"fmt"
	"time"

	mpi "github.com/sbromberger/gompi"
)

func main() {
	str := "Hello this is a message"

	mpi.Start(true)
	fmt.Println("started")
	o := mpi.NewCommunicator(nil)
	myRank := o.Rank()

	if myRank == 0 {
		o.SendString(str, 1, 0)
		_, _ = o.RecvString(1, 0)
		fmt.Println("rank 0")
		t0 := time.Now()
		o.SendString(str, 1, 0)
		t1 := time.Now()
		_, _ = o.RecvString(1, 0)
		t2 := time.Now()

		fmt.Printf("sent in %v, round trip %v\n", t1.Sub(t0), t2.Sub(t0))
	} else {
		s, _ := o.RecvString(0, 0)
		o.SendString(s, 0, 0)
		_, _ = o.RecvString(0, 0)
		o.SendString(s, 0, 0)
	}

	mpi.Stop()
}
