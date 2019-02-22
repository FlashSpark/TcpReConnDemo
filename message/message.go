package message

import (
	"TcpReConnDemo/rlp"
	"fmt"
	"io"
	"io/ioutil"
	"time"
)

type Msg struct {
	Code       uint64
	Size       uint32 // size of the payload
	Payload    io.Reader
	ReceivedAt time.Time
}

func (msg Msg) String() string {
	return fmt.Sprintf("message #%v (%v bytes)", msg.Code, msg.Size)
}

// Discard reads any remaining payload data into a black hole.
func (msg Msg) Discard() error {
	_, err := io.Copy(ioutil.Discard, msg.Payload)
	return err
}

type MsgReader interface {
	ReadMsg() (Msg, error)
}

type MsgWriter interface {
	// WriteMsg sends a message. It will block until the message's
	// Payload has been consumed by the other end.
	//
	// Note that messages can be sent only once because their
	// payload reader is drained.
	WriteMsg(Msg) error
}

// MsgReadWriter provides reading and writing of encoded messages.
// Implementations should ensure that ReadMsg and WriteMsg can be
// called simultaneously from multiple goroutines.
type MsgReadWriter interface {
	MsgReader
	MsgWriter
}

// Send writes an RLP-encoded message with the given code.
// data should encode as an RLP list.
func Send(w MsgWriter, code uint64, data interface{}) error {
	size, r, err := rlp.EncodeToReader(data)
	if err != nil {
		return err
	}
	// fmt.Println("Msg send code:", msgcode, " size:", size)
	return w.WriteMsg(Msg{Code: code, Size: uint32(size), Payload: r})
}
