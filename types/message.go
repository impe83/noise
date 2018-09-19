package types

type ID struct {
	// Address is the network address of the peer
	Address string
	// ID is the computed hash of the public key
	ID []byte
	// PublicKey of the peer (we no longer use the public key as the peer ID, but use it to verify messages)
	PublicKey []byte
}

type Message struct {
	Body []byte
	// MessageNonce is the sequence ID.
	MessageNonce uint64
	// ReplyFlag indicates this is a reply to a request
	ReplyFlag bool
	// RequestNonce is the request/response ID. Null if ID associated with a message is not a request/response.
	RequestNonce uint64
	// Sender's address and public key.
	Sender ID
	// Sender's signature of message.
	Signature []byte
}
