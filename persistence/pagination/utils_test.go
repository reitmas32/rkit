package pagination

import (
	"testing"
)

func TestValidatePageRequest(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		size     int
		wantPage int
		wantSize int
	}{
		{"valid request", 0, 10, 0, 10},
		{"negative page normalized", -1, 10, 0, 10},
		{"zero size normalized", 0, 0, 0, 10},
		{"negative size normalized", 0, -5, 0, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidatePageRequest(tt.page, tt.size)
			if got.Page != tt.wantPage {
				t.Errorf("ValidatePageRequest().Page = %v, want %v", got.Page, tt.wantPage)
			}
			if got.Size != tt.wantSize {
				t.Errorf("ValidatePageRequest().Size = %v, want %v", got.Size, tt.wantSize)
			}
		})
	}
}

func TestCalculateOffset(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		size     int
		expected int
	}{
		{"first page", 0, 10, 0},
		{"second page", 1, 10, 10},
		{"third page", 2, 10, 20},
		{"page 5 size 20", 5, 20, 100},
		{"negative page normalized", -1, 10, 0},
		{"zero size normalized", 0, 0, 0},
		{"negative size normalized", 0, -5, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateOffset(tt.page, tt.size)
			if got != tt.expected {
				t.Errorf("CalculateOffset(%v, %v) = %v, want %v", tt.page, tt.size, got, tt.expected)
			}
		})
	}
}

func TestCalculateTotalPages(t *testing.T) {
	tests := []struct {
		name          string
		totalElements int64
		size          int64
		expected      int
	}{
		{"exact division", 20, 10, 2},
		{"with remainder", 25, 10, 3},
		{"single page", 5, 10, 1},
		{"zero elements", 0, 10, 0},
		{"zero size", 10, 0, 0},
		{"both zero", 0, 0, 0},
		{"large numbers", 1000, 25, 40},
		{"one element", 1, 10, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateTotalPages(tt.totalElements, tt.size)
			if got != tt.expected {
				t.Errorf("CalculateTotalPages(%v, %v) = %v, want %v", tt.totalElements, tt.size, got, tt.expected)
			}
		})
	}
}

func TestIsValidPage(t *testing.T) {
	tests := []struct {
		name      string
		page      int
		totalPages int
		expected  bool
	}{
		{"first page valid", 0, 5, true},
		{"middle page valid", 2, 5, true},
		{"last page valid", 4, 5, true},
		{"page equals total pages", 5, 5, false},
		{"page exceeds total pages", 10, 5, false},
		{"negative page", -1, 5, false},
		{"zero total pages", 0, 0, false},
		{"page 0 with 1 total page", 0, 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidPage(tt.page, tt.totalPages)
			if got != tt.expected {
				t.Errorf("IsValidPage(%v, %v) = %v, want %v", tt.page, tt.totalPages, got, tt.expected)
			}
		})
	}
}
