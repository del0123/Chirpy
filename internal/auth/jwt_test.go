package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	newUserID := uuid.New()
	secret := "123blabla"
	badSecret := "123blabla2"
	expiration := 1 * time.Hour
	badExpiration := -1 * time.Hour

	testCases := []struct {
		testName        string
		userID          uuid.UUID
		tokenSecret     string
		validatedSecret string
		expiresIn       time.Duration
		expectError     bool
	}{
		{
			testName:        "Valid JWT creation",
			userID:          newUserID,
			tokenSecret:     secret,
			validatedSecret: secret,
			expiresIn:       expiration,
			expectError:     false,
		},
		{
			testName:        "Expired JWT",
			userID:          newUserID,
			tokenSecret:     secret,
			validatedSecret: secret,
			expiresIn:       badExpiration,
			expectError:     true,
		},
		{
			testName:        "Invalid secret",
			userID:          newUserID,
			tokenSecret:     secret,
			validatedSecret: badSecret,
			expiresIn:       expiration,
			expectError:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			signedString, err := MakeJWT(tc.userID, tc.tokenSecret)
			if err != nil {
				t.Fatalf("Failed to make JWT: %v", err)
			}

			returnedUserID, err := ValidateJWT(signedString, tc.validatedSecret)
			// Scenario A: We expected an error (e.g., Expired or Bad Secret)
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error, but got none")
				}
				// Exit early! We don't check returnedUserID when validation was supposed to fail.
				return
			}

			// Scenario B: We expected success
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if returnedUserID != tc.userID {
				t.Errorf("Expected userID %v, got %v", tc.userID, returnedUserID)
			}

		})
	}
}
