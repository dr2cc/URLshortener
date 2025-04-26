package storage

import (
	"testing"
)

func TestURLStorage_InsertAndGet(t *testing.T) {
	s := NewStorage()
	testURL := "https://example.com"
	testID := "abc123"

	// Test Insert
	err := s.InsertURL(testID, testURL)
	if err != nil {
		t.Errorf("InsertURL failed: %v", err)
	}

	// Test Get existing
	url, err := s.GetURL(testID)
	if err != nil {
		t.Errorf("GetURL failed: %v", err)
	}
	if url != testURL {
		t.Errorf("Expected %s, got %s", testURL, url)
	}

	// Test Get non-existing
	_, err = s.GetURL("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent URL")
	}
}

func TestMakeNewEntry(t *testing.T) {
	s := NewStorage()
	testURL := "https://example.com"
	testID := "def456"

	MakeEntry(s, testID, testURL)

	url, err := s.GetURL(testID)
	if err != nil {
		t.Errorf("MakeEntry failed: %v", err)
	}
	if url != testURL {
		t.Errorf("Expected %s, got %s", testURL, url)
	}
}

func TestGetEntry(t *testing.T) {
	s := NewStorage()
	testURL := "https://example.com"
	testID := "ghi789"

	s.InsertURL(testID, testURL)

	url, err := GetEntry(s, testID)
	if err != nil {
		t.Errorf("GetEntry failed: %v", err)
	}
	if url != testURL {
		t.Errorf("Expected %s, got %s", testURL, url)
	}
}
