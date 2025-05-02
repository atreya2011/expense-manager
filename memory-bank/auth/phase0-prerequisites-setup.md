# Phase 0: Prerequisites & Setup

1. **Google Cloud Project Setup:**
    * Access the Google Cloud Console.
    * Create or select a Google Cloud project.
    * Navigate to "APIs & Services" -> "Credentials".
    * Configure the "OAuth consent screen" (User Type: External, App Name, Emails).
    * Create an "OAuth client ID" of type "Web application".
    * Add authorized JavaScript origins: `http://localhost:3000` (for local dev) and your production domain (e.g., `https://your-app.com`).
    * Add authorized redirect URIs (optional but good practice): `http://localhost:3000`, `https://your-app.com`.
    * Record the **Client ID**. This is `YOUR_GOOGLE_CLIENT_ID_HERE`.
    * *(Security Note: Client Secret is not typically used directly in this SPA flow).*
<!-- STOP: Please review and approve before proceeding to the next step. -->
