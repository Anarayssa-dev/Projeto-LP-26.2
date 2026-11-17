package icmp

import (
	"encoding/binary"
	"errors"
	"testing"
)

func TestPacketStructure(t *testing.T) {
	msg := Message{
		Header: ICMPHeader{Type: TypeEchoRequest, Code: 0, ID: 1234, Sequence: 1},
		Data:   []byte("PING-DATA"),
	}

	raw, err := Marshal(msg)
	if err != nil {
		t.Fatalf("Marshal() retornou erro inesperado: %v", err)
	}
	if raw[0] != TypeEchoRequest {
		t.Errorf("esperado tipo %d (Echo Request), recebido %d", TypeEchoRequest, raw[0])
	}

	_, err = ParseMessage(nil)
	if !errors.Is(err, ErrBufferTooShort) {
		t.Errorf("esperava ErrBufferTooShort para pacote vazio, recebeu %v", err)
	}
}

func TestMarshal(t *testing.T) {
	msg := Message{
		Header: ICMPHeader{
			Type:     TypeEchoRequest,
			Code:     0,
			ID:       0x1234,
			Sequence: 0x0001,
		},
		Data: []byte("HELLO"),
	}

	raw, err := Marshal(msg)
	if err != nil {
		t.Fatalf("Marshal retornou erro: %v", err)
	}
	if len(raw) != 13 {
		t.Errorf("tamanho incorreto: esperado 13, obtido %d", len(raw))
	}
	if raw[0] != TypeEchoRequest || raw[1] != 0 {
		t.Errorf("Type/Code incorretos: type=%d, code=%d", raw[0], raw[1])
	}
	if got := binary.BigEndian.Uint16(raw[4:6]); got != 0x1234 {
		t.Errorf("ID incorreto: esperado 0x1234, obtido 0x%x", got)
	}
	if got := CalculateChecksum(raw); got != 0 {
		t.Errorf("checksum inválido: verificação retornou 0x%x", got)
	}
}

func TestChecksumOddLength(t *testing.T) {
	data := []byte{0x08, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x41}
	if got := CalculateChecksum(data); got == 0 {
		t.Error("checksum calculado não deveria ser zero para esses dados")
	}
}
