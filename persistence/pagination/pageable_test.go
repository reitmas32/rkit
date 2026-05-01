package pagination

import (
	"testing"
)

func TestNewPageable(t *testing.T) {
	pageRequest := NewPageRequest(0, 10)
	sort := NewSortASC("name")

	pageable := NewPageable(pageRequest, sort)

	if pageable == nil {
		t.Fatal("NewPageable() returned nil")
	}

	if pageable.PageRequest.Page != 0 {
		t.Errorf("NewPageable().PageRequest.Page = %v, want 0", pageable.PageRequest.Page)
	}

	if pageable.PageRequest.Size != 10 {
		t.Errorf("NewPageable().PageRequest.Size = %v, want 10", pageable.PageRequest.Size)
	}

	if pageable.Sort == nil {
		t.Error("NewPageable().Sort should not be nil")
	}

	if pageable.Sort.Field != "name" {
		t.Errorf("NewPageable().Sort.Field = %v, want name", pageable.Sort.Field)
	}
}

func TestNewPageable_NilSort(t *testing.T) {
	pageRequest := NewPageRequest(0, 10)

	pageable := NewPageable(pageRequest, nil)

	if pageable == nil {
		t.Fatal("NewPageable() returned nil")
	}

	if pageable.Sort != nil {
		t.Error("NewPageable() with nil sort should have nil Sort")
	}
}

func TestNewPageableWithoutSort(t *testing.T) {
	pageable := NewPageableWithoutSort(0, 10)

	if pageable == nil {
		t.Fatal("NewPageableWithoutSort() returned nil")
	}

	if pageable.PageRequest.Page != 0 {
		t.Errorf("NewPageableWithoutSort().PageRequest.Page = %v, want 0", pageable.PageRequest.Page)
	}

	if pageable.PageRequest.Size != 10 {
		t.Errorf("NewPageableWithoutSort().PageRequest.Size = %v, want 10", pageable.PageRequest.Size)
	}

	if pageable.Sort != nil {
		t.Error("NewPageableWithoutSort().Sort should be nil")
	}
}

func TestPageable_Offset(t *testing.T) {
	pageable := NewPageableWithoutSort(2, 20)

	if got := pageable.Offset(); got != 40 {
		t.Errorf("Pageable.Offset() = %v, want 40", got)
	}
}

func TestPageable_Limit(t *testing.T) {
	pageable := NewPageableWithoutSort(0, 15)

	if got := pageable.Limit(); got != 15 {
		t.Errorf("Pageable.Limit() = %v, want 15", got)
	}
}

func TestPageable_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		size     int
		expected bool
	}{
		{"valid", 0, 10, true},
		{"valid page 1", 1, 20, true},
		{"normalized negative page becomes valid", -1, 10, true}, // NewPageableWithoutSort normalizes
		{"normalized zero size becomes valid", 0, 0, true},       // NewPageableWithoutSort normalizes
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pageable := NewPageableWithoutSort(tt.page, tt.size)
			if got := pageable.IsValid(); got != tt.expected {
				t.Errorf("Pageable.IsValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPageable_IsValid_WithInvalidValues(t *testing.T) {
	// Test with directly created Pageable (not normalized)
	pageable := &Pageable{
		PageRequest: PageRequest{Page: -1, Size: 10},
	}
	if pageable.IsValid() {
		t.Error("Pageable with negative page should be invalid")
	}

	pageable2 := &Pageable{
		PageRequest: PageRequest{Page: 0, Size: 0},
	}
	if pageable2.IsValid() {
		t.Error("Pageable with zero size should be invalid")
	}
}
