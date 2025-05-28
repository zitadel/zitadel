package query

import "testing"

func TestGetSystemTranslation(t *testing.T) {
	_, err := getSystemTranslation("en")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	_, err = getSystemTranslation("invalid-lang")
	if err == nil {
		t.Error("expected error for invalid language, got nil")
	}
}
