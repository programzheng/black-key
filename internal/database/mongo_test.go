package database

import (
	"testing"
)

func TestMongoCreateOne(t *testing.T) {
	mbr := NewMongoBaseRepository("test")
	id, err := mbr.CreateOne(struct {
		Name string
	}{
		Name: "test",
	})
	if err != nil {
		t.Errorf("TestMongoCreateOne error: %v", err)
		return
	}
	t.Logf("TestMongoCreateOne succeeded result: %v", *id)
}

func TestMongoFind(t *testing.T) {
	mbr := NewMongoBaseRepository("test")
	f := struct {
		Name string
	}{
		Name: "test",
	}
	ms, err := mbr.Find(f)
	if err != nil {
		t.Errorf("TestMongoFind error: %v", err)
		return
	}
	t.Logf("TestMongoFind succeeded result: %v", ms)
}
