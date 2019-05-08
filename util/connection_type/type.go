package connection_type

type ConnectionType uint8

const (
	TCP  ConnectionType = 0
	QUIC ConnectionType = 1
)
