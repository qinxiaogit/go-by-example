package main

type ServerSummary struct {
	nrequests int64
	peer_message_count int64
	group_message_count int64
}

func NewServerSummary()  *ServerSummary{
	return new(ServerSummary)
}