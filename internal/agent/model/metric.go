// Модуль model содержит основные модели данных.
package model

// Metric - модель метрик.
type Metric struct {
	Type  string
	Name  string
	Value any
}
