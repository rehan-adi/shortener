package service

import "shortly-proto/gen/key"

type KeyServiceServer struct {
	key.UnimplementedKeyServiceServer
}

func NewKeyServiceServer() key.KeyServiceServer {
	return &KeyServiceServer{}
}

