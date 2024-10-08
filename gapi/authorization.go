package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/yudanl96/revive/token"
	"google.golang.org/grpc/metadata"
)

const (
	authorizzationHeader = "authorization"
	authorizationBearer  = "bearer"
)

func (server *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {
	mtd, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing meta data")
	}

	values := mtd.Get(authorizzationHeader)
	if len(values) == 0 {
		return nil, fmt.Errorf("missing authotization header")
	}

	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authotization header format")
	}

	authType := strings.ToLower(fields[0])
	if authType != authorizationBearer {
		return nil, fmt.Errorf("unsupported authorization")
	}

	token := fields[1]
	payload, err := server.tokenMaker.VerifyToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %s", err)
	}

	return payload, nil
}
