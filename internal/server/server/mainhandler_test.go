package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golovanevvs/metalecoll/internal/server/constants"
	"github.com/golovanevvs/metalecoll/internal/server/model"
	"github.com/stretchr/testify/assert"
)

func TestMainHandler(t *testing.T) {
	type want struct {
		code        int
		contentType string
		metricCalc  model.Metric
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive test №1",
			want: want{
				code:        200,
				contentType: constants.ContentType,
				metricCalc:  model.Metric{},
			},
		},
	}
	type metCalc struct {
		metric    model.Metric
		CalcValue any
	}

	metrics := make([]metCalc, 7)

	metrics[0] = metCalc{
		metric: model.Metric{
			MetType:  constants.GaugeType,
			MetName:  "Name1",
			MetValue: 5.3,
		},
		CalcValue: 5.3,
	}
	metrics[1] = metCalc{
		metric: model.Metric{
			MetType:  constants.GaugeType,
			MetName:  "Name2",
			MetValue: 100.12,
		},
		CalcValue: 100.12,
	}
	metrics[2] = metCalc{
		metric: model.Metric{
			MetType:  constants.CounterType,
			MetName:  "Name3",
			MetValue: 100,
		},
		CalcValue: 101,
	}
	metrics[3] = metCalc{
		metric: model.Metric{
			MetType:  constants.CounterType,
			MetName:  "Name4",
			MetValue: 301,
		},
		CalcValue: 302,
	}
	metrics[4] = metCalc{
		metric: model.Metric{
			MetType:  constants.GaugeType,
			MetName:  "Name5",
			MetValue: "строка",
		},
		CalcValue: 0,
	}
	metrics[5] = metCalc{
		metric: model.Metric{
			MetType:  constants.CounterType,
			MetName:  "Name6",
			MetValue: "строка",
		},
		CalcValue: 0,
	}
	metrics[6] = metCalc{
		metric: model.Metric{
			MetType:  constants.CounterType,
			MetName:  "Name7",
			MetValue: 5.5,
		},
		CalcValue: 0,
	}

	target := make([]string, 7)
	for i := range metrics {
		target[i] = fmt.Sprintf("http://localhost:8080/update/%s/%s/%v", metrics[i].metric.MetType, metrics[i].metric.MetName, metrics[i].metric.MetValue)
	}

	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			switch i {
			case 0:
				request := httptest.NewRequest(http.MethodPost, target[i], nil)
				request.Header.Set("Content-Type", constants.ContentType)
				w := httptest.NewRecorder()
				MainHandler(w, request)
				res := w.Result()
				assert.Equal(t, test.want.code, res.StatusCode)
				assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}
