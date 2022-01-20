package types

type Action int
type SignalType int
type Unit struct{}

const (
	AsyncSignal = iota
	SyncSignal  = iota
)

const (
	PENDING  = iota
	CLAIMED  = iota
	CONSUMED = iota
)

const (
	MESSAGE        = iota
	AddJoinPattern = iota
	GetNewPortId   = iota
	Shutdown       = iota
)

type Packet struct {
	SignalId Port
	Type     Action
	Payload  Payload
}

type Payload struct {
	Msg    interface{}
	Ch     chan interface{}
	status int
}

type MessageChannel struct {
	Ch chan Payload
}

type Port struct {
	ChannelType     SignalType
	Id              int
	JunctionChannel chan Packet
}

type JoinPatternPacket struct {
	Signals []Port
	Action  interface{}
}

type WrappedJoinPattern struct {
	Pattern JoinPatternPacket
	Bitmask int
}

type UnaryAsync = func(interface{})
type UnarySync = func(interface{}) interface{}

type BinaryAsync = func(interface{}, interface{})
type BinarySync = func(interface{}, interface{}) interface{}

type TernaryAsync = func(interface{}, interface{}, interface{})
type TernarySync = func(interface{}, interface{}, interface{}) interface{}
