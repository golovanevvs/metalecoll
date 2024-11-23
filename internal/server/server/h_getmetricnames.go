package server

import (
	"html/template"
	"net/http"

	"github.com/golovanevvs/metalecoll/internal/server/constants"
)

func (s *server) GetMetricNamesHandler(w http.ResponseWriter, r *http.Request) {
	srv.logger.Debugf("Получение всех известных метрик из хранилища...")
	metricsMap := s.mapStore.GetMetrics()
	srv.logger.Debugf("Получение всех известных метрик из хранилища прошло успешно")

	srv.logger.Debugf("Отправка тела ответа...")
	tmpl := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Метрики</title>
		</head>
		<body>
			<h1>Список метрик</h1>
			<table border="1">
				<tr>
					<th>Имя метрики</th>
					<th>Тип метрики</th>
					<th>Значение</th>
				</tr>
				{{range .}}
				<tr>
					<td>{{.MetName}}</td>
					<td>{{.MetType}}</td>
					<td>{{.MetValue}}</td>
				</tr>
				{{end}}
			</table>
		</body>
		</html>
	`

	t, err := template.New("metricArr").Parse(tmpl)
	if err != nil {
		srv.logger.Errorf("Ошибка создания шаблона: %v", err)
		srv.logger.Debugf("Отправлен код: %v", http.StatusInternalServerError)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", constants.ContentTypeTHUTF8)
	if err := t.Execute(w, metricsMap); err != nil {
		srv.logger.Errorf("Ошибка применения шаблона: %v", err)
		srv.logger.Debugf("Отправлен код: %v", http.StatusInternalServerError)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	srv.logger.Debugf("Отправка тела ответа прошла успешно")
}
