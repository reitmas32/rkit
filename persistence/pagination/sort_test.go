package pagination

import (
	"testing"
)

func TestNewSort(t *testing.T) {
	tests := []struct {
		name      string
		field     string
		direction SortDirection
		expected  SortDirection
	}{
		{"ASC direction", "name", SortDirectionASC, SortDirectionASC},
		{"DESC direction", "name", SortDirectionDESC, SortDirectionDESC},
		{"invalid direction defaults to ASC", "name", SortDirection("INVALID"), SortDirectionASC},
		{"empty direction defaults to ASC", "name", SortDirection(""), SortDirectionASC},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSort(tt.field, tt.direction)
			if got == nil {
				t.Fatal("NewSort() returned nil")
			}
			if got.Field != tt.field {
				t.Errorf("NewSort().Field = %v, want %v", got.Field, tt.field)
			}
			if got.Direction != tt.expected {
				t.Errorf("NewSort().Direction = %v, want %v", got.Direction, tt.expected)
			}
		})
	}
}

func TestNewSortASC(t *testing.T) {
	sort := NewSortASC("name")

	if sort == nil {
		t.Fatal("NewSortASC() returned nil")
	}

	if sort.Field != "name" {
		t.Errorf("NewSortASC().Field = %v, want name", sort.Field)
	}

	if sort.Direction != SortDirectionASC {
		t.Errorf("NewSortASC().Direction = %v, want ASC", sort.Direction)
	}
}

func TestNewSortDESC(t *testing.T) {
	sort := NewSortDESC("age")

	if sort == nil {
		t.Fatal("NewSortDESC() returned nil")
	}

	if sort.Field != "age" {
		t.Errorf("NewSortDESC().Field = %v, want age", sort.Field)
	}

	if sort.Direction != SortDirectionDESC {
		t.Errorf("NewSortDESC().Direction = %v, want DESC", sort.Direction)
	}
}

func TestSort_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		sort     *Sort
		expected bool
	}{
		{"valid ASC", NewSortASC("name"), true},
		{"valid DESC", NewSortDESC("age"), true},
		{"nil sort", nil, false},
		{"empty field", &Sort{Field: "", Direction: SortDirectionASC}, false},
		{"invalid direction", &Sort{Field: "name", Direction: SortDirection("INVALID")}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sort.IsValid(); got != tt.expected {
				t.Errorf("Sort.IsValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSortDirection_Constants(t *testing.T) {
	if SortDirectionASC != "ASC" {
		t.Errorf("SortDirectionASC = %v, want ASC", SortDirectionASC)
	}

	if SortDirectionDESC != "DESC" {
		t.Errorf("SortDirectionDESC = %v, want DESC", SortDirectionDESC)
	}
}
