# Phase 4: Backend Web Layer & Main Application

1. **Implement Authentication Middleware:**
    * Create the file `internal/middleware/auth.go`.
    * Define the `UserIDContextKey`.
    * Implement the `RequireAuth` middleware function:
        * It takes the session store getter (`auth.Store` can be used).
        * It wraps an `http.Handler`.
        * Inside the wrapper, call `auth.GetSessionUserID` to check the session.
        * If valid, add the user ID to the request context using `UserIDContextKey` and call the `next` handler.
        * If invalid, return an `http.StatusUnauthorized` error.
<!-- STOP: Please review and approve before proceeding to the next step. -->

2. **Implement Frontend Embedding:**
    * Create the file `embed.go` in the root directory.
    * Add the `go:embed` directives for `frontend/dist` and `frontend/dist/index.html`. Note that the frontend build output directory is `dist` in this project, not `build`.
<!-- STOP: Please review and approve before proceeding to the next step. -->

3. **Implement Frontend Serving Handler:**
    * Create the file `internal/web/handler.go`.
    * Implement the `FrontendHandler` struct holding the embedded filesystems.
    * Implement the `ServeHTTP` method for `FrontendHandler`:
        * Check if the requested file exists in the `StaticFS` (inside `frontend/dist`).
        * If it exists, serve it using `http.FileServer`.
        * If it *doesn't* exist (likely an SPA route), read and serve the contents of `index.html` from `IndexFS` with `Content-Type: text/html`.
<!-- STOP: Please review and approve before proceeding to the next step. -->

4. **Integrate Auth Service and Middleware into Main Server Logic:**
    * Modify the `main.go` file in the root directory.
    * Ensure `.env` is loaded.
    * Initialize `auth.InitSessionStore()`.
    * Initialize `user.NewInMemoryStore()`.
    * Initialize `auth.NewAuthServiceServer()`.
    * Get the ConnectRPC handler path and handler for the new auth service using `authv1connect.NewAuthServiceHandler()`.
    * Instantiate the `middleware.RequireAuth` middleware using the initialized session store.
    * **Update Routing in `main.go`:**
        * Modify the existing `http.NewServeMux()` to include the new auth service handler.
        * Apply the `middleware.RequireAuth` conditionally to the auth service RPC endpoints that require authentication (`GetCurrentUser`, `Logout`).
        * Integrate the session clearing logic *after* the `Logout` handler runs.
        * Ensure the existing RPC services and the frontend handler are still correctly routed.
    * Ensure CORS middleware is applied correctly to the main mux.
    * Ensure the server is configured and started correctly, including graceful shutdown.
<!-- STOP: Please review and approve before proceeding to the next step. -->
