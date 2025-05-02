# Phase 7: Frontend UI Implementation

1. **Wrap App in AuthProvider:**
    * Edit `frontend/src/main.tsx`. Note that the main entry point is `main.tsx` in this project, not `index.tsx`.
    * Import `AuthProvider`.
    * Wrap the main application component (likely the router provider) with `<AuthProvider>`.
<!-- STOP: Please review and approve before proceeding to the next step. -->

2. **Implement Google Login Button Component:**
    * Create `frontend/src/components/GoogleLoginButton.tsx`.
    * Use `useRef` for the button's container div.
    * Use `useState` for `isScriptLoaded` and local `buttonError`.
    * Use the `useAuth` hook to get `loginUser`, `isLoading`, `error`.
    * Implement the `handleCredentialResponse` callback:
        * Calls `api.verifyGoogleToken` with the received `response.credential`.
        * On success, calls `loginUser` from the auth context.
        * On failure, sets the local `buttonError`.
    * Use `useEffect` to check for the Google script load (`isScriptLoaded` state).
    * Use `useEffect` (dependent on `isScriptLoaded`) to:
        * Initialize `window.google.accounts.id.initialize` with the `REACT_APP_GOOGLE_CLIENT_ID` and `handleCredentialResponse` callback.
        * Render the button using `window.google.accounts.id.renderButton`.
    * Render the placeholder div with the `ref`. Display loading/error states from the auth context or local button state.
<!-- STOP: Please review and approve before proceeding to the next step. -->

3. **Integrate Auth State into Main App Component or Router:**
    * Modify `frontend/src/App.tsx` or the relevant routing component (e.g., `frontend/src/routes/__root.tsx`) to use the `useAuth` hook.
    * Conditionally render content based on the authentication state (`user`, `isLoading`, `error`):
        * Show a loading indicator if `isLoading`.
        * If not loading and no `user`, display the `GoogleLoginButton` component.
        * If not loading and `user` exists, display user-specific content (e.g., welcome message, user details) and a Logout button.
        * Make the Logout button's `onClick` call the `logout` function from the auth context.
<!-- STOP: Please review and approve before proceeding to the next step. -->
