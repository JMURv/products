package dbutils

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestFilterItems(t *testing.T) {
	tests := []struct {
		name         string
		filters      map[string]any
		expected     string
		expectedArgs []any
	}{
		{
			name:         "Single min_price filter",
			filters:      map[string]any{"min_price": 100},
			expected:     " AND i.price >= ?",
			expectedArgs: []any{100},
		},
		{
			name:         "Single max_price filter",
			filters:      map[string]any{"max_price": 500},
			expected:     " AND i.price <= ?",
			expectedArgs: []any{500},
		},
		{
			name:         "Both min_price and max_price filters",
			filters:      map[string]any{"min_price": 100, "max_price": 500},
			expected:     " AND i.price >= ? AND i.price <= ?",
			expectedArgs: []any{100, 500},
		},
		{
			name:         "Range filter with min and max",
			filters:      map[string]any{"memory": map[string]any{"min": "64GB", "max": "256GB"}},
			expected:     " AND (ia.name = ? AND ia.value >= ?) AND (ia.name = ? AND ia.value <= ?)",
			expectedArgs: []any{"memory", "64GB", "memory", "256GB"},
		},
		{
			name:         "List filter with multiple values",
			filters:      map[string]any{"brand": []string{"Apple", "Samsung"}},
			expected:     " AND (ia.name = ? AND ia.value IN (?, ?))",
			expectedArgs: []any{"brand", "Apple", "Samsung"},
		},
		{
			name: "Combined filters",
			filters: map[string]any{
				"min_price": 100,
				"max_price": 500,
				"brand":     []string{"Apple", "Samsung"},
				"memory":    map[string]any{"min": "64GB", "max": "256GB"},
			},
			expected:     " AND i.price >= ? AND i.price <= ? AND (ia.name = ? AND ia.value IN (?, ?)) AND (ia.name = ? AND ia.value >= ?) AND (ia.name = ? AND ia.value <= ?)",
			expectedArgs: []any{100, 500, "brand", "Apple", "Samsung", "memory", "64GB", "memory", "256GB"},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				var q strings.Builder
				args := make([]any, 0)
				args = FilterItems(&q, args, tt.filters)

				require.Equal(t, tt.expected, q.String())
				require.Equal(t, tt.expectedArgs, args)
			},
		)
	}
}
