package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGateawayUserAgentHeader = "grpcgateway-user-agent"
	userAgentHeader             = "user-agent"
	xForwardedForHeader         = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}

	if metadata, ok := metadata.FromIncomingContext(ctx); ok {
		if userAgents := metadata.Get(grpcGateawayUserAgentHeader); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}

		if userAgents := metadata.Get(userAgentHeader); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}

		if clientIPs := metadata.Get(xForwardedForHeader); len(clientIPs) > 0 {
			mtdt.ClientIP = clientIPs[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		mtdt.ClientIP = p.Addr.String()
	}

	return mtdt
}
