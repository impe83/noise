package network

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"sync"

	"github.com/perlin-network/noise/log"
	"github.com/perlin-network/noise/types"

	"github.com/davecgh/go-xdr/xdr2"
	"github.com/pkg/errors"
)

var errEmptyMsg = errors.New("received an empty message from a peer")

// sendMessage marshals, signs and sends a message over a stream.
func (n *Network) sendMessage(w io.Writer, message *types.Message, writerMutex *sync.Mutex) error {
	// Serialize size.
	var msgBuf bytes.Buffer
	_, err := xdr.Marshal(&msgBuf, message)
	if err != nil {
		return errors.Wrap(err, "failed to write message to buffer")
	}

	// Serialize size.
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, uint32(msgBuf.Len()))

	buffer = append(buffer, msgBuf.Bytes()...)
	totalSize := len(buffer)

	log.Info().Interface("msg", message).Msgf("buffer: %+v", buffer)

	// Write until all bytes have been written.
	bytesWritten, totalBytesWritten := 0, 0

	writerMutex.Lock()

	bw, isBuffered := w.(*bufio.Writer)
	if isBuffered && (bw.Buffered() > 0) && (bw.Available() < totalSize) {
		if err := bw.Flush(); err != nil {
			return err
		}
	}

	for totalBytesWritten < len(buffer) && err == nil {
		bytesWritten, err = w.Write(buffer[totalBytesWritten:])
		if err != nil {
			log.Error().Err(err).Msg("stream: failed to write entire buffer")
		}
		totalBytesWritten += bytesWritten
	}

	writerMutex.Unlock()

	if err != nil {
		return errors.Wrap(err, "stream: failed to write to socket")
	}

	return nil
}

// receiveMessage reads, unmarshals and verifies a message from a net.Conn.
func (n *Network) receiveMessage(conn net.Conn) (*types.Message, error) {
	var err error

	// Read until all header bytes have been read.
	buffer := make([]byte, 4)

	bytesRead, totalBytesRead := 0, 0

	for totalBytesRead < 4 && err == nil {
		bytesRead, err = conn.Read(buffer[totalBytesRead:])
		totalBytesRead += bytesRead
	}

	// Decode message size.
	size := binary.BigEndian.Uint32(buffer)

	if size == 0 {
		return nil, errEmptyMsg
	}

	// Message size at most is limited to 4MB. If a big message need be sent,
	// consider partitioning to message into chunks of 4MB.
	if size > 4e+6 {
		return nil, errors.Errorf("message has length of %d which is either broken or too large", size)
	}

	// Read until all message bytes have been read.
	buffer = make([]byte, size)

	bytesRead, totalBytesRead = 0, 0

	for totalBytesRead < int(size) && err == nil {
		bytesRead, err = conn.Read(buffer[totalBytesRead:])
		totalBytesRead += bytesRead
	}

	// Deserialize message.
	msg := new(types.Message)

	_, err = xdr.Unmarshal(bytes.NewReader(buffer), msg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal message")
	}

	// Check if any of the message headers are invalid or null.
	if msg.Body == nil || &msg.Sender == nil || msg.Sender.PublicKey == nil || len(msg.Sender.Address) == 0 {
		log.Info().Msgf("%+v", msg)
		return nil, errors.New("received an invalid message (either no message or no sender) from a peer")
	}

	/*// Verify signature of message.
	if !crypto.Verify(
		n.opts.signaturePolicy,
		n.opts.hashPolicy,
		msg.Sender.PublicKey,
		SerializeMessage(msg.Sender, msg.Message.Value),
		msg.Signature,
	) {
		return nil, errors.New("received message had an malformed signature")
	}*/

	return msg, nil
}
