package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golovanevvs/metalecoll/internal/server/constants"
	"github.com/golovanevvs/metalecoll/internal/server/model"
	"github.com/golovanevvs/metalecoll/internal/server/storage/mapstorage"
	"github.com/stretchr/testify/assert"
)

func TestMainHandler(t *testing.T) {
	store := mapstorage.NewStorage()
	srv = NewServer(store)

	type metCalc struct {
		metric   model.Metric
		contType string
	}

	type want struct {
		code        int
		contentType string
		metricCalc  model.Metric
	}

	tests := []struct {
		name string
		in   metCalc
		want want
	}{
		{
			name: "test №1 (positive)",
			in: metCalc{
				metric: model.Metric{
					MetType:  constants.GaugeType,
					MetName:  "Name1",
					MetValue: 5.3,
				},
				contType: constants.ContentType,
			},
			want: want{
				code:        200,
				contentType: constants.ContentType,
				metricCalc: model.Metric{
					MetType:  constants.GaugeType,
					MetName:  "Name1",
					MetValue: 5.3,
				},
			},
		},
		{
			name: "test №2 (positive)",
			in: metCalc{
				metric: model.Metric{
					MetType:  constants.GaugeType,
					MetName:  "Name2",
					MetValue: 100.12,
				},
				contType: constants.ContentType,
			},
			want: want{
				code:        200,
				contentType: constants.ContentType,
				metricCalc: model.Metric{
					MetType:  constants.GaugeType,
					MetName:  "Name2",
					MetValue: 100.12,
				},
			},
		},
		{
			name: "test №3 (positive)",
			in: metCalc{
				metric: model.Metric{
					MetType:  constants.CounterType,
					MetName:  "Name3",
					MetValue: 100,
				},
				contType: constants.ContentType,
			},
			want: want{
				code:        200,
				contentType: constants.ContentType,
				metricCalc: model.Metric{
					MetType:  constants.CounterType,
					MetName:  "Name3",
					MetValue: 100,
				},
			},
		},
		{
			name: "test №4 (positive)",
			in: metCalc{
				metric: model.Metric{
					MetType:  constants.CounterType,
					MetName:  "Name4",
					MetValue: 5,
				},
				contType: constants.ContentType,
			},
			want: want{
				code:        200,
				contentType: constants.ContentType,
				metricCalc: model.Metric{
					MetType:  constants.CounterType,
					MetName:  "Name4",
					MetValue: 105,
				},
			},
		},
		{
			name: "test №5 (negative)",
			in: metCalc{
				metric: model.Metric{
					MetType:  constants.GaugeType,
					MetName:  "Name5",
					MetValue: "строка",
				},
				contType: constants.ContentType,
			},
			want: want{
				code:        400,
				contentType: constants.ContentType,
				metricCalc: model.Metric{
					MetType:  constants.GaugeType,
					MetName:  "Name5",
					MetValue: 0,
				},
			},
		},
		{
			name: "test №6 (negative)",
			in: metCalc{
				metric: model.Metric{
					MetType:  constants.CounterType,
					MetName:  "Name6",
					MetValue: "строка",
				},
				contType: constants.ContentType,
			},
			want: want{
				code:        400,
				contentType: constants.ContentType,
				metricCalc: model.Metric{
					MetType:  constants.CounterType,
					MetName:  "Name6",
					MetValue: 0,
				},
			},
		},
		{
			name: "test №7 (negative)",
			in: metCalc{
				metric: model.Metric{
					MetType:  constants.CounterType,
					MetName:  "Name7",
					MetValue: 5.5,
				},
				contType: constants.ContentType,
			},
			want: want{
				code:        400,
				contentType: constants.ContentType,
				metricCalc: model.Metric{
					MetType:  constants.CounterType,
					MetName:  "Name7",
					MetValue: 0,
				},
			},
		},
		{
			name: "test №8 (negative)",
			in: metCalc{
				metric: model.Metric{
					MetType:  "Unknown",
					MetName:  "Name8",
					MetValue: 5.5,
				},
				contType: constants.ContentType,
			},
			want: want{
				code:        400,
				contentType: constants.ContentType,
				metricCalc: model.Metric{
					MetType:  constants.CounterType,
					MetName:  "Name8",
					MetValue: 0,
				},
			},
		},
		{
			name: "test №9 (negative)",
			in: metCalc{
				metric: model.Metric{
					MetType:  "Unknown",
					MetName:  "Name9",
					MetValue: 5.5,
				},
				contType: "application/json",
			},
			want: want{
				code:        400,
				contentType: constants.ContentType,
				metricCalc: model.Metric{
					MetType:  constants.CounterType,
					MetName:  "Name9",
					MetValue: 0,
				},
			},
		},
		{
			name: "test №10 (negative)",
			in: metCalc{
				metric: model.Metric{
					MetType:  constants.GaugeType,
					MetName:  "",
					MetValue: 5.5,
				},
				contType: "application/json",
			},
			want: want{
				code:        404,
				contentType: constants.ContentType,
				metricCalc: model.Metric{
					MetType:  constants.GaugeType,
					MetName:  "Name10",
					MetValue: 0,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			target := fmt.Sprintf("http://localhost:8080/update/%s/%s/%v", test.in.metric.MetType, test.in.metric.MetName, test.in.metric.MetValue)
			request := httptest.NewRequest(http.MethodPost, target, nil)
			request.Header.Set("Content-Type", constants.ContentType)
			w := httptest.NewRecorder()
			//MainHandler(w, request)
			srv.ServeHTTP(w, request)
			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))

		})
	}
}
