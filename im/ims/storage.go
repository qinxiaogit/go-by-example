package main

type Storage struct {
	*StorageFile
	*PeerStorage
	*GroupStorage
}
