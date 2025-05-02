# Phase 3: Backend Core Logic Implementation

1. **Implement User Store:**
    * Create the file `internal/user/store.go`.
    * Implement the `User` struct (internal representation).
    * Implement the `Store` interface.
    * Implement the `InMemoryStore` type satisfying the `Store` interface (`FindByGoogleID`, `FindByID`, `CreateOrUpdateFromGoogle`). Use `uuid.NewString()` for new user IDs. This will be a temporary in-memory store; a persistent database solution will be addressed in a later phase. Adjust import paths if necessary to align with the current project structure.
<!-- STOP: Please review and approve before proceeding to the next step. -->

2. **Implement Session Management:**
    * Create the file `internal/auth/session.go`.
    * Define constants `SessionName` and `UserIDKey`.
    * Declare the global `Store` variable (`sessions.Store`).
    * Implement `InitSessionStore` function to initialize `sessions.NewCookieStore` using `SESSION_SECRET_KEY` and `SESSION_SECURE_COOKIE` from the `.env` file. Configure `HttpOnly`, `Secure`, `SameSite`, `MaxAge`, `Path` options.
    * Implement `GetSessionUserID` function to retrieve the user ID from the session cookie in an `http.Request`.
<!-- STOP: Please review and approve before proceeding to the next step. -->

3. **Implement Auth Service Logic:**
    * Create the file `internal/rpc/services/auth.go`.
    * Implement the `AuthServiceServer` struct holding the `GoogleClientID` and `UserStore`.
    * Implement `NewAuthServiceServer` constructor.
    * Implement the `VerifyGoogleToken` RPC handler:
        * Get token from request.
        * Validate token using `idtoken.Validate` and the backend's `GoogleClientID`.
        * Call `UserStore.CreateOrUpdateFromGoogle` to get/create the application user.
        * Send the `authv1.User` back in the response. *(Note: Session creation will be handled in `main.go` or potentially with interceptors)*.
    * Implement the `GetCurrentUser` RPC handler:
        * Retrieve `userID` from the request context (added by middleware).
        * Call `UserStore.FindByID` to fetch the user.
        * Return the `authv1.User` or `CodeUnauthenticated` / `CodeInternal` error.
    * Implement the `Logout` RPC handler:
        * This handler primarily signals intent. Session clearing will happen in `main.go` or with interceptors.
        * Send an empty response on success.
<!-- STOP: Please review and approve before proceeding to the next step. -->
