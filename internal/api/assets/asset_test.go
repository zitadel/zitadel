package assets

import "testing"

func TestRefineContentType(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		filename    string
		want        string
	}{
		{
			name:        "octet-stream ttf is refined to font/ttf",
			contentType: "application/octet-stream",
			filename:    "custom-font.TTF",
			want:        "font/ttf",
		},
		{
			name:        "octet-stream otf is refined to font/otf",
			contentType: "application/octet-stream",
			filename:    "custom-font.otf",
			want:        "font/otf",
		},
		{
			name:        "octet-stream woff2 is refined to font/woff2",
			contentType: "application/octet-stream",
			filename:    "custom-font.woff2",
			want:        "font/woff2",
		},
		{
			name:        "octet-stream with unknown extension is left untouched",
			contentType: "application/octet-stream",
			filename:    "custom-font.bin",
			want:        "application/octet-stream",
		},
		{
			name:        "correctly detected content type is left untouched",
			contentType: "font/ttf",
			filename:    "custom-font.ttf",
			want:        "font/ttf",
		},
		{
			name:        "non-font content type is left untouched",
			contentType: "image/png",
			filename:    "logo.ttf",
			want:        "image/png",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := refineContentType(tt.contentType, tt.filename); got != tt.want {
				t.Errorf("refineContentType() = %v, want %v", got, tt.want)
			}
		})
	}
}
