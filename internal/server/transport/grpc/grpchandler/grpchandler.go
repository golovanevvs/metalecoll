package grpchandler

import (
	pb "github.com/golovanevvs/metalecoll/internal/proto"
	"github.com/golovanevvs/metalecoll/internal/server/service"
	"github.com/sirupsen/logrus"
)

type gRPCHandler struct {
	sv *service.Service
	lg *logrus.Logger
	pb.UnimplementedMetricsServer
}

func NewGrpcHandler(sv *service.Service, lg *logrus.Logger) *gRPCHandler {
	return &gRPCHandler{
		sv: sv,
		lg: lg,
	}
}
