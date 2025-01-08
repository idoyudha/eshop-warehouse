package v1

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/idoyudha/eshop-warehouse/config"
	"github.com/stretchr/testify/assert"
)

func TestCognitoMiddleware(t *testing.T) {
	// t.Parallell()

	mockUserID := uuid.New()

	var mockServer *httptest.Server
	mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		token := strings.TrimPrefix(r.URL.Path, "/v1/auth/")

		switch token {
		case "valid_token":
			response := authSuccessResponse{
				Code: http.StatusOK,
				Data: authResponse{
					UserID: mockUserID,
					Role:   "user",
				},
				Message: "success",
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		case "invalid_token":
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"message": "unauthorized",
			})
		default:
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"message": "unauthorized",
			})
		}
	}))

	t.Cleanup(func() {
		mockServer.Close()
	})

	tests := []struct {
		name           string
		authHeader     string
		expectedCode   int
		expectedUserID *uuid.UUID
	}{
		{
			name:           "success - valid token",
			authHeader:     "Bearer valid_token",
			expectedCode:   http.StatusOK,
			expectedUserID: &mockUserID,
		},
		{
			name:           "error - no auth header",
			authHeader:     "",
			expectedCode:   http.StatusUnauthorized,
			expectedUserID: nil,
		},
		{
			name:           "error - invalid token",
			authHeader:     "Bearer invalid_token",
			expectedCode:   http.StatusUnauthorized,
			expectedUserID: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			authConfig := config.AuthService{
				BaseURL: mockServer.URL,
			}

			router := gin.New()

			var capturedUserID uuid.UUID
			router.Use(cognitoMiddleware(authConfig))
			router.GET("/test", func(c *gin.Context) {
				if id, exists := c.Get(UserIDKey); exists {
					if uid, ok := id.(uuid.UUID); ok {
						capturedUserID = uid
					}
				}
				c.Status(http.StatusOK)
			})

			// create request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code, "status code mismatch")

			if tt.expectedUserID != nil {
				assert.Equal(t, *tt.expectedUserID, capturedUserID, "userID mismatch")
			} else {
				assert.Equal(t, uuid.UUID{}, capturedUserID, "expected empty UUID")
			}
		})
	}
}
