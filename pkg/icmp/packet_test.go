package icmp

import (
	"encoding/binary"
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

	_, err := Marshal(msg)
	if err != nil {
		t.Fatalf("Marshal retornou erro inesperado: %v", err)
	}

	_, err = ParseMessage([]byte{})
	if err != nil && err != ErrBufferTooShort {
	}
}

func TestMarshal(t *testing.T) {
	msg := Message{
		Header: ICMPHeader{
			Type:     TypeEchoRequest,
			Code:     0,
			Checksum: 0,
			ID:       0x1234,
			Sequence: 0x0001,
		},
		Data: []byte("HELLO"),
	}

	raw, err := Marshal(msg)
	if err != nil {
		t.Fatalf("Erro ao executar Marshal: %v", err)
	}

	if len(raw) != 13 {
		t.Errorf("Tamanho do buffer incorreto. Esperado 13, obtido %d", len(raw))
	}

	if raw[0] != TypeEchoRequest || raw[1] != 0 {
		t.Errorf("Header Type/Code incorretos: type=%d, code=%d", raw[0], raw[1])
	}

	idParsed := binary.BigEndian.Uint16(raw[4:6])
	if idParsed != 0x1234 {
		t.Errorf("ID incorreto: esperado 0x1234, obtido 0x%x", idParsed)
	}

	verifyChk := CalculateChecksum(raw)
	if verifyChk != 0 {
		t.Errorf("Checksum inválido! A verificação do pacote gerado retornou 0x%x (esperado 0x0000)", verifyChk)
	}
}

func TestChecksumOddLength(t *testing.T) {
	// Teste de integridade com payload de tamanho ímpar
	data := []byte{0x08, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x41} // 9 bytes
	chk := CalculateChecksum(data)
	if chk == 0 {
		t.Errorf("Checksum calculado não pode ser zero para dados não nulos")
	}
}
