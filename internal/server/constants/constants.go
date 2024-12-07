// Модуль constants содержит в себе все константы приложений server и agent.
package constants

const (
	GaugeType   = "gauge"
	CounterType = "counter"

	AddrS = ":8080"
	AddrA = "localhost:8080"

	ContentTypeTPUTF8 = "text/plain; charset=utf-8"
	ContentTypeTP     = "text/plain"
	ContentTypeAJ     = "application/json"
	ContentTypeTH     = "text/html"
	ContentTypeTHUTF8 = "text/html; charset=utf-8"

	UpdateMethod   = "update"
	GetValueMethod = "value"
)
