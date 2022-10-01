package database

import (
	"testing"
)

func TestMongoCreateOne(t *testing.T) {
	mbr := NewMongoBaseRepository()
	mbr.CollectionName = "test"
	r, err := mbr.CreateOne(struct {
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
