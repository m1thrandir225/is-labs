 Authentication Go + Gin
Built using Go, Gin, SQLC, SQLite and HOTP. 

## How to Run locally
1. Create a local .db file 
2. Setup the file `app.env`
3. Install go-migrate and run the migrations
4. Run `make server` to start the http server locally.

### How Does It Work?

When a user registers, two key actions occur:

1. **Password Hashing and Secret Key Creation**
   - The user's password is hashed using bcrypt and stored in the database.
   - A unique secret key is generated using Base32, which is used alongside a counter to create OTP codes.

2. **Counter Initialization**
   - A relationship is established with the `hotp_counter` table to track the current counter for generating HMAC-based one-time passwords (HOTPs).

After successful registration, the user can proceed to log in. Depending on whether two-factor authentication (2FA) is enabled during registration, the user either logs in directly or is prompted for additional 2FA verification.

#### Login Process
1. Verify the user exists using their email. If not, an error is returned.
2. Compare the submitted password with the stored hashed password. If they do not match, an error is returned.
3. If 2FA is enabled:
   - Retrieve the user's associated counter.
   - Generate an OTP using the counter and the user's secret key.
   - Email the OTP to the user.
   - Return a response indicating the OTP was sent.
4. If 2FA is not enabled:
   - Generate an access token and a refresh token with different expiration times.
   - Send the tokens to the user.

#### OTP Verification
For OTP verification, the process is similar to login but includes:
- Sending the `user_id` (provided during registration) and the received OTP code.
- Retrieving the user, incrementing the counter, and verifying the OTP.
- If the OTP is valid, generate and send access and refresh tokens. Otherwise, return an error.

#### Refresh Token
The `/refresh-token` route allows the user to refresh their access token, which has a shorter expiration time.

### Available Routes
The application provides the following routes:
- `/register`: Register a new user.
- `/login`: Log in as an existing user.
- `/verify-2fa`: Verify a 2FA OTP.
- `/refresh-token`: Refresh the access token.
- `/`: A protected route accessible only with valid authentication.
