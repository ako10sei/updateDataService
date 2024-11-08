package token

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"log/slog"

	"github.com/stretchr/testify/assert"
)

func TestDigitalProfileTokenProvider_GetToken(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(Token{
			AccessToken: "mock-access-token",
			ExpiresIn:   3600,
			TokenType:   "Bearer",
			Scope:       "digital_profile",
		})
		if err != nil {
			return
		}
	}))
	defer mockServer.Close()

	provider := &DigitalProfileTokenProvider{
		URL:          mockServer.URL,
		ClientID:     "mock-client-id",
		ClientSecret: "mock-client-secret",
		Scope:        "digital_profile",
	}

	log := slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))

	token, err := provider.GetToken(log)
	assert.NoError(t, err)
	assert.Equal(t, "mock-access-token", token)
}

func TestVisiologyTokenProvider_GetToken(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(Token{
			AccessToken: "mock-visiology-token",
			ExpiresIn:   3600,
			TokenType:   "Bearer",
			Scope:       "openid profile email roles",
		})
		if err != nil {
			return
		}
	}))
	defer mockServer.Close()

	provider := &VisiologyTokenProvider{
		URL:      mockServer.URL,
		Username: "mock-username",
		Password: "mock-password",
		Scope:    "openid profile email roles",
	}

	log := slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))

	token, err := provider.GetToken(log)
	assert.NoError(t, err)
	assert.Equal(t, "mock-visiology-token", token)
}

func TestRequestToken_NonOKResponse(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	log := slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))

	_, err := requestToken(mockServer.URL, "mock-body", log, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "неверный статус ответа: 500")
}

func TestRequestToken_InvalidResponseBody(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("invalid json"))
		if err != nil {
			return
		}
	}))
	defer mockServer.Close()

	log := slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))

	_, err := requestToken(mockServer.URL, "mock-body", log, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ошибка десериализации тела ответа")
}

type mockReadCloser struct {
	closeError error
}

func (m *mockReadCloser) Read(_ []byte) (int, error) {
	return 0, nil
}

func (m *mockReadCloser) Close() error {
	return m.closeError
}

func TestCloseResponse_Error(t *testing.T) {
	log := slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))

	mockBody := &mockReadCloser{closeError: assert.AnError}

	closeResponse(mockBody, log)
}

func TestHandleNonOKResponse(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       io.NopCloser(bytes.NewBufferString("bad request")),
	}

	log := slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))

	handleNonOKResponse(resp, log)
}
