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

// runAggregation is the pipeline's aggregation stage. It reads every
// record off in, keeps a running sum/count per group as they arrive,
// and once the channel closes it computes the final results and sends
// them (along with every record it saw) on the returned channel.
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

// getValueIgnoringSpace looks up a field by name, and if there's no
// exact match, tries again ignoring leading/trailing whitespace on the
// keys - handy for CSV headers with stray spaces.
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

// accumulate folds one record into the running sums/counts for every
// configured aggregation, keyed by group.
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

// groupKeyFor builds the map key a record's value gets accumulated
// under: "<op>:<field>" normally, or "<op>:<field>:<groupBy value>" when
// the aggregation has a group-by field.
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

// buildResults turns the accumulated sums/counts into final Result rows,
// resolving each group's operation (sum, avg, or count) from its key.
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

// containsOp reports whether a group key starts with the given
// operation name (group keys are built as "<op>:...").
func containsOp(key, op string) bool {
	return len(key) >= len(op) && key[:len(op)] == op
}

// toFloat tries to coerce a record's field value (which comes in as
// float64, int, or a string from JSON/CSV) into a float64 for math.
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
