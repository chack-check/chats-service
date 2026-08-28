package dtos

import (
	"reflect"
	"testing"
)

func TestNewPaginatedResponse(t *testing.T) {
	data := []string{"first", "second"}
	response := NewPaginatedResponse(2, 10, 4, 32, data)

	if response.GetPage() != 2 {
		t.Errorf("GetPage() = %d, want 2", response.GetPage())
	}
	if response.GetPerPage() != 10 {
		t.Errorf("GetPerPage() = %d, want 10", response.GetPerPage())
	}
	if response.GetPagesCount() != 4 {
		t.Errorf("GetPagesCount() = %d, want 4", response.GetPagesCount())
	}
	if response.GetTotal() != 32 {
		t.Errorf("GetTotal() = %d, want 32", response.GetTotal())
	}
	if !reflect.DeepEqual(response.GetData(), data) {
		t.Errorf("GetData() = %v, want %v", response.GetData(), data)
	}
}
