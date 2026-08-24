package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/WhiteKiwi/pushman-cli/internal/cli"
	"github.com/WhiteKiwi/pushman-cli/internal/credential"
)

type memoryCredentials struct{ token string }

func (m *memoryCredentials) Get() (string, error) {
	if m.token == "" {
		return "", credential.ErrNotFound
	}
	return m.token, nil
}
func (m *memoryCredentials) Set(value string) error { m.token = value; return nil }
func (m *memoryCredentials) Delete() error          { m.token = ""; return nil }

func TestPairStoresCredentialOnlyAfterApproval(t *testing.T) {
	expiresAt := time.Now().Add(time.Minute).UTC().Truncate(time.Second)
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		switch {
		case request.Method == http.MethodPost && request.URL.Path == "/v1/pairings":
			response.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(response).Encode(map[string]any{
				"id": "pair_01K00000000000000000000000", "userCode": "ABCD-EFGH",
				"verificationUri": "https://app.pushman.example/pair", "expiresAt": expiresAt,
				"interval": 5, "status": "pending", "pairingSecret": "pair-secret",
			})
		case request.Method == http.MethodGet && request.URL.Path == "/v1/pairings/pair_01K00000000000000000000000":
			if request.Header.Get("pairingSecret") != "pair-secret" {
				t.Fatal("missing pairing secret header")
			}
			_ = json.NewEncoder(response).Encode(map[string]any{
				"id": "pair_01K00000000000000000000000", "userCode": "ABCD-EFGH",
				"verificationUri": "https://app.pushman.example/pair", "expiresAt": expiresAt,
				"interval": 5, "status": "approved",
				"credential": map[string]string{"token": "pm_cli_secret", "name": "Build Mac"},
			})
		default:
			http.NotFound(response, request)
		}
	}))
	defer server.Close()

	credentials := new(memoryCredentials)
	service, err := New(server.URL+"/v1", credentials, "", server.Client())
	if err != nil {
		t.Fatal(err)
	}
	service.wait = func(context.Context, time.Duration) error { return nil }
	var challenge cli.PairChallenge
	result, err := service.Pair(context.Background(), cli.PairRequest{
		Platform: "darwin", SuggestedName: "Build Mac",
		OnChallenge: func(value cli.PairChallenge) error { challenge = value; return nil },
	})
	if err != nil {
		t.Fatal(err)
	}
	if challenge.UserCode != "ABCD-EFGH" || result.Nickname != "Build Mac" || credentials.token != "pm_cli_secret" {
		t.Fatalf("challenge/result/token = %#v / %#v / %q", challenge, result, credentials.token)
	}
}

func TestPushPrefersAutomationToken(t *testing.T) {
	var authorization string
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		authorization = request.Header.Get("Authorization")
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusAccepted)
		_, _ = response.Write([]byte(`{"id":"msg_01K00000000000000000000000","logicalMessageId":"lmsg_01K00000000000000000000000","acceptedAt":"2026-08-25T00:00:00Z","targetCount":1}`))
	}))
	defer server.Close()
	credentials := &memoryCredentials{token: "pm_cli_paired"}
	service, err := New(server.URL+"/v1", credentials, "pm_auto_token", server.Client())
	if err != nil {
		t.Fatal(err)
	}
	result, err := service.Push(context.Background(), cli.PushRequest{Body: "done", Sound: "default", Format: "plain"})
	if err != nil {
		t.Fatal(err)
	}
	if authorization != "Bearer pm_auto_token" || result.Status != "accepted" {
		t.Fatalf("authorization/result = %q / %#v", authorization, result)
	}
}

func TestUnauthorizedResponseMapsToAuthorizationExit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusUnauthorized)
		_, _ = response.Write([]byte(`{"error":{"code":"unauthorized","message":"revoked","requestId":"req_test"}}`))
	}))
	defer server.Close()
	service, err := New(server.URL+"/v1", &memoryCredentials{token: "revoked"}, "", server.Client())
	if err != nil {
		t.Fatal(err)
	}
	_, err = service.Status(context.Background())
	var authorization *cli.AuthorizationError
	if !errors.As(err, &authorization) || cli.ExitCode(err) != 4 {
		t.Fatalf("error/exit = %T %v / %d", err, err, cli.ExitCode(err))
	}
}

func TestValidateBaseURL(t *testing.T) {
	tests := []struct {
		value string
		ok    bool
	}{
		{"https://api.pushman.example/v1", true},
		{"http://127.0.0.1:8080/v1", true},
		{"http://localhost:8080/v1/", true},
		{"http://api.pushman.example/v1", false},
		{"https://api.pushman.example/v2", false},
		{"https://user@example.com/v1", false},
		{"https://api.pushman.example/v1?token=secret", false},
	}
	for _, test := range tests {
		t.Run(strings.ReplaceAll(test.value, "/", "_"), func(t *testing.T) {
			_, err := ValidateBaseURL(test.value)
			if (err == nil) != test.ok {
				t.Fatalf("ValidateBaseURL(%q) error = %v", test.value, err)
			}
		})
	}
}
