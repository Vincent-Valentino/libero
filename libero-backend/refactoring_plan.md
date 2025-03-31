# Plan: AuthService Implementation (JWT & Password Login)

This document outlines the plan to implement `AuthService` correctly, focusing on password-based login, JWT generation/validation, and integrating with the existing OAuth flow.

## 1. Dependencies & Setup

*   **Add JWT Library:** Ensure `github.com/golang-jwt/jwt/v5` is added to `go.mod`.
*   **Import JWT Library:** Import `github.com/golang-jwt/jwt/v5` in `internal/service/auth_service.go`.
*   **Inject Config:** Modify `NewAuthService` in `auth_service.go` to accept `config.JWTConfig` as a parameter and store it in the `authService` struct.
*   **Custom Claims:** Define a custom struct (e.g., `JWTClaims`) within `auth_service.go` that embeds `jwt.RegisteredClaims` and includes `UserID`, `Email`, and `Role`.

## 2. Interface (`AuthService`)

Update the interface definition in `internal/service/auth_service.go`:

*   **Add:** `LoginByPassword(ctx context.Context, email, password string) (string, error)` - Returns JWT string.
*   **Add:** `ValidateJWTToken(tokenString string) (*JWTClaims, error)` - Returns custom claims or error.
*   **Modify:** `LoginOrRegisterViaProvider(ctx context.Context, userInfo *UserInfo) (string, error)` - Ensure it returns a JWT string upon success.

## 3. Implementation (`authService`)

Implement the methods within the `authService` struct in `internal/service/auth_service.go`:

*   **`generateJWTToken` (Private Helper):**
    *   Input: `*models.User`.
    *   Create `JWTClaims` instance.
    *   Set standard claims (Issuer, Subject, Audience - optional).
    *   Set expiration using `time.Now().Add(time.Second * time.Duration(s.jwtCfg.ExpiresIn))`.
    *   Set custom claims (`UserID`, `Email`, `Role`) from the user object.
    *   Create token: `jwt.NewWithClaims(jwt.SigningMethodHS256, claims)`.
    *   Sign token: `token.SignedString([]byte(s.jwtCfg.Secret))`.
    *   Return signed token string and nil error, or empty string and error.
*   **`LoginByPassword`:**
    *   Call `s.userService.FindUserByEmail(ctx, email)`.
    *   Handle user not found error.
    *   Call `user.ComparePassword(password)`. Handle mismatch error.
    *   If valid, call private `generateJWTToken(user)`.
    *   Return token string or error.
*   **`LoginOrRegisterViaProvider`:**
    *   Implement user lookup/creation/linking using `s.userService`.
    *   Replace placeholder return strings with calls to private `generateJWTToken(foundOrCreatedUser)`.
*   **`ValidateJWTToken`:**
    *   Parse using `jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) { return []byte(s.jwtCfg.Secret), nil })`.
    *   Check for parsing errors and claim type assertion errors.
    *   Return validated `*JWTClaims` or appropriate error.

## 4. Middleware (`internal/middleware/auth.go` - Related Task)

*   The auth middleware needs to be updated (or created) to:
    *   Extract the "Bearer <token>" from the `Authorization` header.
    *   Call `authService.ValidateJWTToken`.
    *   Inject `AuthService` dependency into the middleware.
    *   Store valid claims in the request context.
    *   Respond with `401 Unauthorized` on failure.

## 5. Conceptual Flow Diagram (JWT)

```mermaid
sequenceDiagram
    participant Client
    participant Controller
    participant Middleware
    participant AuthService
    participant UserService
    participant User Model

    Client->>Controller: POST /login (email, pass)
    Controller->>AuthService: LoginByPassword(email, pass)
    AuthService->>UserService: FindUserByEmail(email)
    UserService-->>AuthService: User or Error
    alt User Found
        AuthService->>User Model: ComparePassword(pass)
        User Model-->>AuthService: bool (match)
        alt Password Match
            AuthService->>AuthService: generateJWTToken(user)
            AuthService-->>Controller: JWT Token String
            Controller-->>Client: 200 OK {token: "..."}
        else Password Mismatch
            AuthService-->>Controller: Error (Invalid Credentials)
            Controller-->>Client: 401 Unauthorized
        end
    else User Not Found
        AuthService-->>Controller: Error (User Not Found)
        Controller-->>Client: 401 Unauthorized
    end

    Client->>Controller: GET /profile (Authorization: Bearer <token>)
    Controller->>Middleware: AuthMiddleware(request)
    Middleware->>AuthService: ValidateJWTToken(token)
    AuthService-->>Middleware: Claims or Error
    alt Token Valid
        Middleware->>Middleware: Add Claims to Context
        Middleware->>Controller: Next()
        Controller->>Controller: Get Claims from Context
        Controller->>UserService: GetUserByID(claims.UserID)
        UserService-->>Controller: User Profile
        Controller-->>Client: 200 OK {user profile}
    else Token Invalid
        Middleware-->>Client: 401 Unauthorized
    end