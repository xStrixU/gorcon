package gorcon

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

// https://developer.valvesoftware.com/wiki/Source_RCON_Protocol

type Rcon struct {
	conn net.Conn
}

func Connect(address string, password string) (*Rcon, error) {
	conn, err := net.Dial("tcp", address)

	if err != nil {
		return nil, err
	}

	rcon := &Rcon {
		conn: conn,
	}

	if err = rcon.auth(password); err != nil {
		return nil, err
	}

	return rcon, nil
}

func (rcon *Rcon) sendPacket(packet *packet) (*packet, error) {
	buf := bytes.Buffer{}

	_ = binary.Write(&buf, binary.LittleEndian, packet.calculateSize())
	_ = binary.Write(&buf, binary.LittleEndian, packet.packetId)
	_ = binary.Write(&buf, binary.LittleEndian, packet.packetType)
	_ = binary.Write(&buf, binary.LittleEndian, packet.packetBody)
	_ = binary.Write(&buf, binary.LittleEndian, make([]byte, 2))

	if buf.Len() > MaxPacketSize {
		return nil, fmt.Errorf("Packet is too big: (expected maximum value of %d, got %d)", MaxPacketSize, buf.Len())
	}

	if _, err := rcon.conn.Write(buf.Bytes()); err != nil {
		return nil, err
	}

	if response, err := rcon.readPacket(); err != nil {
		return nil, err
	} else {
		return response, nil
	}
}

func(rcon *Rcon) readPacket() (*packet, error) {
	var packetSize int32
	var packetId int32
	var packetType int32

	if err := binary.Read(rcon.conn, binary.LittleEndian, &packetSize); err != nil {
		return nil, fmt.Errorf("Failed to read response packet size: %w", err)
	}

	if err := binary.Read(rcon.conn, binary.LittleEndian, &packetId); err != nil {
		return nil, fmt.Errorf("Failed to read response packet id: %w", err)
	}

	if err := binary.Read(rcon.conn, binary.LittleEndian, &packetType); err != nil {
		return nil, fmt.Errorf("Failed to read response packet type: %w", err)
	}

	packetBody := make([]byte, packetSize - 8)

	if err := binary.Read(rcon.conn, binary.LittleEndian, &packetBody); err != nil {
		return nil, fmt.Errorf("Failed to read response packet body: %w", err)
	}

	return &packet {
		packetId: packetId,
		packetType: packetType,
		packetBody: packetBody[:len(packetBody)-2],
	}, nil
}

func (rcon *Rcon) auth(password string) error {
	packet := newPacket(ServerdataAuth, password)
	response, err := rcon.sendPacket(packet)

	if err != nil {
		return err
	}

	if response.packetType == ServerdataAuthResponse && response.packetId == -1 {
		return errors.New("Failed to authorize!")
	}

	return nil
}

func (rcon *Rcon) SendCommand(command string) (string, error) {
	packet := newPacket(ServerdataExeccommand, command)
	response, err := rcon.sendPacket(packet)

	if err != nil {
		return "", err
	}

	return string(response.packetBody), nil
}

func (rcon *Rcon) Close() error {
	return rcon.conn.Close()
}