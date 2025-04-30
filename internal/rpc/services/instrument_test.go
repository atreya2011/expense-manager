package services

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	expensesv1 "github.com/atreya2011/expense-manager/internal/rpc/gen/expenses/v1"
)

// TestCreateInstrument tests the CreateInstrument RPC method
func TestCreateInstrument(t *testing.T) {
	// Reset the test database
	resetTestDB(t)

	// Create a new InstrumentService with the test repositories
	service := NewInstrumentService(instrumentRepo, testClock)

	// Define test cases
	tests := []struct {
		name        string
		request     *expensesv1.CreateInstrumentRequest
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid instrument creation",
			request: &expensesv1.CreateInstrumentRequest{
				Name: "Cash",
			},
			expectError: false,
		},
		{
			name: "Missing name",
			request: &expensesv1.CreateInstrumentRequest{
				Name: "",
			},
			expectError: true,
			errorMsg:    "name is required",
		},
		{
			name: "Duplicate name",
			request: &expensesv1.CreateInstrumentRequest{
				Name: "Duplicate",
			},
			expectError: false, // First insertion should succeed
		},
		{
			name: "Duplicate name (should fail)",
			request: &expensesv1.CreateInstrumentRequest{
				Name: "Duplicate", // Same name as above
			},
			expectError: true,
			errorMsg:    "already exists",
		},
	}

	// Run tests
	ctx := context.Background()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := connect.NewRequest(tc.request)
			resp, err := service.CreateInstrument(ctx, req)

			// Check errors
			assertError(t, err, tc.expectError, tc.errorMsg)

			// Verify response for successful cases
			if !tc.expectError {
				if resp == nil || resp.Msg == nil || resp.Msg.Instrument == nil {
					t.Fatalf("Expected valid response, got nil")
				}

				if resp.Msg.Instrument.Id == "" {
					t.Errorf("Expected instrument ID to be generated, got empty string")
				}

				if resp.Msg.Instrument.Name != tc.request.Name {
					t.Errorf("Expected name=%s, got %s", tc.request.Name, resp.Msg.Instrument.Name)
				}
			}
		})
	}
}

// TestGetInstrument tests the GetInstrument RPC method
func TestGetInstrument(t *testing.T) {
	// Reset the test database
	resetTestDB(t)

	// Create a new InstrumentService with the test repositories
	service := NewInstrumentService(instrumentRepo, testClock)

	// Create a test instrument
	testInstrument := createTestInstrument(t, "Bank Account")

	// Define test cases
	tests := []struct {
		name         string
		instrumentID string
		expectError  bool
		errorMsg     string
	}{
		{
			name:         "Valid instrument retrieval",
			instrumentID: testInstrument.ID,
			expectError:  false,
		},
		{
			name:         "Non-existent instrument",
			instrumentID: "ins_nonexistent",
			expectError:  true,
			errorMsg:     "not found",
		},
		{
			name:         "Empty instrument ID",
			instrumentID: "",
			expectError:  true,
			errorMsg:     "id is required",
		},
	}

	// Run tests
	ctx := context.Background()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := connect.NewRequest(&expensesv1.GetInstrumentRequest{
				Id: tc.instrumentID,
			})

			resp, err := service.GetInstrument(ctx, req)

			// Check errors
			assertError(t, err, tc.expectError, tc.errorMsg)

			// Verify response for successful cases
			if !tc.expectError {
				if resp == nil || resp.Msg == nil || resp.Msg.Instrument == nil {
					t.Fatalf("Expected valid response, got nil")
				}

				if resp.Msg.Instrument.Id != testInstrument.ID {
					t.Errorf("Expected ID=%s, got %s", testInstrument.ID, resp.Msg.Instrument.Id)
				}

				if resp.Msg.Instrument.Name != testInstrument.Name {
					t.Errorf("Expected name=%s, got %s", testInstrument.Name, resp.Msg.Instrument.Name)
				}
			}
		})
	}
}

// TestListInstruments tests the ListInstruments RPC method
func TestListInstruments(t *testing.T) {
	// Reset the test database
	resetTestDB(t)

	// Create a new InstrumentService with the test repositories
	service := NewInstrumentService(instrumentRepo, testClock)

	// Create test instruments
	_ = createTestInstrument(t, "Cash")
	_ = createTestInstrument(t, "Bank Account")
	_ = createTestInstrument(t, "Credit Card")

	// Define test cases
	tests := []struct {
		name        string
		pageSize    int32
		pageToken   string
		expectedN   int
		expectError bool
	}{
		{
			name:      "List all instruments (default page size)",
			pageSize:  0, // Default
			pageToken: "",
			expectedN: 3,
		},
		{
			name:      "List with pagination (page 1)",
			pageSize:  2,
			pageToken: "",
			expectedN: 2,
		},
		{
			name:      "List with pagination (page 2)",
			pageSize:  2,
			pageToken: "2", // Start from offset 2
			expectedN: 1,   // Only 1 instrument left
		},
		{
			name:      "Empty result",
			pageSize:  2,
			pageToken: "10", // Start from offset 10 (beyond available data)
			expectedN: 0,
		},
	}

	// Run tests
	ctx := context.Background()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var pagination *expensesv1.Pagination
			if tc.pageSize > 0 || tc.pageToken != "" {
				pagination = &expensesv1.Pagination{
					PageSize:  tc.pageSize,
					PageToken: tc.pageToken,
				}
			}

			req := connect.NewRequest(&expensesv1.ListInstrumentsRequest{
				Pagination: pagination,
			})

			resp, err := service.ListInstruments(ctx, req)

			// Check errors
			assertError(t, err, tc.expectError, "")

			// Verify response for successful cases
			if !tc.expectError {
				if resp == nil || resp.Msg == nil {
					t.Fatalf("Expected valid response, got nil")
				}

				if len(resp.Msg.Instruments) != tc.expectedN {
					t.Errorf("Expected %d instruments, got %d", tc.expectedN, len(resp.Msg.Instruments))
				}

				// Verify pagination for first page (when we expect more data)
				if tc.pageSize > 0 && tc.expectedN == int(tc.pageSize) && len(resp.Msg.Instruments) == int(tc.pageSize) {
					if resp.Msg.PaginationResponse == nil || resp.Msg.PaginationResponse.NextPageToken == "" {
						t.Errorf("Expected next page token, got empty")
					}
				}
			}
		})
	}
}

// TestUpdateInstrument tests the UpdateInstrument RPC method
func TestUpdateInstrument(t *testing.T) {
	// Reset the test database
	resetTestDB(t)

	// Create a new InstrumentService with the test repositories
	service := NewInstrumentService(instrumentRepo, testClock)

	// Create test instruments
	testInstrument := createTestInstrument(t, "Original Name")
	otherInstrument := createTestInstrument(t, "Other Instrument")

	// Define test cases
	tests := []struct {
		name        string
		request     *expensesv1.UpdateInstrumentRequest
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid instrument update",
			request: &expensesv1.UpdateInstrumentRequest{
				Id:   testInstrument.ID,
				Name: "Updated Name",
			},
			expectError: false,
		},
		{
			name: "Non-existent instrument",
			request: &expensesv1.UpdateInstrumentRequest{
				Id:   "ins_nonexistent",
				Name: "Updated Name",
			},
			expectError: true,
			errorMsg:    "not found",
		},
		{
			name: "Empty instrument ID",
			request: &expensesv1.UpdateInstrumentRequest{
				Id:   "",
				Name: "Updated Name",
			},
			expectError: true,
			errorMsg:    "id is required",
		},
		{
			name: "Missing name",
			request: &expensesv1.UpdateInstrumentRequest{
				Id:   testInstrument.ID,
				Name: "",
			},
			expectError: true,
			errorMsg:    "name is required",
		},
		{
			name: "Duplicate name",
			request: &expensesv1.UpdateInstrumentRequest{
				Id:   testInstrument.ID,
				Name: otherInstrument.Name, // Duplicate name
			},
			expectError: true,
			errorMsg:    "already exists",
		},
	}

	// Run tests
	ctx := context.Background()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := connect.NewRequest(tc.request)
			resp, err := service.UpdateInstrument(ctx, req)

			// Check errors
			assertError(t, err, tc.expectError, tc.errorMsg)

			// Verify response for successful cases
			if !tc.expectError {
				if resp == nil || resp.Msg == nil || resp.Msg.Instrument == nil {
					t.Fatalf("Expected valid response, got nil")
				}

				if resp.Msg.Instrument.Id != tc.request.Id {
					t.Errorf("Expected ID=%s, got %s", tc.request.Id, resp.Msg.Instrument.Id)
				}

				if resp.Msg.Instrument.Name != tc.request.Name {
					t.Errorf("Expected name=%s, got %s", tc.request.Name, resp.Msg.Instrument.Name)
				}
			}
		})
	}
}

// TestDeleteInstrument tests the DeleteInstrument RPC method
func TestDeleteInstrument(t *testing.T) {
	// Reset the test database
	resetTestDB(t)

	// Create a new InstrumentService with the test repositories
	service := NewInstrumentService(instrumentRepo, testClock)

	// Create a test instrument
	testInstrument := createTestInstrument(t, "Delete Test Instrument")

	// Define test cases
	tests := []struct {
		name         string
		instrumentID string
		expectError  bool
		errorMsg     string
	}{
		{
			name:         "Valid instrument deletion",
			instrumentID: testInstrument.ID,
			expectError:  false,
		},
		{
			name:         "Non-existent instrument",
			instrumentID: "ins_nonexistent",
			expectError:  true,
			errorMsg:     "not found",
		},
		{
			name:         "Empty instrument ID",
			instrumentID: "",
			expectError:  true,
			errorMsg:     "id is required",
		},
	}

	// Run tests
	ctx := context.Background()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Reset test data for each test case so we can delete the same instrument
			// in the "Valid instrument deletion" case each time
			if tc.name != "Valid instrument deletion" {
				// Only reset for non-deletion cases
				resetTestDB(t)
				createTestInstrument(t, testInstrument.Name)
			}

			req := connect.NewRequest(&expensesv1.DeleteInstrumentRequest{
				Id: tc.instrumentID,
			})

			resp, err := service.DeleteInstrument(ctx, req)

			// Check errors
			assertError(t, err, tc.expectError, tc.errorMsg)

			// Verify response for successful cases
			if !tc.expectError {
				if resp == nil || resp.Msg == nil {
					t.Fatalf("Expected valid response, got nil")
				}

				if !resp.Msg.Success {
					t.Errorf("Expected success=true, got false")
				}

				// Verify the instrument was actually deleted
				var count int
				countErr := testDB.QueryRow("SELECT COUNT(*) FROM instruments WHERE id = ?", tc.instrumentID).Scan(&count)
				if countErr != nil {
					t.Fatalf("Failed to query deleted instrument: %v", countErr)
				}

				if count != 0 {
					t.Errorf("Instrument was not deleted, found %d records", count)
				}
			}
		})
	}
}
