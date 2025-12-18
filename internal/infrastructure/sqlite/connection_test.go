package sqlite

import (
	"testing"
)

func TestNewConnection(t *testing.T) {
	db, err := NewConnection(":memory:")
	if err != nil {
		t.Fatalf("Failed to open connection: %v", err)
	}

	if db == nil {
		t.Fatal("Expected non-nil *gorm.DB")
	}

	type TestModel struct {
		ID   uint
		Name string
	}
	err = db.AutoMigrate(&TestModel{})
	if err != nil {
		t.Fatalf("Failed to automigrate: %v", err)
	}
}
