// Package membership implements the service.StudentClient driven port over
// HTTP — LoanRegistrationService can no longer call membership's
// StudentRepository directly now that Student lives in a different database
// (library-docs/09-microservices/service-boundary-rules.md).
package membership

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ErrStudentNotFound mirrors membership-service's 404 for a missing student.
var ErrStudentNotFound = fmt.Errorf("student not found")

type Client struct {
	baseURL   string
	jwtSecret string
	http      *http.Client
}

func NewClient(baseURL, jwtSecret string) *Client {
	return &Client{baseURL: baseURL, jwtSecret: jwtSecret, http: &http.Client{Timeout: 5 * time.Second}}
}

type studentResponse struct {
	ID             string  `json:"id"`
	SuspendedUntil *string `json:"suspendedUntil,omitempty"`
}

// ResolveByDocumentID translates the identifier an Administrator actually
// has on hand (the student's document ID) into the internal UUID the domain
// operates on — an Administrator has no way to know the UUID
// (library-docs/09-microservices/data-ownership-matrix.md).
func (c *Client) ResolveByDocumentID(ctx context.Context, documentID string) (string, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/students/by-document/%s", documentID), nil)
	if err != nil {
		return "", err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("calling membership-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", ErrStudentNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("membership-service returned %d", resp.StatusCode)
	}

	var body studentResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", fmt.Errorf("decoding membership-service response: %w", err)
	}
	return body.ID, nil
}

// IsEligible mirrors membership.Student.IsEligibleForLoan — a student with no
// active suspension, or one that already ended, is eligible.
func (c *Client) IsEligible(ctx context.Context, studentID string) (bool, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/students/%s", studentID), nil)
	if err != nil {
		return false, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return false, fmt.Errorf("calling membership-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, ErrStudentNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("membership-service returned %d", resp.StatusCode)
	}

	var body studentResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return false, fmt.Errorf("decoding membership-service response: %w", err)
	}
	if body.SuspendedUntil == nil {
		return true, nil
	}
	suspendedUntil, err := time.Parse(time.RFC3339, *body.SuspendedUntil)
	if err != nil {
		return false, fmt.Errorf("parsing suspendedUntil: %w", err)
	}
	return suspendedUntil.Before(time.Now().UTC()), nil
}

// Suspend applies the flat suspension (INV-006) via membership-service —
// the caller (LoanRegistrationService) decides *when* (a late return);
// membership-service owns *how* (mutating and persisting Student).
func (c *Client) Suspend(ctx context.Context, studentID string, days int) error {
	body, err := json.Marshal(map[string]int{"days": days})
	if err != nil {
		return err
	}

	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/api/v1/students/%s/suspend", studentID), bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("calling membership-service: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNotFound:
		return ErrStudentNotFound
	case http.StatusOK:
		return nil
	default:
		return fmt.Errorf("membership-service returned %d", resp.StatusCode)
	}
}

func (c *Client) newRequest(ctx context.Context, method, path string, body *bytes.Reader) (*http.Request, error) {
	token, err := c.mintServiceToken()
	if err != nil {
		return nil, fmt.Errorf("minting internal service token: %w", err)
	}

	var req *http.Request
	if body != nil {
		req, err = http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	} else {
		req, err = http.NewRequestWithContext(ctx, method, c.baseURL+path, nil)
	}
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	return req, nil
}

// mintServiceToken signs a short-lived internal token with the JWT_SECRET
// shared by every service (see access-service's JWTIssuer) — v1 has no
// separate service-to-service auth scope, an accepted trade-off at this size.
func (c *Client) mintServiceToken() (string, error) {
	now := time.Now().UTC()
	claims := jwt.MapClaims{
		"sub": "circulation-service",
		"iat": now.Unix(),
		"exp": now.Add(1 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(c.jwtSecret))
}
