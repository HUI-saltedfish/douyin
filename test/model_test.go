package test

import (
	"simple-demo/model"
	"testing"
)

func TestUserCreate(t *testing.T) {
	_, err := model.GetDB()
	if err != nil {
		t.Errorf("Init DB failed: %v", err)
	}
}

func TestGetUserByName(t *testing.T) {
	_, err := model.GetDB()
	if err != nil {
		t.Errorf("Init DB failed: %v", err)
	}
	user, err := model.GetUserByName("liuhui")
	if err != nil {
		t.Errorf("Get user failed: %v", err)
	}
	t.Logf("Get user: %v", user)
}