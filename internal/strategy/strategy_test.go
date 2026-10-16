package strategy

import (
	"errors"
	"math"
	"testing"
	"time"

	"github.com/ChuUpa/stream-metrics-aggregation/internal/model"
)

var t0 = time.Date(2026, 10, 16, 12, 0, 0, 0, time.UTC)

func point(name string, value float64, offset time.Duration) model.Point {
	return model.Point{
		Name:      name,
		Labels:    model.Labels{"host": "a"},
		Value:     value,
		Timestamp: t0.Add(offset),
	}
}

func TestResult_Merge(t *testing.T) {
	a := Result{Count: 2, Sum: 3, Min: 1, Max: 2}
	b := Result{Count: 1, Sum: 10, Min: 10, Max: 10}

	got := a.Merge(b)
	want := Result{Count: 3, Sum: 13, Min: 1, Max: 10}
	if got != want {
		t.Errorf("a.Merge(b) = %+v, want %+v", got, want)
	}
	if b.Merge(a) != want {
		t.Errorf("Merge не коммутативен: b.Merge(a) = %+v", b.Merge(a))
	}
}

func TestResult_MergeIdentity(t *testing.T) {
	x := Result{Count: 1, Sum: 5, Min: 5, Max: 5}
	if got := (Result{}).Merge(x); got != x {
		t.Errorf("Result{}.Merge(x) = %+v, want %+v", got, x)
	}
	if got := x.Merge(Result{}); got != x {
		t.Errorf("x.Merge(Result{}) = %+v, want %+v", got, x)
	}
}

func TestResult_Avg(t *testing.T) {
	if got := (Result{Count: 4, Sum: 10}).Avg(); got != 2.5 {
		t.Errorf("Avg = %v, want 2.5", got)
	}
	if got := (Result{}).Avg(); !math.IsNaN(got) {
		t.Errorf("Avg пустого агрегата = %v, want NaN", got)
	}
}

func TestTumbling_SameWindow(t *testing.T) {
	agg := NewTumbling(10 * time.Second)
	for i, v := range []float64{1, 2, 3} {
		if err := agg.Add(point("cpu", v, time.Duration(i+1)*time.Second)); err != nil {
			t.Fatalf("Add: неожиданная ошибка: %v", err)
		}
	}

	results := agg.Flush(t0.Add(10 * time.Second))
	if len(results) != 1 {
		t.Fatalf("ожидался 1 агрегат, получено %d", len(results))
	}
	r := results[0]
	if r.Count != 3 || r.Sum != 6 || r.Min != 1 || r.Max != 3 {
		t.Errorf("неверный агрегат: %+v", r)
	}
	if !r.Window.Start.Equal(t0) || !r.Window.End.Equal(t0.Add(10*time.Second)) {
		t.Errorf("неверное окно: %+v", r.Window)
	}
	if want := model.SeriesKey("cpu", model.Labels{"host": "a"}); r.SeriesKey != want {
		t.Errorf("SeriesKey = %q, want %q", r.SeriesKey, want)
	}
}

func TestTumbling_FlushBeforeWindowEnd(t *testing.T) {
	agg := NewTumbling(10 * time.Second)
	if err := agg.Add(point("cpu", 1, time.Second)); err != nil {
		t.Fatal(err)
	}
	if got := agg.Flush(t0.Add(5 * time.Second)); len(got) != 0 {
		t.Fatalf("окно ещё не закрыто, но получено %d агрегатов", len(got))
	}
	if got := agg.Flush(t0.Add(10 * time.Second)); len(got) != 1 {
		t.Fatalf("ожидался 1 агрегат после закрытия окна, получено %d", len(got))
	}
	if got := agg.Flush(t0.Add(20 * time.Second)); len(got) != 0 {
		t.Fatalf("окно уже выгружено, повторно получено %d агрегатов", len(got))
	}
}

func TestTumbling_TwoWindows(t *testing.T) {
	agg := NewTumbling(10 * time.Second)
	_ = agg.Add(point("cpu", 1, 1*time.Second))
	_ = agg.Add(point("cpu", 5, 11*time.Second))

	results := agg.Flush(t0.Add(20 * time.Second))
	if len(results) != 2 {
		t.Fatalf("ожидалось 2 агрегата, получено %d", len(results))
	}
	sums := map[time.Duration]float64{}
	for _, r := range results {
		sums[r.Window.Start.Sub(t0)] = r.Sum
	}
	if sums[0] != 1 || sums[10*time.Second] != 5 {
		t.Errorf("неверное распределение по окнам: %v", sums)
	}
}

func TestTumbling_SeriesSeparated(t *testing.T) {
	agg := NewTumbling(10 * time.Second)
	_ = agg.Add(point("cpu", 1, time.Second))
	_ = agg.Add(point("mem", 2, time.Second))

	if got := agg.Flush(t0.Add(10 * time.Second)); len(got) != 2 {
		t.Fatalf("разные ряды должны агрегироваться отдельно, получено %d агрегатов", len(got))
	}
}

func TestTumbling_LatePoint(t *testing.T) {
	agg := NewTumbling(10 * time.Second)
	_ = agg.Add(point("cpu", 1, time.Second))
	agg.Flush(t0.Add(10 * time.Second))

	err := agg.Add(point("cpu", 2, 2*time.Second))
	if !errors.Is(err, ErrLatePoint) {
		t.Errorf("ожидалась ErrLatePoint, получено: %v", err)
	}
	if err := agg.Add(point("cpu", 3, 12*time.Second)); err != nil {
		t.Errorf("точка в открытом окне не должна быть опоздавшей: %v", err)
	}
}
