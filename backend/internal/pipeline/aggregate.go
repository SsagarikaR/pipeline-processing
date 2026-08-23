package pipeline

import (
	"context"
	"fmt"
	"strings"

	"github.com/SsagarikaR/pipeline-processing/internal/models"
)

type aggOutput struct {
	Records []Record
	Results []models.Result
}

func runAggregation(ctx context.Context, jobID int, in <-chan Record, configs []AggregationConfig, progressCh chan<- struct{}) <-chan aggOutput {
	outCh := make(chan aggOutput, 1)

	go func() {
		defer close(outCh)
		records := make([]Record, 0, 1024)
		sums := map[string]float64{}
		counts := map[string]int{}

		for {
			select {
			case <-ctx.Done():
				return
			case r, ok := <-in:
				if !ok {
					outCh <- aggOutput{Records: records, Results: buildResults(jobID, sums, counts)}
					return
				}
				records = append(records, r)
				accumulate(r, configs, sums, counts)
				progressCh <- struct{}{}
			}
		}
	}()

	return outCh
}

func getValueIgnoringSpace(data map[string]any, target string) (any, bool) {
	if v, ok := data[target]; ok {
		return v, true
	}
	t := strings.TrimSpace(target)
	for k, v := range data {
		if strings.TrimSpace(k) == t {
			return v, true
		}
	}
	return nil, false
}

func accumulate(r Record, configs []AggregationConfig, sums map[string]float64, counts map[string]int) {
	for _, c := range configs {
		key := groupKeyFor(c, r)
		counts[key]++
		if v, ok := getValueIgnoringSpace(r.Data, c.Field); ok {
			if f, ok := toFloat(v); ok {
				sums[key] += f
			}
		}
	}
}

func groupKeyFor(c AggregationConfig, r Record) string {
	base := c.Op + ":" + c.Field
	if c.GroupBy == "" {
		return base
	}
	if gv, ok := getValueIgnoringSpace(r.Data, c.GroupBy); ok {
		return fmt.Sprintf("%s:%v", base, gv)
	}
	return base
}

func buildResults(jobID int, sums map[string]float64, counts map[string]int) []models.Result {
	var results []models.Result
	for key, count := range counts {
		var val float64
		switch {
		case containsOp(key, "avg"):
			if count > 0 {
				val = sums[key] / float64(count)
			}
		case containsOp(key, "count"):
			val = float64(count)
		default:
			val = sums[key]
		}
		results = append(results, models.Result{JobID: jobID, GroupKey: key, AggregatedValue: val})
	}
	return results
}

func containsOp(key, op string) bool {
	return len(key) >= len(op) && key[:len(op)] == op
}

func toFloat(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	case string:
		var f float64
		_, err := fmt.Sscanf(n, "%f", &f)
		return f, err == nil
	default:
		return 0, false
	}
}
