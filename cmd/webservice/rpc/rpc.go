package rpc

import (
	insvc "github.com/nmluci/go-backend/internal/service"
	fRPC "github.com/nmluci/go-backend/pkg/rpc/fileop"
)

type FileRPC struct {
	fRPC.UnimplementedStellarFileServer
	service insvc.Service
}

func Init(svc insvc.Service) fRPC.StellarFileServer {
	return &FileRPC{
		service: svc,
	}
}
