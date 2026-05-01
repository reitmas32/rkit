package criteria

import (
	"testing"
)

func TestCriteria(t *testing.T) {
	filters := NewFilters([]Filter{
		{Field: "name", Operator: OperatorEqual, Value: "test"},
		{Field: "age", Operator: OperatorGreaterThan, Value: 18},
	})

	criteria := Criteria{
		Filters: *filters,
	}

	if criteria.Filters.Get() == nil {
		t.Error("Criteria should have filters")
	}

	if len(criteria.Filters.Get()) != 2 {
		t.Errorf("Criteria should have 2 filters, got %d", len(criteria.Filters.Get()))
	}
}

func TestCriteria_Empty(t *testing.T) {
	criteria := Criteria{
		Filters: *NewFilters([]Filter{}),
	}

	if len(criteria.Filters.Get()) != 0 {
		t.Errorf("Empty Criteria should have 0 filters, got %d", len(criteria.Filters.Get()))
	}
}
