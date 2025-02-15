package grpchandler

import (
	"github.com/golovanevvs/metalecoll/internal/proto"
	"github.com/golovanevvs/metalecoll/internal/server/service"
	"github.com/sirupsen/logrus"
)

type gRPCHandler struct {
	sv             *service.Service
	lg             *logrus.Logger
	hashKey        string
	privateKeyPath string
	trustedSubnet  string
	proto.UnimplementedMetricsServer
}

func NewGrpcHandler(sv *service.Service, lg *logrus.Logger, hashKey string, privateKeyPath string, trustedSubnet string) *gRPCHandler {
	return &gRPCHandler{
		sv:             sv,
		lg:             lg,
		hashKey:        hashKey,
		privateKeyPath: privateKeyPath,
		trustedSubnet:  trustedSubnet,
	}
}
