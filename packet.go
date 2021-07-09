package gorcon

import "math/rand"

const (
	ServerdataAuth            = 3
	ServerdataAuthResponse    = 2
	ServerdataExeccommand     = 2
	SERVERDATA_RESPONSE_VALUE = 0

	MaxPacketSize = 4096
)

type packet struct {
	packetId int32
	packetType int32
	packetBody []byte
}

func newPacket(packetType int32, packetBody string) *packet {
	return &packet {
		packetId: rand.Int31(),
		packetType: packetType,
		packetBody: []byte(packetBody),
	}
}

func (packet *packet) calculateSize() int32 {
	return 4 + 4 + int32(len(packet.packetBody)) + 2
}