package pagination

import (
	"testing"
)

func TestNewPageRequest(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		size     int
		wantPage int
		wantSize int
	}{
		{"valid request", 0, 10, 0, 10},
		{"valid request page 1", 1, 20, 1, 20},
		{"negative page normalized", -1, 10, 0, 10},
		{"zero size normalized", 0, 0, 0, 10},
		{"negative size normalized", 0, -5, 0, 10},
		{"both invalid", -1, -5, 0, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewPageRequest(tt.page, tt.size)
			if got.Page != tt.wantPage {
				t.Errorf("NewPageRequest().Page = %v, want %v", got.Page, tt.wantPage)
			}
			if got.Size != tt.wantSize {
				t.Errorf("NewPageRequest().Size = %v, want %v", got.Size, tt.wantSize)
			}
		})
	}
}

func TestPageRequest_Offset(t *testing.T) {
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
		{"page 0 size 5", 0, 5, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := NewPageRequest(tt.page, tt.size)
			if got := pr.Offset(); got != tt.expected {
				t.Errorf("PageRequest.Offset() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPageRequest_Limit(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		size     int
		expected int
	}{
		{"size 10", 0, 10, 10},
		{"size 20", 1, 20, 20},
		{"size 5", 2, 5, 5},
		{"size 100", 0, 100, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := NewPageRequest(tt.page, tt.size)
			if got := pr.Limit(); got != tt.expected {
				t.Errorf("PageRequest.Limit() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPageRequest_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		size     int
		expected bool
	}{
		{"valid request", 0, 10, true},
		{"valid page 1", 1, 20, true},
		{"valid large page", 100, 50, true},
		{"invalid negative page", -1, 10, false},
		{"invalid zero size", 0, 0, false},
		{"invalid negative size", 0, -5, false},
		{"both invalid", -1, -5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := PageRequest{Page: tt.page, Size: tt.size}
			if got := pr.IsValid(); got != tt.expected {
				t.Errorf("PageRequest.IsValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}
