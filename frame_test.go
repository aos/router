package main

import (
	"testing"
)

func TestMarshalBinary(t *testing.T) {
	// Technically they should be atleast 64 bytes long including
	// the 4-byte CRC at the end
	t.Run("Frames should be at least 60 bytes long", func(t *testing.T) {
		f := &Frame{
			Destination: []byte{0x6a, 0xd2, 0x58, 0x62, 0xa5, 0xeb},
			Source:      []byte{0x0a, 0xae, 0x42, 0xdb, 0xf5, 0x48},
			EtherType:   EtherTypeARP,
			Payload:     []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64},
		}
		b, err := f.MarshalBinary()
		if err != nil {
			t.Error(err)
		}
		if len(b) < 60 {
			t.Errorf("got: %d, want: 60", len(b))
		}
	})
}
