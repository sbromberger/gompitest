package main

import (
	"fmt"

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
		t0 := mpi.WorldTime()
		o.SendString(str, 1, 0)
		t1 := mpi.WorldTime()
		_, _ = o.RecvString(1, 0)
		t2 := mpi.WorldTime()

		fmt.Printf("sent in %v µs, round trip %v µs\n", (t1-t0) * 100000, (t2-t0) * 100000)
	} else {
		s, _ := o.RecvString(0, 0)
		o.SendString(s, 0, 0)
		_, _ = o.RecvString(0, 0)
		o.SendString(s, 0, 0)
	}

	mpi.Stop()
}
