package contact

import (
	"path/filepath"
	"testing"

	"github.com/hkdb/aerion/internal/database"
)

func openTestDB(t *testing.T) *database.DB {
	t.Helper()
	db, err := database.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Migrate(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestUpsert(t *testing.T) {
	db := openTestDB(t)
	store := NewStore(db.DB)

	if err := store.AddOrUpdate("alice@example.com", "Alice"); err != nil {
		t.Fatalf("AddOrUpdate failed: %v", err)
	}

	got, err := store.Get("alice@example.com")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got == nil {
		t.Fatal("expected contact, got nil")
	}
	if got.DisplayName != "Alice" {
		t.Errorf("DisplayName = %q, want %q", got.DisplayName, "Alice")
	}
	if got.SendCount != 1 {
		t.Errorf("SendCount = %d, want 1", got.SendCount)
	}

	// Upsert again to increment send count
	if err := store.AddOrUpdate("alice@example.com", "Alice Smith"); err != nil {
		t.Fatalf("AddOrUpdate (second) failed: %v", err)
	}

	got, err = store.Get("alice@example.com")
	if err != nil {
		t.Fatalf("Get after second upsert failed: %v", err)
	}
	if got.SendCount != 2 {
		t.Errorf("SendCount = %d, want 2", got.SendCount)
	}
	if got.DisplayName != "Alice Smith" {
		t.Errorf("DisplayName = %q, want %q", got.DisplayName, "Alice Smith")
	}
}

func TestSearch(t *testing.T) {
	db := openTestDB(t)
	store := NewStore(db.DB)

	contacts := []struct {
		email, name string
	}{
		{"alice@example.com", "Alice"},
		{"bob@example.com", "Bob"},
		{"alicia@test.com", "Alicia"},
	}
	for _, c := range contacts {
		if err := store.AddOrUpdate(c.email, c.name); err != nil {
			t.Fatalf("AddOrUpdate failed: %v", err)
		}
	}

	results, err := store.Search("ali", 10)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("Search returned %d results, want 2", len(results))
	}

	// Verify both alice and alicia are returned
	emails := make(map[string]bool)
	for _, r := range results {
		emails[r.Email] = true
	}
	if !emails["alice@example.com"] {
		t.Error("expected alice@example.com in results")
	}
	if !emails["alicia@test.com"] {
		t.Error("expected alicia@test.com in results")
	}
}

func TestSearchEmpty(t *testing.T) {
	db := openTestDB(t)
	store := NewStore(db.DB)

	results, err := store.Search("nonexistent", 10)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("Search returned %d results, want 0", len(results))
	}
}
