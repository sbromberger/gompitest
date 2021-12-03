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
	Source        int
	Inbox, Outbox chan *Msg
	comm          *mpi.Communicator
}

func NewNode(source int, comm *mpi.Communicator, bufsize int) *Node {
	inbox := make(chan *Msg, bufsize)
	outbox := make(chan *Msg, bufsize)
	node := Node{Source: source, comm: comm, Inbox: inbox, Outbox: outbox}
	return &node
}

func send(node *Node) {
	runtime.LockOSThread()
	log.Debugf("    %d: starting up send goroutine", node.Source)
	defer log.Debugf("    %d: closing down send goroutine", node.Source)
	for {
		msg := <-node.Outbox
		log.Debugf("    %d: send: sending msg %v from outbox", node.Source, msg)
		node.comm.SendBytes(msg.Bytes, msg.Remote, msg.Tag)
		if msg.Tag == node.comm.MaxTag {
			log.Debugf("    %d send: terminating", node.Source)
			return
		}
	}
}

func recv(node *Node) {
	runtime.LockOSThread()
	log.Debugf("    %d: starting up recv goroutine", node.Source)
	// defer close(node.Inbox)
	// defer close(node.Outbox)
	defer log.Debugf("    %d: closing down recv goroutine", node.Source)
	for {
		recvbytes, status := node.comm.MrecvBytes(mpi.MPI_ANY_SOURCE, mpi.MPI_ANY_TAG)
		log.Debugf("    %d: recv: received bytes %s from inbox", node.Source, string(recvbytes))
		tag := status.GetTag()
		msg := &Msg{Bytes: recvbytes, Remote: status.GetSource(), Tag: tag}
		node.Inbox <- msg
		if tag == node.comm.MaxTag {
			log.Debugf("    %d recv: terminating", node.Source)
			return
		}
	}
}

func (node *Node) Launch() {
	runtime.LockOSThread()
	go send(node)
	go recv(node)
}
func (node *Node) Terminate() {
	log.Debugf("    %d: Terminate", node.Source)
	// node.comm.SendBytes([]byte{0}, node.Source, node.comm.MaxTag)
	node.Outbox <- &Msg{Bytes: []byte{0}, Remote: node.Source, Tag: node.comm.MaxTag}
	<-node.Inbox
	close(node.Inbox)
	close(node.Outbox)
	log.Debugf("    %d: Terminate: all channels closed", node.Source)
}
