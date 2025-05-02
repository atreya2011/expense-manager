# Phase 2: Backend API Definition & Code Generation

1. **Define Protobuf Service:**
    * Create the file `proto/auth/v1/auth.proto`.
    * Paste the Protobuf definition provided in the example (`AuthService` with `VerifyGoogleToken`, `GetCurrentUser`, `Logout` RPCs and `User` message).
    * **Crucially:** Adjust the `go_package` option in the proto file to match the current project's Go module path (`github.com/atreya2011/expense-manager`) and the desired output directory (`internal/rpc/gen/auth/v1`): `option go_package = "github.com/atreya2011/expense-manager/internal/rpc/gen/auth/v1;authv1";`.
<!-- STOP: Please review and approve before proceeding to the next step. -->

2. **Generate Go Code from Proto:**
    * Run `buf generate` in the root directory.
    * Verify that the `internal/rpc/gen/auth/v1/` directory and its contents (`auth.pb.go`, `authv1connect/auth.connect.go`) are created.
<!-- STOP: Please review and approve before proceeding to the next step. -->
