package pagination

import (
	"testing"
)

func TestNewPageResult(t *testing.T) {
	content := []string{"item1", "item2", "item3"}
	totalElements := int64(25)
	pageRequest := NewPageRequest(0, 10)

	result := NewPageResult(content, totalElements, pageRequest)

	if len(result.Content) != 3 {
		t.Errorf("NewPageResult().Content length = %v, want 3", len(result.Content))
	}

	if result.TotalElements != totalElements {
		t.Errorf("NewPageResult().TotalElements = %v, want %v", result.TotalElements, totalElements)
	}

	if result.TotalPages != 3 {
		t.Errorf("NewPageResult().TotalPages = %v, want 3", result.TotalPages)
	}

	if result.Page != 0 {
		t.Errorf("NewPageResult().Page = %v, want 0", result.Page)
	}

	if result.Size != 10 {
		t.Errorf("NewPageResult().Size = %v, want 10", result.Size)
	}

	if !result.IsFirst {
		t.Error("NewPageResult().IsFirst should be true for page 0")
	}

	if result.IsLast {
		t.Error("NewPageResult().IsLast should be false for page 0 with 3 pages")
	}

	if !result.HasNext {
		t.Error("NewPageResult().HasNext should be true for page 0")
	}

	if result.HasPrevious {
		t.Error("NewPageResult().HasPrevious should be false for page 0")
	}
}

func TestNewPageResult_LastPage(t *testing.T) {
	content := []string{"item1", "item2"}
	totalElements := int64(12)
	pageRequest := NewPageRequest(1, 10) // Second page (0-indexed)

	result := NewPageResult(content, totalElements, pageRequest)

	if !result.IsLast {
		t.Error("NewPageResult().IsLast should be true for last page")
	}

	if result.HasNext {
		t.Error("NewPageResult().HasNext should be false for last page")
	}

	if !result.HasPrevious {
		t.Error("NewPageResult().HasPrevious should be true for page 1")
	}
}

func TestNewPageResult_MiddlePage(t *testing.T) {
	content := []string{"item1", "item2"}
	totalElements := int64(30)
	pageRequest := NewPageRequest(1, 10) // Second page

	result := NewPageResult(content, totalElements, pageRequest)

	if result.IsFirst {
		t.Error("NewPageResult().IsFirst should be false for middle page")
	}

	if result.IsLast {
		t.Error("NewPageResult().IsLast should be false for middle page")
	}

	if !result.HasNext {
		t.Error("NewPageResult().HasNext should be true for middle page")
	}

	if !result.HasPrevious {
		t.Error("NewPageResult().HasPrevious should be true for middle page")
	}
}

func TestNewPageResult_ExactDivision(t *testing.T) {
	content := []string{"item1", "item2"}
	totalElements := int64(20)
	pageRequest := NewPageRequest(1, 10) // Second page

	result := NewPageResult(content, totalElements, pageRequest)

	if result.TotalPages != 2 {
		t.Errorf("NewPageResult().TotalPages = %v, want 2", result.TotalPages)
	}

	if !result.IsLast {
		t.Error("NewPageResult().IsLast should be true for last page with exact division")
	}
}

func TestEmptyPageResult(t *testing.T) {
	pageRequest := NewPageRequest(0, 10)

	result := EmptyPageResult[string](pageRequest)

	if len(result.Content) != 0 {
		t.Errorf("EmptyPageResult().Content length = %v, want 0", len(result.Content))
	}

	if result.TotalElements != 0 {
		t.Errorf("EmptyPageResult().TotalElements = %v, want 0", result.TotalElements)
	}

	if result.TotalPages != 0 {
		t.Errorf("EmptyPageResult().TotalPages = %v, want 0", result.TotalPages)
	}

	if !result.IsFirst {
		t.Error("EmptyPageResult().IsFirst should be true")
	}

	if !result.IsLast {
		t.Error("EmptyPageResult().IsLast should be true")
	}

	if result.HasNext {
		t.Error("EmptyPageResult().HasNext should be false")
	}

	if result.HasPrevious {
		t.Error("EmptyPageResult().HasPrevious should be false")
	}
}
