# Phase 1: Backend Setup & Configuration

1. **Install Go Dependencies:**
    * Install the necessary dependencies that are not already present:

        ```bash
        go get golang.org/x/net/http2 golang.org/x/net/http2/h2c github.com/rs/cors google.golang.org/api/idtoken github.com/gorilla/sessions github.com/google/uuid
        ```

    * Install the required protobuf generation tools:

        ```bash
        go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
        go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest
        ```
<!-- STOP: Please review and approve before proceeding to the next step. -->

2. **Update Backend Environment File:**
    * Add the following variables to the existing `.env` file in the root directory, replacing placeholders:

        ```dotenv
        GOOGLE_CLIENT_ID=YOUR_GOOGLE_CLIENT_ID_HERE
        SESSION_SECRET_KEY=YOUR_STRONG_RANDOM_SESSION_SECRET_KEY # Generate using 'openssl rand -base64 32'
        SESSION_SECURE_COOKIE=false # Use 'true' for production HTTPS
        CORS_ALLOWED_ORIGIN=http://localhost:3000 # Adjust for production if needed, or keep for dev flexibility
        ```
<!-- STOP: Please review and approve before proceeding to the next step. -->
