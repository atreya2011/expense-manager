# Phase 6: Frontend Core Logic

1. **Implement API Service:**
    * Create `frontend/src/services/api.ts`.
    * Define the `AppUser` interface matching the backend's `authv1.User` proto message.
    * Implement the `apiFetch` helper function using `fetch` to call backend endpoints, handling JSON request/response and basic error parsing (checking `response.ok` and `responseBody.message/code`). Use `REACT_APP_API_BASE_URL` from the environment.
    * Implement exported functions: `verifyGoogleToken`, `getCurrentUser`, `logoutUser`, each calling `apiFetch` with the correct endpoint path (e.g., `/auth.v1.AuthService/VerifyGoogleToken`) and request/response types. These functions will interact with the backend's ConnectRPC service.
<!-- STOP: Please review and approve before proceeding to the next step. -->

2. **Implement Auth Context & Provider:**
    * Create `frontend/src/AuthProvider.tsx`.
    * Define the `AuthState` interface (user, isLoading, error, methods).
    * Create `AuthContext`.
    * Implement the `AuthProvider` component:
        * Use React hooks like `useState` and `useEffect` for managing authentication state (`user`, `isLoading`, `error`).
        * Implement `checkSession`: Calls `api.getCurrentUser` to check for an existing session, updates state accordingly, and handles errors (especially ignoring expected 'unauthenticated' errors).
        * Implement `loginUser`: Sets the user state after successful verification via the Google login button callback.
        * Implement `logout`: Calls `api.logoutUser` to invalidate the backend session, clears frontend user state, and calls `window.google.accounts.id.disableAutoSelect()` to prevent automatic sign-in. Handles errors during logout.
        * Use `useEffect` to call `checkSession` on initial component mount to check for a persistent session.
        * Provide the authentication state and methods via `AuthContext.Provider` to the rest of the application.
<!-- STOP: Please review and approve before proceeding to the next step. -->

3. **Implement Auth Hook:**
    * Create `frontend/src/hooks/useAuth.ts`.
    * Implement the `useAuth` hook that consumes `AuthContext` using `useContext` and throws an error if used outside the `AuthProvider`, ensuring the hook is used within the authentication context.
<!-- STOP: Please review and approve before proceeding to the next step. -->
