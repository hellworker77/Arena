package udp_types

import "net"

type InputPacket struct {
	Data []byte
	From *net.UDPAddr
}
