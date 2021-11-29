package messages

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	mpi "github.com/sbromberger/gompi"
)

const DEFAULT_TAG = 0
const TERMINATE_TAG = 1<<31 - 1

type Msg struct {
	Remote int
	Tag    int
	Bytes  []byte
}

func (m Msg) String() string {
	return fmt.Sprintf("Message: Remote %d, Tag %d, \"%s\"", m.Remote, m.Tag, string(m.Bytes))
}

type Node struct {
	Source        int
	Inbox, Outbox chan Msg
	comm          *mpi.Communicator
}

func NewNode(source int, comm *mpi.Communicator, bufsize int) *Node {
	inbox := make(chan Msg, bufsize)
	outbox := make(chan Msg, bufsize)
	node := Node{Source: source, comm: comm, Inbox: inbox, Outbox: outbox}
	return &node
}

func send(node *Node) {
	log.Warnf("    %d: starting up send goroutine", node.Source)
	defer close(node.Outbox)
	defer log.Warnf("    %d: closing down send goroutine", node.Source)
	for {
		msg := <-node.Outbox
		log.Warnf("    %d: send: sending msg %v from outbox", node.Source, msg)
		node.comm.SendBytes(msg.Bytes, msg.Remote, msg.Tag)
		if msg.Tag == TERMINATE_TAG {
			log.Warnf("    %d send: terminating", node.Source)
			return
		}
	}
}

func recv(node *Node) {
	log.Warnf("    %d: starting up recv goroutine", node.Source)
	defer close(node.Inbox)
	defer log.Warnf("    %d: closing down recv goroutine", node.Source)
	for {
		recvbytes, status := node.comm.RecvBytes(mpi.MPI_ANY_SOURCE, mpi.MPI_ANY_TAG)
		log.Warnf("    %d: recv: received bytes %s from inbox", node.Source, string(recvbytes))
		tag := status.GetTag()
		if tag == TERMINATE_TAG {
			log.Warnf("    %d recv: terminating", node.Source)
			return
		}
		msg := Msg{Bytes: recvbytes, Remote: status.GetSource(), Tag: tag}
		node.Inbox <- msg
	}
}

func (node *Node) Launch() {
	go send(node)
	go recv(node)
}
func (node *Node) Terminate() {
	log.Warnf("    %d: Terminate", node.Source)
	// node.comm.SendBytes([]byte{0}, node.Source, TERMINATE_TAG)
	node.Outbox <- Msg{Bytes: []byte{0}, Remote: node.Source, Tag: TERMINATE_TAG}
}
