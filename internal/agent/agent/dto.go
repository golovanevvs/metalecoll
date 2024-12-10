// Модуль dto предназначен для хранения структур, необходимых для передачи данных между объектами.
package agent

// Metrics - структура для передачи данных между объектами.
type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}
