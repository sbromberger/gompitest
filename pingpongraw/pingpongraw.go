package main

import (
	"fmt"
	"os"
	"strconv"

	mpi "github.com/sbromberger/gompi"
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

	t0 := mpi.WorldTime()
	for i := 0; i < n; i++ {
		if myRank == 0 {
			o.SendString(str, 1, 0)
			_, _ = o.RecvString(1, 0)
			o.SendString(str, 1, 0)
			_, _ = o.RecvString(1, 0)

		} else {
			s, _ := o.RecvString(0, 0)
			o.SendString(s, 0, 0)
			_, _ = o.RecvString(0, 0)
			o.SendString(s, 0, 0)
		}
	}
	t1 := mpi.WorldTime()

	fmt.Printf("elapsed %v s, average %v µs\n", t1-t0, (t1-t0)/float64(n)*1000000)
	mpi.Stop()
}
