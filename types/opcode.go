package types

type Opcode uint16

const (
	PingCode               Opcode = 0
	PongCode               Opcode = 1
	LookupNodeRequestCode  Opcode = 2
	LookupNodeResponseCode Opcode = 3
)

var (
	opcodeTable = map[Opcode]interface{}{
		PingCode:               Ping{},
		PongCode:               Pong{},
		LookupNodeRequestCode:  LookupNodeRequest{},
		LookupNodeResponseCode: LookupNodeResponse{},
	}
)

func AddMessageType(i interface{}) {

}

func GetMessageType(code Opcode) interface{} {
	if i, ok := opcodeTable[code]; ok {
		return i
	}
	return nil
}
