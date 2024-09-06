package server

import (
	"html/template"
	"net/http"

	"github.com/golovanevvs/metalecoll/internal/server/constants"
	"github.com/golovanevvs/metalecoll/internal/server/storage"
)

func GetMetricNamesHandler(w http.ResponseWriter, r *http.Request, store storage.Storage) {
	srv.logger.Debugf("Получение всех известных метрик из хранилища...")
	metricsMap := storage.GMs(store)
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
					<td>{{.Name}}</td>
					<td>{{.Type}}</td>
					<td>{{.Value}}</td>
				</tr>
				{{end}}
			</table>
		</body>
		</html>
	`

	t, err := template.New("metrics").Parse(tmpl)
	if err != nil {
		srv.logger.Errorf("Ошибка создания шаблона: %v", err)
		srv.logger.Debugf("Отправлен код: %v", http.StatusInternalServerError)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", constants.ContentTypeTPUTF8)
	if err := t.Execute(w, metricsMap); err != nil {
		srv.logger.Errorf("Ошибка применения шаблона: %v", err)
		srv.logger.Debugf("Отправлен код: %v", http.StatusInternalServerError)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	srv.logger.Debugf("Отправка тела ответа прошла успешно")

	srv.logger.Debugf("Отправляем код: %v", http.StatusOK)
	w.WriteHeader(http.StatusOK)
}
