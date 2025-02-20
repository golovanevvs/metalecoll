package grpchandler

import (
	"context"

	pb "github.com/golovanevvs/metalecoll/internal/proto"
	"github.com/golovanevvs/metalecoll/internal/server/constants"
	"github.com/golovanevvs/metalecoll/internal/server/dto"
	"github.com/golovanevvs/metalecoll/internal/server/model"
)

// UpdateMetricsJSON - обновление метрики
func (grpchd *gRPCHandler) UpdateMetrics(ctx context.Context, in *pb.UpdateMetricsRequest) (*pb.UpdateMetricsResponse, error) {
	grpchd.lg.Debugf("Запуск gRPC handler UpdateMetricsJSON...")

	var response pb.UpdateMetricsResponse

	var mVParse any
	var req dto.Metrics

	req.ID = in.Id
	req.MType = in.Type
	req.Value = &in.Value
	req.Delta = &in.Delta

	grpchd.lg.Debugf("Проверка наличия имени метрики...")
	if req.ID == "" {
		grpchd.lg.Errorf("Имя метрики не задано")
		response.Success = false
		return &response, nil
	}
	grpchd.lg.Debugf("Проверка наличия имени метрики прошла успешно")

	switch req.MType {
	case constants.GaugeType:
		mVParse = *req.Value
	case constants.CounterType:
		mVParse = *req.Delta
	default:
		grpchd.lg.Errorf("Неизвестный тип метрики")
		response.Success = false
		return &response, nil
	}

	receivedMetric := model.Metric{MetType: req.MType, MetName: req.ID, MetValue: mVParse}
	grpchd.lg.Debugf("Полученная метрика: %v", receivedMetric)

	grpchd.lg.Debugf("Запуск сервиса обновления метрики...")
	calcMetric := grpchd.sv.UpdateMetric(receivedMetric)
	grpchd.lg.Debugf("Обновление метрики прошло успешно, обновлённая метрика: %v", calcMetric)

	updatedMap := grpchd.sv.GetMetricsFromMap()
	grpchd.lg.Debugf("Обновлённая мапа: %v", updatedMap)

	response.Success = true

	return &response, nil
}
