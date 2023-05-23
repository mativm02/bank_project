package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/mativm02/bank_system/token"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	authorizationBearer = "bearer"
)

func (server *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing context metadata")
	}

	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) != 2 {
		return nil, fmt.Errorf("invalid authorization header")
	}

	authType := fields[0]
	if strings.ToLower(authType) != authorizationBearer {
		return nil, fmt.Errorf("invalid authorization type")
	}

	tokenString := fields[1]
	payload, err := server.tokenMaker.VerifyToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return payload, nil
}
