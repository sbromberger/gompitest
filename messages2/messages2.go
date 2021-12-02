package messages

import (
	"fmt"
	"runtime"

	log "github.com/sirupsen/logrus"

	mpi "github.com/sbromberger/gompi"
)

const DEFAULT_TAG = 0

type Msg struct {
	Remote int
	Tag    int
	Bytes  []byte
}

func (m Msg) String() string {
	return fmt.Sprintf("Message: Remote %d, Tag %d, \"%s\"", m.Remote, m.Tag, string(m.Bytes))
}

type Node struct {
	Source int
	comm   *mpi.Communicator
}

func NewNode(source int, comm *mpi.Communicator, bufsize int) *Node {
	node := Node{Source: source, comm: comm}
	return &node
}

func (node *Node) Send(msg Msg) {
	node.comm.SendBytes(msg.Bytes, msg.Remote, msg.Tag)
}

func (node *Node) Recv() Msg {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	recvbytes, status := node.comm.MrecvBytes(mpi.MPI_ANY_SOURCE, mpi.MPI_ANY_TAG)
	log.Debugf("    %d: recv: received bytes %s from inbox", node.Source, string(recvbytes))
	tag := status.GetTag()
	msg := Msg{Bytes: recvbytes, Remote: status.GetSource(), Tag: tag}
	return msg
}
