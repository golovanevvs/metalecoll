package grpchandler

import (
	"context"
	"encoding/json"

	pb "github.com/golovanevvs/metalecoll/internal/proto"
	"github.com/golovanevvs/metalecoll/internal/server/constants"
	"github.com/golovanevvs/metalecoll/internal/server/dto"
	"github.com/golovanevvs/metalecoll/internal/server/model"
)

// UpdateMetricsJSON - обновление метрики, полученной в JSON.
func (grpchd *gRPCHandler) UpdateMetricsJSON(ctx context.Context, in *pb.UpdateMetricsJSONRequest) (*pb.UpdateMetricsJSONResponse, error) {
	grpchd.lg.Debugf("Запуск gRPC handler UpdateMetricsJSON...")

	var response pb.UpdateMetricsJSONResponse

	decFromCTX, ok := ctx.Value(constants.DecryptKey).([]byte)
	if !ok {
		grpchd.lg.Errorf("Ошибка декодирования crypto\n")
		response.Success = false
		return &response, nil
	}

	var mVParse any
	var req dto.Metrics
	// var resp dto.Metrics

	grpchd.lg.Debugf("Декодирование JSON...")
	if err := json.Unmarshal(decFromCTX, &req); err != nil {
		grpchd.lg.Errorf("Ошибка декодирования JSON: %v", err)
		response.Success = false
		return &response, nil
	}
	grpchd.lg.Debugf("Декодирование JSON прошло успешно")

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

	// grpchd.lg.Debugf("Формирование тела ответа...")
	// switch req.MType {
	// case constants.GaugeType:
	// 	if v, ok := calcMetric.MetValue.(float64); ok {
	// 		resp = dto.Metrics{
	// 			ID:    calcMetric.MetName,
	// 			MType: calcMetric.MetType,
	// 			Value: &v,
	// 		}
	// 	}
	// case constants.CounterType:
	// 	if v, ok := calcMetric.MetValue.(int64); ok {
	// 		resp = dto.Metrics{
	// 			ID:    calcMetric.MetName,
	// 			MType: calcMetric.MetType,
	// 			Delta: &v,
	// 		}
	// 	}
	// }
	// grpchd.lg.Debugf("Формирование тела ответа прошло успешно")

	// grpchd.lg.Debugf("Кодирование в JSON...")
	// enc, err := json.Marshal(resp)
	// if err != nil {
	// 	grpchd.lg.Errorf("Ошибка кодирования: %v", err)
	// 	return
	// }
	// grpchd.lg.Debugf("Кодирование в JSON прошло успешно")

	// grpchd.lg.Debugf("Отправлен Content-Type: %v", constants.ContentTypeAJ)
	// w.Header().Set("Content-Type", constants.ContentTypeAJ)

	// grpchd.lg.Debugf("Отправлен код: %v", http.StatusOK)
	// w.WriteHeader(http.StatusOK)

	// w.Write(enc)

	response.Success = true

	return &response, nil
}
