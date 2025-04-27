package services

import (
	"testing"
	"time"

	"github.com/atreya2011/expense-manager/internal/clock"
)

// TestMockClock verifies that we can use the MockClock in tests
func TestMockClock(t *testing.T) {
	// Get the test clock as a MockClock
	mockClock, ok := testClock.(*clock.MockClock)
	if !ok {
		t.Fatalf("testClock is not a *clock.MockClock")
	}

	// Verify initial time
	initialTime := testClock.Now()
	expectedInitialTime := time.Date(2025, 4, 26, 12, 0, 0, 0, time.UTC)
	if !initialTime.Equal(expectedInitialTime) {
		t.Errorf("Expected initial time %v, got %v", expectedInitialTime, initialTime)
	}

	// Set a new time
	newTime := time.Date(2025, 5, 1, 15, 30, 0, 0, time.UTC)
	mockClock.SetTime(newTime)

	// Verify the time was updated
	updatedTime := testClock.Now()
	if !updatedTime.Equal(newTime) {
		t.Errorf("Expected updated time %v, got %v", newTime, updatedTime)
	}

	// Reset to original time for other tests
	mockClock.SetTime(expectedInitialTime)
}
