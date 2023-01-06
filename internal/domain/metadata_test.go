package domain

import "testing"

func TestMetadatas_Append(t *testing.T) {
	type args struct {
		md *Metadata
	}
	tests := []struct {
		name string
		m    Metadatas
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Append(tt.args.md)
		})
	}
}
