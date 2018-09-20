package types

import (
	"github.com/perlin-network/noise/peer"
)

type Message struct {
	Body []byte
	// MessageNonce is the sequence ID.
	MessageNonce uint64
	// ReplyFlag indicates this is a reply to a request
	ReplyFlag bool
	// RequestNonce is the request/response ID. Null if ID associated with a message is not a request/response.
	RequestNonce uint64
	// Sender's address and public key.
	Sender *peer.ID
	// Sender's signature of message.
	Signature []byte
}
