package network

import (
	"encoding/binary"
	"net"

	"github.com/perlin-network/noise/internal/protobuf"
)

// SerializeMessage compactly packs all bytes of a message together for cryptographic signing purposes.
func SerializeMessage(id *types.ID, message []byte) []byte {
	const uint32Size = 4

	serialized := make([]byte, uint32Size+len(id.Address)+uint32Size+len(id.ID)+len(message))
	pos := 0

	binary.LittleEndian.PutUint32(serialized[pos:], uint32(len(id.Address)))
	pos += uint32Size

	copy(serialized[pos:], []byte(id.Address))
	pos += len(id.Address)

	binary.LittleEndian.PutUint32(serialized[pos:], uint32(len(id.ID)))
	pos += uint32Size

	copy(serialized[pos:], id.ID)
	pos += len(id.ID)

	copy(serialized[pos:], message)
	pos += len(message)

	if pos != len(serialized) {
		panic("internal error: invalid serialization output")
	}

	return serialized
}

// FilterPeers filters out duplicate/empty addresses.
func FilterPeers(address string, peers []string) (filtered []string) {
	visited := make(map[string]struct{})
	visited[address] = struct{}{}

	for _, peerAddress := range peers {
		if len(peerAddress) == 0 {
			continue
		}

		resolved, err := ToUnifiedAddress(peerAddress)
		if err != nil {
			continue
		}
		if _, exists := visited[resolved]; !exists {
			filtered = append(filtered, resolved)
			visited[resolved] = struct{}{}
		}
	}
	return filtered
}

// GetRandomUnusedPort returns a random unused port
func GetRandomUnusedPort() int {
	listener, _ := net.Listen("tcp", ":0")
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port
}

// ProtoPeerToPeer converts a protobuf.ID type to types.ID
func ProtoPeerToPeer(p *protobuf.ID) *types.ID {
	if p == nil {
		return nil
	}
	return &types.ID{
		Address:   p.Address,
		ID:        p.Id,
		PublicKey: p.PublicKey,
	}
}

// PeerToProtoPeer converts a types.ID to protobuf.ID
func PeerToProtoPeer(p *types.ID) *protobuf.ID {
	if p == nil {
		return nil
	}
	return &protobuf.ID{
		Address:   p.Address,
		Id:        p.ID,
		PublicKey: p.PublicKey,
	}
}
