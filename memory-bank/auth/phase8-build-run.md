# Phase 8: Build & Run

1. **Build Frontend Artifacts:**
    * Navigate to the `frontend/` directory.
    * Run `npm run build` (or `yarn build`). This must complete successfully, creating the `frontend/dist` directory.
<!-- STOP: Please review and approve before proceeding to the next step. -->

2. **Build Go Binary:**
    * Navigate back to the root directory of the project.
    * Run `go build -o expense-manager .` (or `go build -o expense-manager.exe .` on Windows). This embeds the `frontend/dist` contents. Note that the output binary name is `expense-manager` in this project.
<!-- STOP: Please review and approve before proceeding to the next step. -->

3. **Run Application:**
    * Ensure the `.env` file (with backend variables) is present in the root directory.
    * Run the compiled binary: `./expense-manager` (or `.\expense-manager.exe`).
    * Access the application in your browser at `http://localhost:8080` (or the port configured in `.env`).
    * Test the Google Sign-In flow and logout functionality.
<!-- STOP: Please review and approve before proceeding to the next step. -->
