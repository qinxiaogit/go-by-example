package common

type OpenMessage struct {
	Addr string
}

type OpenAckMessage struct {
	Status bool
}

type RelayMessage struct {
	Data []byte
}
