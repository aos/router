package main

import (
	"encoding/binary"
	"log"
	"net"

	"github.com/songgao/water"
)

// TAP interface works on L2 (read/write ethernet frames)
// TUN interface works on L3 (read/write IP packets)

// EtherType ...
type EtherType uint16

const (
	// EtherTypeARP ...
	EtherTypeARP EtherType = 0x0806
	// EtherTypeIPv4 ...
	EtherTypeIPv4 EtherType = 0x0800

	minPayload = 46
)

func main() {
	config := water.Config{
		DeviceType: water.TAP,
	}
	config.Name = "louie0"

	ifce, err := water.New(config)
	if err != nil {
		log.Fatal(err)
	}
	defer ifce.Close()

	_ = &Frame{
		Destination: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		Source:      []byte{0x3c, 0x97, 0x0e, 0x9e, 0x31, 0x96},
		EtherType:   EtherTypeARP,
		Payload:     []byte("hello world"),
	}

	buffer := make([]byte, 64)
	for {
		n, err := ifce.Read(buffer)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Read: ", buffer[:n])
	}

}

// Frame is an IEEE 802.3 Ethernet II frame (WIP: VLAN tagging)
type Frame struct {
	Destination net.HardwareAddr // 6 octets
	Source      net.HardwareAddr // 6 octets
	EtherType   EtherType        // 2 octets (no VLAN tagging)
	Payload     []byte           // minPayload + n
}

// MarshalBinary converts the frame into its binary format
func (f *Frame) MarshalBinary() ([]byte, error) {
	if len(f.Payload) < minPayload {
		padding := make([]byte, minPayload-len(f.Payload))
		f.Payload = append(f.Payload, padding...)
	}

	b := make([]byte, 14+len(f.Payload))
	copy(b[0:6], f.Destination)
	copy(b[6:12], f.Source)
	binary.BigEndian.PutUint16(b[12:14], uint16(f.EtherType))
	copy(b[14:], f.Payload)
	return b, nil
}
