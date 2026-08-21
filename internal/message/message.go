package message

type Message struct {
	Topic     string
	Key       []byte
	Value     []byte
	Partition int32
	Offset    int64
	ack       func()
}

func (m *Message) SetAck(ack func()) {
	m.ack = ack
}

func (m *Message) Ack() {
	if m.ack != nil {
		m.ack()
	}
}
