package services

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	expensesv1 "github.com/atreya2011/expense-manager/internal/rpc/gen/expenses/v1"
)

// TestCreateUser tests the CreateUser RPC method
func TestCreateUser(t *testing.T) {
	// Reset the test database
	resetTestDB(t)

	// Create a new UserService with the test repositories
	service := NewUserService(userRepo, testClock)

	// Define test cases
	tests := []struct {
		name        string
		request     *expensesv1.CreateUserRequest
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid user creation",
			request: &expensesv1.CreateUserRequest{
				Name:  "Test User",
				Email: "test@example.com",
			},
			expectError: false,
		},
		{
			name: "Missing name",
			request: &expensesv1.CreateUserRequest{
				Name:  "",
				Email: "no-name@example.com",
			},
			expectError: true,
			errorMsg:    "name is required",
		},
		{
			name: "Missing email",
			request: &expensesv1.CreateUserRequest{
				Name:  "No Email User",
				Email: "",
			},
			expectError: true,
			errorMsg:    "email is required",
		},
		{
			name: "Duplicate email",
			request: &expensesv1.CreateUserRequest{
				Name:  "Duplicate Email",
				Email: "duplicate@example.com",
			},
			expectError: false, // First insertion should succeed
		},
		{
			name: "Duplicate email (should fail)",
			request: &expensesv1.CreateUserRequest{
				Name:  "Another User",
				Email: "duplicate@example.com", // Same email as above
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
			resp, err := service.CreateUser(ctx, req)

			// Check errors
			assertError(t, err, tc.expectError, tc.errorMsg)

			// Verify response for successful cases
			if !tc.expectError {
				if resp == nil || resp.Msg == nil || resp.Msg.User == nil {
					t.Fatalf("Expected valid response, got nil")
				}

				if resp.Msg.User.Id == "" {
					t.Errorf("Expected user ID to be generated, got empty string")
				}

				if resp.Msg.User.Name != tc.request.Name {
					t.Errorf("Expected name=%s, got %s", tc.request.Name, resp.Msg.User.Name)
				}

				if resp.Msg.User.Email != tc.request.Email {
					t.Errorf("Expected email=%s, got %s", tc.request.Email, resp.Msg.User.Email)
				}
			}
		})
	}
}

// TestGetUser tests the GetUser RPC method
func TestGetUser(t *testing.T) {
	// Reset the test database
	resetTestDB(t)

	// Create a new UserService with the test repositories
	service := NewUserService(userRepo, testClock)

	// Create a test user (using the main DB connection for setup)
	testUser := createTestUser(t, testDB, "Get Test User", "get@example.com")

	// Define test cases
	tests := []struct {
		name        string
		userID      string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid user retrieval",
			userID:      testUser.ID,
			expectError: false,
		},
		{
			name:        "Non-existent user",
			userID:      "usr_nonexistent",
			expectError: true,
			errorMsg:    "not found",
		},
		{
			name:        "Empty user ID",
			userID:      "",
			expectError: true,
			errorMsg:    "id is required",
		},
	}

	// Run tests
	ctx := context.Background()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := connect.NewRequest(&expensesv1.GetUserRequest{
				Id: tc.userID,
			})

			// Read operations can use the main DB connection
			resp, err := service.GetUser(ctx, req)

			// Check errors
			assertError(t, err, tc.expectError, tc.errorMsg)

			// Verify response for successful cases
			if !tc.expectError {
				if resp == nil || resp.Msg == nil || resp.Msg.User == nil {
					t.Fatalf("Expected valid response, got nil")
				}

				if resp.Msg.User.Id != testUser.ID {
					t.Errorf("Expected ID=%s, got %s", testUser.ID, resp.Msg.User.Id)
				}

				if resp.Msg.User.Name != testUser.Name {
					t.Errorf("Expected name=%s, got %s", testUser.Name, resp.Msg.User.Name)
				}

				if resp.Msg.User.Email != testUser.Email {
					t.Errorf("Expected email=%s, got %s", testUser.Email, resp.Msg.User.Email)
				}
			}
		})
	}
}

// TestListUsers tests the ListUsers RPC method
func TestListUsers(t *testing.T) {
	// Reset the test database
	resetTestDB(t)

	// Create a new UserService with the test repositories
	service := NewUserService(userRepo, testClock)

	// Create test users (using the main DB connection for setup)
	_ = createTestUser(t, testDB, "List User 1", "list1@example.com")
	_ = createTestUser(t, testDB, "List User 2", "list2@example.com")
	_ = createTestUser(t, testDB, "List User 3", "list3@example.com")

	// Define test cases
	tests := []struct {
		name        string
		pageSize    int32
		pageToken   string
		expectedN   int
		expectError bool
	}{
		{
			name:      "List all users (default page size)",
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
			expectedN: 1,   // Only 1 user left
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

			req := connect.NewRequest(&expensesv1.ListUsersRequest{
				Pagination: pagination,
			})

			// Read operations can use the main DB connection
			resp, err := service.ListUsers(ctx, req)

			// Check errors
			assertError(t, err, tc.expectError, "")

			// Verify response for successful cases
			if !tc.expectError {
				if resp == nil || resp.Msg == nil {
					t.Fatalf("Expected valid response, got nil")
				}

				if len(resp.Msg.Users) != tc.expectedN {
					t.Errorf("Expected %d users, got %d", tc.expectedN, len(resp.Msg.Users))
				}

				// Verify pagination for first page (when we expect more data)
				if tc.pageSize > 0 && tc.expectedN == int(tc.pageSize) && len(resp.Msg.Users) == int(tc.pageSize) {
					if resp.Msg.PaginationResponse == nil || resp.Msg.PaginationResponse.NextPageToken == "" {
						t.Errorf("Expected next page token, got empty")
					}
				}
			}
		})
	}
}

// TestUpdateUser tests the UpdateUser RPC method
func TestUpdateUser(t *testing.T) {
	// Reset the test database
	resetTestDB(t)

	// Create a new UserService with the test repositories
	service := NewUserService(userRepo, testClock)

	// Create test users (using the main DB connection for setup)
	testUser := createTestUser(t, testDB, "Original Name", "original@example.com")
	otherUser := createTestUser(t, testDB, "Other User", "other@example.com")

	// Define test cases
	tests := []struct {
		name        string
		request     *expensesv1.UpdateUserRequest
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid user update",
			request: &expensesv1.UpdateUserRequest{
				Id:    testUser.ID,
				Name:  "Updated Name",
				Email: "updated@example.com",
			},
			expectError: false,
		},
		{
			name: "Non-existent user",
			request: &expensesv1.UpdateUserRequest{
				Id:    "usr_nonexistent",
				Name:  "Updated Name",
				Email: "updated@example.com",
			},
			expectError: true,
			errorMsg:    "not found",
		},
		{
			name: "Empty user ID",
			request: &expensesv1.UpdateUserRequest{
				Id:    "",
				Name:  "Updated Name",
				Email: "updated@example.com",
			},
			expectError: true,
			errorMsg:    "id is required",
		},
		{
			name: "Missing name",
			request: &expensesv1.UpdateUserRequest{
				Id:    testUser.ID,
				Name:  "",
				Email: "updated@example.com",
			},
			expectError: true,
			errorMsg:    "name is required",
		},
		{
			name: "Missing email",
			request: &expensesv1.UpdateUserRequest{
				Id:    testUser.ID,
				Name:  "Updated Name",
				Email: "",
			},
			expectError: true,
			errorMsg:    "email is required",
		},
		{
			name: "Duplicate email",
			request: &expensesv1.UpdateUserRequest{
				Id:    testUser.ID,
				Name:  "Updated Name",
				Email: otherUser.Email, // Duplicate email
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
			resp, err := service.UpdateUser(ctx, req)

			// Check errors
			assertError(t, err, tc.expectError, tc.errorMsg)

			// Verify response for successful cases
			if !tc.expectError {
				if resp == nil || resp.Msg == nil || resp.Msg.User == nil {
					t.Fatalf("Expected valid response, got nil")
				}

				if resp.Msg.User.Id != tc.request.Id {
					t.Errorf("Expected ID=%s, got %s", tc.request.Id, resp.Msg.User.Id)
				}

				if resp.Msg.User.Name != tc.request.Name {
					t.Errorf("Expected name=%s, got %s", tc.request.Name, resp.Msg.User.Name)
				}

				if resp.Msg.User.Email != tc.request.Email {
					t.Errorf("Expected email=%s, got %s", tc.request.Email, resp.Msg.User.Email)
				}
			}
		})
	}
}

// TestDeleteUser tests the DeleteUser RPC method
func TestDeleteUser(t *testing.T) {
	// Reset the test database
	resetTestDB(t)

	// Create a new UserService with the test repositories
	service := NewUserService(userRepo, testClock)

	// Create a test user (using the main DB connection for setup)
	testUser := createTestUser(t, testDB, "Delete Test User", "delete@example.com")

	// Define test cases
	tests := []struct {
		name        string
		userID      string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid user deletion",
			userID:      testUser.ID,
			expectError: false,
		},
		{
			name:        "Non-existent user",
			userID:      "usr_nonexistent",
			expectError: true,
			errorMsg:    "not found",
		},
		{
			name:        "Empty user ID",
			userID:      "",
			expectError: true,
			errorMsg:    "id is required",
		},
	}

	// Run tests
	ctx := context.Background()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Reset test data for each test case so we can delete the same user
			// in the "Valid user deletion" case each time
			if tc.name != "Valid user deletion" {
				// Only reset for non-deletion cases
				resetTestDB(t)
				createTestUser(t, testDB, testUser.Name, testUser.Email)
			}

			req := connect.NewRequest(&expensesv1.DeleteUserRequest{
				Id: tc.userID,
			})

			resp, err := service.DeleteUser(ctx, req)

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

				// Verify the user was actually deleted
				var count int
				countErr := testDB.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", tc.userID).Scan(&count)
				if countErr != nil {
					t.Fatalf("Failed to query deleted user: %v", countErr)
				}

				if count != 0 {
					t.Errorf("User was not deleted, found %d records", count)
				}
			}
		})
	}
}
