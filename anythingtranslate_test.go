package anythingtranslate

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type mockTransport struct {
	mockServer    *httptest.Server
	origTransport http.RoundTripper
}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if strings.Contains(req.URL.Host, "anythingtranslate.com") {
		newURL := t.mockServer.URL + req.URL.Path

		newReq, err := http.NewRequest(req.Method, newURL, req.Body)
		if err != nil {
			return nil, err
		}

		newReq.Header = req.Header

		return t.origTransport.RoundTrip(newReq)
	}

	return t.origTransport.RoundTrip(req)
}

func TestTranslate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			t.Fatalf("Failed to parse form: %v", err)
		}

		postID := r.FormValue("post_id")
		toTranslate := r.FormValue("to_translate")

		var translation string
		switch postID {
		case string(Target_Warhammer40kOrk):
			translation = fmt.Sprintf("ORK TRANSLATION: %s", toTranslate)
		case string(Target_PidginEnglish):
			translation = fmt.Sprintf("PIDGIN TRANSLATION: %s", toTranslate)
		default:
			translation = fmt.Sprintf("UNKNOWN TRANSLATION: %s", toTranslate)
		}

		expireTime := time.Now().Add(24 * time.Hour)
		expireTimeStr := expireTime.Format("Mon, 02-Jan-2006 15:04:05 MST")
		w.Header().Set("Set-Cookie", "f_t=test-cookie-value; Path=/; expires="+expireTimeStr)

		response := requestResponse{
			Success: true,
			Data:    translation,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	originalTransport := http.DefaultTransport

	http.DefaultTransport = &mockTransport{
		mockServer:    server,
		origTransport: originalTransport,
	}

	defer func() {
		http.DefaultTransport = originalTransport
	}()

	testCases := []struct {
		name     string
		text     string
		target   TranslatorTarget
		expected string
	}{
		{
			name:     "Ork Translation",
			text:     "Hello world",
			target:   Target_Warhammer40kOrk,
			expected: "ORK TRANSLATION: Hello world",
		},
		{
			name:     "Pidgin Translation",
			text:     "How are you",
			target:   Target_PidginEnglish,
			expected: "PIDGIN TRANSLATION: How are you",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := NewAnythingTranslateClient()

			if client.nonces == nil {
				client.nonces = make(map[string]string)
			}

			client.nonces[string(tc.target)] = "test-nonce"

			result, err := client.Translate(tc.text, tc.target)
			if err != nil {
				t.Fatalf("Translate returned an error: %v", err)
			}

			if result != tc.expected {
				t.Errorf("Expected translation '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestTranslateWithCustomTarget(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			t.Fatalf("Failed to parse form: %v", err)
		}

		postID := r.FormValue("post_id")
		toTranslate := r.FormValue("to_translate")

		translation := fmt.Sprintf("CUSTOM TRANSLATION (%s): %s", postID, toTranslate)

		expireTime := time.Now().Add(24 * time.Hour)
		expireTimeStr := expireTime.Format("Mon, 02-Jan-2006 15:04:05 MST")
		w.Header().Set("Set-Cookie", "f_t=test-cookie-value; Path=/; expires="+expireTimeStr)

		response := requestResponse{
			Success: true,
			Data:    translation,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	originalTransport := http.DefaultTransport

	http.DefaultTransport = &mockTransport{
		mockServer:    server,
		origTransport: originalTransport,
	}

	defer func() {
		http.DefaultTransport = originalTransport
	}()

	client := NewAnythingTranslateClient()

	if client.nonces == nil {
		client.nonces = make(map[string]string)
	}

	customTarget := "12345"
	client.nonces[customTarget] = "test-nonce"

	result, err := client.TranslateWithCustomTarget("Custom text", customTarget)
	if err != nil {
		t.Fatalf("TranslateWithCustomTarget returned an error: %v", err)
	}

	expected := "CUSTOM TRANSLATION (12345): Custom text"
	if result != expected {
		t.Errorf("Expected translation '%s', got '%s'", expected, result)
	}
}

func TestTranslateError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expireTime := time.Now().Add(24 * time.Hour)
		expireTimeStr := expireTime.Format("Mon, 02-Jan-2006 15:04:05 MST")
		w.Header().Set("Set-Cookie", "f_t=test-cookie-value; Path=/; expires="+expireTimeStr)

		response := requestResponse{
			Success: false,
			Data:    "Error translating text",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	originalTransport := http.DefaultTransport

	http.DefaultTransport = &mockTransport{
		mockServer:    server,
		origTransport: originalTransport,
	}

	defer func() {
		http.DefaultTransport = originalTransport
	}()

	client := NewAnythingTranslateClient()

	if client.nonces == nil {
		client.nonces = make(map[string]string)
	}

	client.nonces[string(Target_Warhammer40kOrk)] = "test-nonce"

	_, err := client.Translate("Hello world", Target_Warhammer40kOrk)

	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
}
