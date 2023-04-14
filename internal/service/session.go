package service

import (
	"context"

	"github.com/nmluci/go-backend/internal/commonkey"
	"github.com/nmluci/go-backend/internal/indto"
	"github.com/nmluci/go-backend/internal/util/ctxutil"
	"github.com/nmluci/go-backend/internal/util/rpcutil"
)

var (
	tagLoggerAuthenticateSession = "[AuthenticateSession]"
	tagLoggerAuthenticateService = "[AuthenticateService]"
)

func (s *service) AuthenticateSession(ctx context.Context, token string) (access context.Context, err error) {
	ctx = rpcutil.AppendMetaContext(ctx)

	// conf := config.Get()
	// usr, err := s.stellarRPC.Auth.AuthorizeToken(ctx, &rpc.UserAccess{
	// 	AccessToken: token,
	// 	Requester:   conf.ServiceID,
	// })
	// if err != nil {
	// 	s.logger.Errorf("%s stellarRPC error: %+v", tagLoggerAuthenticateSession, err)
	// 	return
	// }

	scopeMap := indto.UserScopeMap{}
	// for _, scope := range usr.UserScope {

	// }

	access = ctxutil.WrapCtx(ctx, commonkey.SCOPE_CTX_KEY, scopeMap)
	return
}

func (s *service) AuthenticateService(ctx context.Context, name string) (access context.Context, err error) {
	ctx = rpcutil.AppendMetaContext(ctx)

	// conf := config.Get()
	// svcMeta, err := s.stellarRPC.Auth.AuthorizeService(ctx, &rpc.ServiceAccess{
	// 	ServiceName: name,
	// 	Requester:   conf.ServiceID,
	// })
	// if err != nil {
	// 	s.logger.Errorf("%s stellarRPC error: %+v", tagLoggerAuthenticateService, err)
	// 	return
	// }

	scopeMap := indto.UserScopeMap{}
	// for _, scope := range svcMeta.ServiceScope {

	// }

	access = ctxutil.WrapCtx(ctx, commonkey.SCOPE_CTX_KEY, scopeMap)
	return
}
