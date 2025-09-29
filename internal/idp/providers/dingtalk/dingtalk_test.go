package dingtalk

import (
	"testing"
)

func TestProvider_Name(t *testing.T) {
	provider, err := New("clientID", "clientSecret", "https://example.com/callback", []string{"openid"})
	if err != nil {
		t.Errorf("Error creating provider: %v", err)
	}
	if provider.Name() != name {
		t.Errorf("Expected name %s, got %s", name, provider.Name())
	}
}

func TestUser_GetID(t *testing.T) {
	user := &User{
		UnionID: "test_union_id",
		Nick:    "Test User",
		Email:   "test@example.com",
	}
	if user.GetID() != "test_union_id" {
		t.Errorf("Expected ID test_union_id, got %s", user.GetID())
	}
}

func TestUser_GetDisplayName(t *testing.T) {
	user := &User{
		UnionID: "test_union_id",
		Nick:    "Test User",
		Email:   "test@example.com",
	}
	if user.GetDisplayName() != "Test User" {
		t.Errorf("Expected display name Test User, got %s", user.GetDisplayName())
	}
}

func TestUser_GetEmail(t *testing.T) {
	user := &User{
		UnionID: "test_union_id",
		Nick:    "Test User",
		Email:   "test@example.com",
	}
	if string(user.GetEmail()) != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", user.GetEmail())
	}
}

func TestUser_IsEmailVerified(t *testing.T) {
	user := &User{
		UnionID: "test_union_id",
		Nick:    "Test User",
		Email:   "test@example.com",
	}
	if !user.IsEmailVerified() {
		t.Error("Expected email to be verified")
	}
}

func TestUser_GetPhone(t *testing.T) {
	user := &User{
		UnionID: "test_union_id",
		Nick:    "Test User",
		Email:   "test@example.com",
		Mobile:  "13800138000",
	}
	if string(user.GetPhone()) != "13800138000" {
		t.Errorf("Expected phone 13800138000, got %s", user.GetPhone())
	}
}

func TestUser_IsPhoneVerified(t *testing.T) {
	userWithPhone := &User{
		UnionID: "test_union_id",
		Nick:    "Test User",
		Email:   "test@example.com",
		Mobile:  "13800138000",
	}
	if !userWithPhone.IsPhoneVerified() {
		t.Error("Expected phone to be verified when mobile is present")
	}

	userWithoutPhone := &User{
		UnionID: "test_union_id",
		Nick:    "Test User",
		Email:   "test@example.com",
		Mobile:  "",
	}
	if userWithoutPhone.IsPhoneVerified() {
		t.Error("Expected phone to not be verified when mobile is empty")
	}
}
