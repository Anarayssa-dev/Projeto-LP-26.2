package icmp

import (
	"encoding/binary"
	"errors"
)

// Tipos de mensagem que o ICMP pode enviar no IPV4
const (
	TypeEchoReply   uint8 = 0
	TypeEchoRequest uint8 = 8
)

// Estrutura do protocolo ICMP sem o payload

type ICMPHeader struct {
	Type     uint8  // Tipo da mensagem ou Echo Request ou Echo Reply
	Code     uint8  // Código da mensagem
	Checksum uint16 // Auto explicativo
	ID       uint16 // Identificador
	Sequence uint16 // Numero da sequencia do pacote
}

// Agrupa o protocolo ICMP com o payload

type Message struct {
	Header ICMPHeader
	Data   []byte
}

var (
	ErrBufferTooShort = errors.New("buffer muito curto para conter um cabeçalho ICMP válido")
	ErrInvalidType    = errors.New("tipo de pacote ICMP não suportado")
)

func CalculateChecksum(data []byte) uint16 {
	var sum uint32
	length := len(data)

	// Soma palavras 2 bytes por vez
	for i := 0; i < length-1; i += 2 {
		sum += uint32(binary.BigEndian.Uint16(data[i : i+2]))
	}

	// Se for ímpar, adiciona o ultimo byte restante ajustado a esquerda
	if length%2 != 0 {
		sum += uint32(data[length-1]) << 8
	}

	// Dobra a soma de 32 bits em 16 bits
	for (sum >> 16) > 0 {
		sum = (sum & 0xffff) + (sum >> 16)
	}

	// Retorna negado
	return uint16(^sum)
}

func Marshal(msg Message) ([]byte, error) {
	// Calcula o tamanho da mensagem 8 bytes do cabeçalho + o payload
	totalLen := 8 + len(msg.Data)
	buf := make([]byte, totalLen)

	// Preenche os campos fixos do cabeçalho
	buf[0] = msg.Header.Type
	buf[1] = msg.Header.Code

	// Checksum iniciado em 0
	binary.BigEndian.PutUint16(buf[2:4], 0)

	// Preenche o ID e Sequence do protocolo ICMP
	binary.BigEndian.PutUint16(buf[4:6], msg.Header.ID)
	binary.BigEndian.PutUint16(buf[6:8], msg.Header.Sequence)

	// Copia o payload para depois dos 8 bytes
	copy(buf[8:], msg.Data)

	// Calcula o checksum com o campo do checksum zerado
	chk := CalculateChecksum(buf)

	// Escreve o checksum calculado nos bytes 2 e 3 do buffer
	binary.BigEndian.PutUint16(buf[2:4], chk)

	return buf, nil
}

// TODO: Transformar uma fatia de bytes em Message

func ParseMessage(b []byte) (*Message, error) {
	return nil, nil
}
