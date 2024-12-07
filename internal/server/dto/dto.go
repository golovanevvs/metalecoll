// Модуль dto предназначен для хранения структур, необходимых для передачи данных между объектами.
package dto

// Metrics - структура для передачи данных между объектами.
type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

// MetricsPost - структура для передачи данных между объектами.
type MetricsPost struct {
	Name  string  `json:"name"`
	Type  string  `json:"type"`
	Value float64 `json:"value"`
}
