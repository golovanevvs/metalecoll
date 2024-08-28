package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golovanevvs/metalecoll/internal/server/constants"
	"github.com/golovanevvs/metalecoll/internal/server/model"
	"github.com/golovanevvs/metalecoll/internal/server/storage"
	"github.com/golovanevvs/metalecoll/internal/server/storage/mapstorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateMetricsHandler(t *testing.T) {
	configtest := &Config{
		Addr:           constants.AddrS,
		GaugeType:      constants.GaugeType,
		CounterType:    constants.CounterType,
		UpdateMethod:   constants.UpdateMethod,
		GetValueMethod: constants.GetValueMethod,
	}
	store := mapstorage.NewStorage()
	srv = NewServer(store, configtest)

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
					MetValue: float64(5.3),
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
					MetValue: float64(100.12),
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
					MetValue: int64(100),
				},
			},
		},
		{
			name: "test №4 (positive)",
			in: metCalc{
				metric: model.Metric{
					MetType:  constants.CounterType,
					MetName:  "Name3",
					MetValue: 5,
				},
				contType: constants.ContentType,
			},
			want: want{
				code:        200,
				contentType: constants.ContentType,
				metricCalc: model.Metric{
					MetType:  constants.CounterType,
					MetName:  "Name3",
					MetValue: int64(105),
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
					MetValue: float64(100.12),
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
					MetValue: int64(105),
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
					MetValue: int64(105),
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
					MetType:  "Unknown",
					MetName:  "Name8",
					MetValue: 0,
				},
			},
		},
		{
			name: "test №9 (negative)",
			in: metCalc{
				metric: model.Metric{
					MetType:  constants.GaugeType,
					MetName:  "Name9",
					MetValue: 5.5,
				},
				contType: "application/json",
			},
			want: want{
				code:        400,
				contentType: constants.ContentType,
				metricCalc: model.Metric{
					MetType:  constants.GaugeType,
					MetName:  "Name9",
					MetValue: float64(100.12),
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
				contType: constants.ContentType,
			},
			want: want{
				code:        404,
				contentType: constants.ContentType,
				metricCalc: model.Metric{
					MetType:  constants.GaugeType,
					MetName:  "Name10",
					MetValue: float64(100.12),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			target := fmt.Sprintf("http://%s/%s/%s/%s/%v", constants.AddrA, constants.UpdateMethod, test.in.metric.MetType, test.in.metric.MetName, test.in.metric.MetValue)
			request := httptest.NewRequest(http.MethodPost, target, nil)
			request.Header.Set("Content-Type", test.in.contType)
			w := httptest.NewRecorder()
			srv.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			switch test.name {
			case "test №1 (positive)", "test №2 (positive)", "test №3 (positive)", "test №4 (positive)":
				assert.Equal(t, test.want.code, res.StatusCode)
				assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
				v, err := storage.GM(store, test.want.metricCalc.MetName)
				require.NoError(t, err)
				assert.Equal(t, test.want.metricCalc.MetValue, v.MetValue)
			case "test №8 (negative)":
				assert.Equal(t, test.want.code, res.StatusCode)
				_, err := storage.GM(store, test.want.metricCalc.MetName)
				assert.Error(t, err)
			default:
				assert.Equal(t, test.want.code, res.StatusCode)
			}
		})
	}
}
