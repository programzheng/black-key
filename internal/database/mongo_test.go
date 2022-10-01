package database

import (
	"testing"
)

func TestMongoCreateOne(t *testing.T) {
	md := NewMongoInstance()
	r, err := md.CreateOne("test", struct {
		Name string
	}{
		Name: "test",
	})
	if err != nil {
		t.Errorf("TestMongoCreateOne error: %v", err)
		return
	}
	t.Logf("TestMongoCreateOne succeeded result: %v", r)
}
