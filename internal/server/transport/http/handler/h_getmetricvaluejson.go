package handler

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"strconv"

// 	"github.com/golovanevvs/metalecoll/internal/server/constants"
// 	"github.com/golovanevvs/metalecoll/internal/server/dto"
// )

// func (s *server) GetMetricValueJSONHandler(w http.ResponseWriter, r *http.Request) {
// 	srv.logger.Debugf("Декодирование JSON...")

// 	req := dto.Metrics{}
// 	dec := json.NewDecoder(r.Body)
// 	if err := dec.Decode(&req); err != nil {
// 		srv.logger.Errorf("Ошибка декодирования JSON: %v", err)
// 		srv.logger.Errorf("Отправлен код: %v", http.StatusInternalServerError)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	defer r.Body.Close()
// 	srv.logger.Debugf("Декодирование JSON прошло успешно: %v", req)

// 	srv.logger.Debugf("Получение данных из хранилища по name %v...", req.ID)
// 	metric, err := s.mapStore.GetMetric(req.ID)
// 	if err != nil {
// 		fmt.Println(err)
// 		srv.logger.Errorf("Ошибка получения данных из хранилища: %v", err)
// 		srv.logger.Errorf("Отправлен код: %v", http.StatusNotFound)
// 		w.WriteHeader(http.StatusNotFound)
// 		return
// 	}
// 	srv.logger.Debugf("Получение данных из хранилища прошло успешно: %v", metric)

// 	srv.logger.Debugf("Формирование тела ответа...")
// 	resp := dto.Metrics{}
// 	switch metric.MetType {
// 	case constants.GaugeType:
// 		v := metric.MetValue.(float64)
// 		resp = dto.Metrics{
// 			ID:    metric.MetName,
// 			MType: metric.MetType,
// 			Value: &v,
// 		}
// 	case constants.CounterType:
// 		if d, ok := metric.MetValue.(int64); ok {
// 			srv.logger.Debugf("d: %v", d)
// 			resp = dto.Metrics{
// 				ID:    metric.MetName,
// 				MType: metric.MetType,
// 				Delta: &d,
// 			}
// 		} else {
// 			d := metric.MetValue.(float64)
// 			d1 := fmt.Sprintf("%.0f", d)
// 			d2, _ := strconv.Atoi(d1)
// 			d3 := int64(d2)
// 			resp = dto.Metrics{
// 				ID:    metric.MetName,
// 				MType: metric.MetType,
// 				Delta: &d3,
// 			}
// 		}
// 	}
// 	srv.logger.Debugf("Формирование тела ответа прошло успешно: %v", resp)

// 	srv.logger.Debugf("Отправлен Content-Type: %v", constants.ContentTypeAJ)
// 	w.Header().Set("Content-Type", constants.ContentTypeAJ)

// 	srv.logger.Debugf("Кодирование в JSON...")
// 	enc, err := json.Marshal(resp)
// 	if err != nil {
// 		srv.logger.Errorf("Ошибка кодирования: %v", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	srv.logger.Debugf("Кодирование в JSON прошло успешно")

// 	srv.logger.Debugf("Отправлен код: %v", http.StatusOK)
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(enc)
// }
