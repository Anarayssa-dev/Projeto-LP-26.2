package icmp

import (
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

// TODO: transformar a Message em uma fatia de bytes para o envio na rede. Logica Checksum tambem

func Marshal(msg Message) ([]byte, error) {
	return nil, nil
}

// TODO: Transformar uma fatia de bytes em Message

func ParseMessage(b []byte) (*Message, error) {
	return nil, nil
}
