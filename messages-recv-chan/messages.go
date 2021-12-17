// messages implements inbox via a channel. Send is not a channel.
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
	Bytes  *([]byte)
}

func (m Msg) String() string {
	return fmt.Sprintf("Message: Remote %d, Tag %d, \"%s\"", m.Remote, m.Tag, string(*m.Bytes))
}

type Node struct {
	Source              int
	Inbox               chan *Msg
	comm                *mpi.Communicator
	localGlobalMsgCount uint64
}

func NewNode(source int, comm *mpi.Communicator, bufsize int) *Node {
	inbox := make(chan *Msg, bufsize)
	node := Node{Source: source, comm: comm, Inbox: inbox}
	return &node
}

func (node *Node) Send(msg *Msg) {
	node.comm.SendBytes(*msg.Bytes, msg.Remote, msg.Tag)
	node.localGlobalMsgCount++
}

func recv(node *Node) {
	runtime.LockOSThread()
	log.Debugf("    %d: starting up recv goroutine", node.Source)
	defer log.Debugf("    %d: closing down recv goroutine", node.Source)
	for {
		recvbytes, status := node.comm.MrecvBytes(mpi.AnySource, mpi.AnyTag)
		log.Debugf("    %d: recv: received bytes %s from inbox", node.Source, string(recvbytes))
		tag := status.GetTag()
		msg := &Msg{Bytes: &recvbytes, Remote: status.GetSource(), Tag: tag}
		node.Inbox <- msg
		node.localGlobalMsgCount--

		if tag == node.comm.MaxTag {
			log.Debugf("    %d recv: terminating", node.Source)
			return
		}
	}
}

func (node *Node) Launch() {
	runtime.LockOSThread()
	go recv(node)
}
func (node *Node) Terminate() {
	log.Debugf("    %d: Terminate", node.Source)
	// node.comm.SendBytes([]byte{0}, node.Source, node.comm.MaxTag)
	node.Send(&Msg{Bytes: &([]byte{0}), Remote: node.Source, Tag: node.comm.MaxTag})
	<-node.Inbox
	close(node.Inbox)
	log.Debugf("    %d: Terminate: all channels closed", node.Source)
}
