package test

import (
	"simple-demo/model"
	"testing"
	"time"
)

func TestUserCreate(t *testing.T) {
	_, err := model.GetDB()
	if err != nil {
		t.Errorf("Init DB failed: %v", err)
	}
}

func TestGetUserByName(t *testing.T) {
	pre_time := time.Now()
	user, err := model.GetUserByName("liuhui")
	if err != nil {
		t.Errorf("Get user failed: %v", err)
	}
	t.Logf("Get user: %v within %v", user.Name, time.Since(pre_time))
}
