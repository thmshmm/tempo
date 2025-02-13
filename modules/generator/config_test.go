package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grafana/tempo/modules/generator/processor/servicegraphs"
	"github.com/grafana/tempo/modules/generator/processor/spanmetrics"
)

func TestProcessorConfig_copyWithOverrides(t *testing.T) {
	original := &ProcessorConfig{
		ServiceGraphs: servicegraphs.Config{
			HistogramBuckets: []float64{1},
			Dimensions:       []string{},
		},
		SpanMetrics: spanmetrics.Config{
			HistogramBuckets:    []float64{1, 2},
			Dimensions:          []string{"namespace"},
			IntrinsicDimensions: spanmetrics.IntrinsicDimensions{Service: true},
		},
	}

	t.Run("overrides buckets and dimension", func(t *testing.T) {
		o := &mockOverrides{
			serviceGraphsHistogramBuckets:  []float64{1, 2},
			serviceGraphsDimensions:        []string{"namespace"},
			spanMetricsHistogramBuckets:    []float64{1, 2, 3},
			spanMetricsDimensions:          []string{"cluster", "namespace"},
			spanMetricsIntrinsicDimensions: map[string]bool{"status_code": true},
		}

		copied, err := original.copyWithOverrides(o, "tenant")
		require.NoError(t, err)

		assert.NotEqual(t, *original, copied)

		// assert nothing changed
		assert.Equal(t, []float64{1}, original.ServiceGraphs.HistogramBuckets)
		assert.Equal(t, []string{}, original.ServiceGraphs.Dimensions)
		assert.Equal(t, []float64{1, 2}, original.SpanMetrics.HistogramBuckets)
		assert.Equal(t, []string{"namespace"}, original.SpanMetrics.Dimensions)
		assert.Equal(t, spanmetrics.IntrinsicDimensions{Service: true}, original.SpanMetrics.IntrinsicDimensions)

		// assert overrides were applied
		assert.Equal(t, []float64{1, 2}, copied.ServiceGraphs.HistogramBuckets)
		assert.Equal(t, []string{"namespace"}, copied.ServiceGraphs.Dimensions)
		assert.Equal(t, []float64{1, 2, 3}, copied.SpanMetrics.HistogramBuckets)
		assert.Equal(t, []string{"cluster", "namespace"}, copied.SpanMetrics.Dimensions)
		assert.Equal(t, spanmetrics.IntrinsicDimensions{Service: true, StatusCode: true}, copied.SpanMetrics.IntrinsicDimensions)
	})

	t.Run("empty overrides", func(t *testing.T) {
		o := &mockOverrides{}

		copied, err := original.copyWithOverrides(o, "tenant")
		require.NoError(t, err)

		assert.Equal(t, *original, copied)
	})

	t.Run("invalid overrides", func(t *testing.T) {
		o := &mockOverrides{
			spanMetricsIntrinsicDimensions: map[string]bool{"invalid": true},
		}

		_, err := original.copyWithOverrides(o, "tenant")
		require.Error(t, err)
	})
}
