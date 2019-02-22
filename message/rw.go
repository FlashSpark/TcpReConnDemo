package message

import (
	"TcpReConnDemo/rlp"
	"bytes"
	"errors"
	"io"
)

const (
	maxUint24 = ^uint32(0) >> 8
)

var (
	// this is used in place of actual frame header data.
	// TODO: replace this when Msg contains the protocol type code.
	zeroHeader = []byte{0xC2, 0x80, 0x80}
	// sixteen zero bytes
	zero16 = make([]byte, 16)
)

//NOTE: it's without any crypto
func DataRWIns(conn io.ReadWriter) *DataRW {
	return &DataRW{
		conn: conn,
	}
}

// data rw
type DataRW struct {
	conn io.ReadWriter
}

// implement of MsgReader
// bits[0:23] is the total size of the frame
func (rw *DataRW) ReadMsg() (msg Msg, err error) {
	// read the header
	headBuf := make([]byte, 32)
	if _, err := io.ReadFull(rw.conn, headBuf); err != nil {
		return msg, err
	}
	// NOTE fSize is the sum of len(message.code) and len(message.payload)
	fSize := readInt24(headBuf)

	// read the frame content
	var rSize = fSize // frame size rounded up to 16 byte boundary
	if padding := fSize % 16; padding > 0 {
		rSize += 16 - padding
	}

	// NOTE: read data from frame
	frameBuf := make([]byte, rSize)
	if _, err := io.ReadFull(rw.conn, frameBuf); err != nil {
		return msg, err
	}

	// get fSize bytes needed
	// decode message code
	content := bytes.NewReader(frameBuf[:fSize])
	if err := rlp.Decode(content, &msg.Code); err != nil {
		return msg, err
	}
	msg.Size = uint32(content.Len())
	msg.Payload = content

	return msg, nil
}

// implement of MsgWriter
func (rw *DataRW) WriteMsg(msg Msg) error {
	// calculate length of message
	pType, _ := rlp.EncodeToBytes(msg.Code)
	payload := rw.payload(msg)
	msg.Size = uint32(len(payload))
	// write header
	headBuf := make([]byte, 32)
	fSize := uint32(len(pType)) + msg.Size
	// it is not allowed if size if larger than 24bits
	if fSize > maxUint24 {
		return errors.New("message size overflows uint24")
	}
	putInt24(fSize, headBuf)
	copy(headBuf[3:], zeroHeader)

	// write headBuf which contains fSize and a bunch of following zeros
	if _, err := rw.conn.Write(headBuf); err != nil {
		return err
	}

	// write type code
	if _, err := rw.conn.Write(pType); err != nil {
		return err
	}
	// write payload
	if _, err := rw.conn.Write(payload); err != nil {
		return err
	}
	// fill bytes up
	if padding := fSize % 16; padding > 0 {
		if _, err := rw.conn.Write(zero16[:16-padding]); err != nil {
			return err
		}
	}
	return nil
}

// it will cost more cpu time
// but can save some memory
func (rw *DataRW) payload(msg Msg) []byte {
	buffer := bytes.NewBuffer(make([]byte, 0, 65536))
	_, err := io.Copy(buffer, msg.Payload)
	if err != nil {
		return nil
	}
	temp := buffer.Bytes()
	length := len(temp)
	var body []byte
	//are we wasting more than 5% space?
	if cap(temp) > (length + length/5) {
		body = make([]byte, length)
		copy(body, temp)
	} else {
		body = temp
	}
	return body
}

func readInt24(b []byte) uint32 {
	return uint32(b[2]) | uint32(b[1])<<8 | uint32(b[0])<<16
}

// save by Big-Endian
func putInt24(v uint32, b []byte) {
	b[0] = byte(v >> 16)
	b[1] = byte(v >> 8)
	b[2] = byte(v)
}
