package types

import (
	"bytes"

	"github.com/davecgh/go-xdr/xdr2"
)

// Message defines a Noise message
type Message struct {
	// Body is the serialized message body
	Body []byte
	// MessageNonce is the sequence ID.
	MessageNonce uint64
	// Opcode defines the Opcode of the message
	Opcode Opcode
	// ReplyFlag indicates this is a reply to a request
	ReplyFlag bool
	// RequestNonce is the request/response ID. Null if ID associated with a message is not a request/response.
	RequestNonce uint64
	// Sender's address and public key.
	Sender *ID
	// Sender's signature of message.
	Signature []byte
}

func (m *Message) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	_, err := xdr.Marshal(&buf, m)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type Ping struct {
}

type Pong struct {
}

type LookupNodeRequest struct {
	target *ID
}

type LookupNodeResponse struct {
	peers []*ID
}
