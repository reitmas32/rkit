package criteria

import (
	"testing"
)

func TestFilter(t *testing.T) {
	filter := Filter{
		Field:    "name",
		Operator: OperatorEqual,
		Value:    "test",
	}

	if filter.Field != "name" {
		t.Errorf("Filter.Field = %v, want name", filter.Field)
	}
	if filter.Operator != OperatorEqual {
		t.Errorf("Filter.Operator = %v, want =", filter.Operator)
	}
	if filter.Value != "test" {
		t.Errorf("Filter.Value = %v, want test", filter.Value)
	}
}

func TestOperator_Constants(t *testing.T) {
	tests := []struct {
		name     string
		operator Operator
		expected string
	}{
		{"Equal", OperatorEqual, "="},
		{"NotEqual", OperatorNotEqual, "<>"},
		{"GreaterThan", OperatorGreaterThan, ">"},
		{"GreaterEqual", OperatorGreaterEqual, ">="},
		{"LessThan", OperatorLessThan, "<"},
		{"LessEqual", OperatorLessEqual, "<="},
		{"Like", OperatorLike, "LIKE"},
		{"NotLike", OperatorNotLike, "NOT LIKE"},
		{"In", OperatorIn, "IN"},
		{"NotIn", OperatorNotIn, "NOT IN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.operator) != tt.expected {
				t.Errorf("Operator %s = %v, want %v", tt.name, tt.operator, tt.expected)
			}
		})
	}
}

func TestFilter_WithIntValue(t *testing.T) {
	filter := Filter{
		Field:    "age",
		Operator: OperatorGreaterThan,
		Value:    18,
	}

	if filter.Value != 18 {
		t.Errorf("Filter.Value = %v, want 18", filter.Value)
	}
}

func TestFilter_WithStringValue(t *testing.T) {
	filter := Filter{
		Field:    "name",
		Operator: OperatorLike,
		Value:    "%test%",
	}

	if filter.Value != "%test%" {
		t.Errorf("Filter.Value = %v, want %%test%%", filter.Value)
	}
}
