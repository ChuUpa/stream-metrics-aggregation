// Package strategy содержит стратегии потоковой агрегации метрик.
package strategy

import (
	"errors"
	"time"

	"github.com/ChuUpa/stream-metrics-aggregation/internal/model"
)

// ErrLatePoint возвращается, если точка пришла после закрытия её окна.
var ErrLatePoint = errors.New("strategy: точка пришла после закрытия её окна")

// Window — полуинтервал времени [Start, End).
type Window struct {
	Start time.Time
	End   time.Time
}

// Result — агрегат одного временного ряда за одно окно.
type Result struct {
	SeriesKey string
	Window    Window
	Count     uint64
	Sum       float64
	Min       float64
	Max       float64
}

// Merge объединяет два частичных агрегата одного ряда за одно окно.
//
// Требования (свойства коммутативного моноида):
//   - ассоциативность и коммутативность: агрегаты с разных агентов
//     можно объединять в любом порядке;
//   - нейтральный элемент: пустой Result{} (Count == 0) ничего не меняет.
func (r Result) Merge(other Result) Result {
	// TODO: реализуй сам
	panic("not implemented")
}

// Avg возвращает среднее. Для пустого агрегата среднее не определено: NaN.
// Среднее — алгебраический агрегат: оно не объединяется напрямую,
// но вычисляется из объединяемой пары (Sum, Count).
func (r Result) Avg() float64 {
	// TODO: реализуй сам
	panic("not implemented")
}

// Aggregator — общий интерфейс стратегии агрегации.
// Любая стратегия, реализующая его, подключается к системе без изменения
// остальных компонентов.
type Aggregator interface {
	// Name возвращает название стратегии для логов и отчётов.
	Name() string
	// Add учитывает одну точку.
	Add(p model.Point) error
	// Flush возвращает агрегаты всех окон, закрытых к моменту now,
	// и освобождает занятую ими память.
	Flush(now time.Time) []Result
}
