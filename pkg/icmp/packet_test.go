package icmp

import (
	"testing"
)

func TestPacketStructure(t *testing.T) {
	msg := Message{
		Header: ICMPHeader{
			Type:     TypeEchoRequest,
			Code:     0,
			Checksum: 0,
			ID:       1234,
			Sequence: 1,
		},
		Data: []byte("PING-DATA"),
	}

	if msg.Header.Type != 8 {
		t.Errorf("Esperado tipo 8 (Echo Request), recebido %d", msg.Header.Type)
	}

	// Valida se as assinaturas das funções estão funcionando
	_, err := Marshal(msg)
	if err != nil {
		t.Fatalf("Marshal retornou erro inesperado: %v", err)
	}

	_, err = ParseMessage([]byte{})
	if err != nil && err != ErrBufferTooShort {
	}
}
