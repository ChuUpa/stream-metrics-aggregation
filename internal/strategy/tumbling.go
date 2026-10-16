package strategy

import (
	"sync"
	"time"

	"github.com/ChuUpa/stream-metrics-aggregation/internal/model"
)

// Tumbling — стратегия фиксированных непересекающихся окон (tumbling windows).
// Каждая точка попадает ровно в одно окно длины size.
type Tumbling struct {
	size time.Duration
	mu   sync.Mutex
	// TODO: добавь состояние:
	//   - открытые окна: агрегаты по ключу (начало окна, ключ ряда);
	//   - водяной знак: момент последнего Flush, чтобы распознавать опоздавшие точки.
}

// NewTumbling создаёт стратегию с окнами длины size.
func NewTumbling(size time.Duration) *Tumbling {
	return &Tumbling{size: size}
}

// Name возвращает название стратегии.
func (t *Tumbling) Name() string { return "tumbling" }

// Add учитывает точку в её окне.
func (t *Tumbling) Add(p model.Point) error {
	// TODO
	panic("not implemented")
}

// Flush возвращает агрегаты окон, у которых End <= now.
func (t *Tumbling) Flush(now time.Time) []Result {
	// TODO
	panic("not implemented")
}

// Проверка на этапе компиляции, что Tumbling реализует Aggregator.
var _ Aggregator = (*Tumbling)(nil)
