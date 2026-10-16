// Package model содержит базовые типы данных системы.
package model

import (
	"maps"
	"slices"
	"strings"
	"time"
)

// Labels — набор меток временного ряда, например {"host": "a", "service": "api"}.
type Labels map[string]string

// Point — одно измерение метрики в момент времени.
type Point struct {
	Name      string
	Labels    Labels
	Value     float64
	Timestamp time.Time
}

// SeriesKey строит однозначный ключ временного ряда: имя + отсортированные метки.
// Пример: cpu_usage{host=a,service=api}.
//
// Сортировка обязательна: порядок обхода map в Go случаен,
// и без неё один и тот же ряд получал бы разные ключи.
func SeriesKey(name string, labels Labels) string {
	var b strings.Builder
	b.WriteString(name)
	b.WriteByte('{')
	for i, k := range slices.Sorted(maps.Keys(labels)) {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(k)
		b.WriteByte('=')
		b.WriteString(labels[k])
	}
	b.WriteByte('}')
	return b.String()
}
