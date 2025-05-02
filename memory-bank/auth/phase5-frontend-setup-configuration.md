# Phase 5: Frontend Setup & Configuration

1. **Install Frontend Dependencies:**
    * Navigate into the `frontend/` directory (`cd frontend`).
    * Install the necessary Google types dependency: `npm install --save-dev @types/google.accounts` (or `yarn add --dev`).
<!-- STOP: Please review and approve before proceeding to the next step. -->

2. **Add Google Script to HTML:**
    * Edit `frontend/index.html`.
    * Add `<script src="https://accounts.google.com/gsi/client" async defer></script>` inside the `<head>`.
<!-- STOP: Please review and approve before proceeding to the next step. -->

3. **Create Frontend Environment File:**
    * Create a file named `frontend/.env` in the `frontend/` directory.
    * Add the following content, replacing placeholders:

        ```dotenv
        REACT_APP_GOOGLE_CLIENT_ID=YOUR_GOOGLE_CLIENT_ID_HERE # MUST match backend/Google Cloud
        REACT_APP_API_BASE_URL=/api # Use relative path for embedded deployment
        ```

    * *(Security Note: Ensure `frontend/.env` is in `.gitignore`)*.
<!-- STOP: Please review and approve before proceeding to the next step. -->

4. **Add TypeScript Global Types:**
    * Create a new file, for example `frontend/src/global.d.ts`.
    * Add the following content to declare the global `google` object:

        ```typescript
        declare global {
          interface Window {
            google?: typeof google;
          }
        }
        ```
<!-- STOP: Please review and approve before proceeding to the next step. -->
