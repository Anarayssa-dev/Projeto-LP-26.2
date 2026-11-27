package icmp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"testing"
)

func TestParseMessageRoundTrip(t *testing.T) {
	want := Message{
		Header: ICMPHeader{Type: TypeEchoRequest, Code: 0, ID: 0x1234, Sequence: 7},
		Data:   []byte("teste ICMP"),
	}

	raw, err := Marshal(want)
	if err != nil {
		t.Fatalf("Marshal() erro: %v", err)
	}

	got, err := ParseMessage(raw)
	if err != nil {
		t.Fatalf("ParseMessage() erro: %v", err)
	}
	if got.Header.Type != want.Header.Type || got.Header.Code != want.Header.Code ||
		got.Header.ID != want.Header.ID || got.Header.Sequence != want.Header.Sequence {
		t.Fatalf("cabeçalho diferente: obtido %+v, esperado %+v", got.Header, want.Header)
	}
	if !bytes.Equal(got.Data, want.Data) {
		t.Fatalf("payload diferente: %q", got.Data)
	}
	if got.Header.Checksum != binary.BigEndian.Uint16(raw[2:4]) {
		t.Fatal("checksum não foi preservado")
	}
}

func TestParseMessageBufferTooShort(t *testing.T) {
	_, err := ParseMessage([]byte{TypeEchoRequest, 0, 0})
	if !errors.Is(err, ErrBufferTooShort) {
		t.Fatalf("esperava ErrBufferTooShort, recebeu %v", err)
	}
}

func TestParseMessageInvalidType(t *testing.T) {
	packet := make([]byte, 8)
	packet[0] = 3 // Destination Unreachable: não é Echo Request/Reply nesta implementação.
	_, err := ParseMessage(packet)
	if !errors.Is(err, ErrInvalidType) {
		t.Fatalf("esperava ErrInvalidType, recebeu %v", err)
	}
}

func TestParseMessageEchoReply(t *testing.T) {
	packet := make([]byte, 8)
	packet[0] = TypeEchoReply
	msg, err := ParseMessage(packet)
	if err != nil {
		t.Fatalf("ParseMessage() erro: %v", err)
	}
	if msg.Header.Type != TypeEchoReply {
		t.Fatalf("tipo obtido %d; esperado %d", msg.Header.Type, TypeEchoReply)
	}
}
