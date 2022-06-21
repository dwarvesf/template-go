package util

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPagination_Standardize(t *testing.T) {
	type fields struct {
		Page  int
		Size  int
		Total int
	}
	tests := []struct {
		name   string
		fields fields
		expect Pagination
	}{
		{
			name:   "case all valid",
			fields: fields{Page: 1, Size: 10},
			expect: Pagination{Page: 1, Size: 10},
		},
		{
			name:   "success",
			fields: fields{Page: -1, Size: 60},
			expect: Pagination{Page: 0, Size: 24},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Pagination{
				Page:  tt.fields.Page,
				Size:  tt.fields.Size,
				Total: tt.fields.Total,
			}
			p.Standardize()

			if !cmp.Equal(p, tt.expect) {
				t.Errorf("GetConfig() = %v, want %v \n diff: %v", p, tt.expect, cmp.Diff(p, tt.expect))
				t.FailNow()
			}
		})
	}
}
