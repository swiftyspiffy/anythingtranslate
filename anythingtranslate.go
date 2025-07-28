// Package anythingtranslate provides a client for the AnythingTranslate service.
//
// AnythingTranslate is a web service that offers various text translation options,
// including fun transformations like Warhammer 40k Ork speech, Pidgin English,
// and many others. This package allows programmatic access to these translators.
package anythingtranslate

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// AnythingTranslateClient is the main client for interacting with the AnythingTranslate service.
// It handles authentication, session management, and translation requests.
type AnythingTranslateClient struct {
	cookie           string
	cookieExpiration *time.Time
	nonces           map[string]string
	logger           Logger
}

// requestResponse represents the structure of API responses from the AnythingTranslate service.
type requestResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

// NewAnythingTranslateClient creates a new client with a no-op logger.
// This is a convenience function for users who don't need logging.
func NewAnythingTranslateClient() *AnythingTranslateClient {
	return NewAnythingTranslateClientWithLogger(NoopLogger{})
}

// NewAnythingTranslateClientWithLogger creates a new client with the specified logger.
// This allows users to provide their own logger implementation for custom logging.
func NewAnythingTranslateClientWithLogger(logger Logger) *AnythingTranslateClient {
	return &AnythingTranslateClient{
		nonces: make(map[string]string),
		logger: logger,
	}
}

// Translate sends the provided text to the AnythingTranslate service to be translated
// using the specified translator target.
//
// The target parameter should be one of the predefined TranslatorTarget constants.
// Returns the translated text or an error if the translation failed.
func (c *AnythingTranslateClient) Translate(text string, target TranslatorTarget) (string, error) {
	c.logger.Info("translating text", Field{"text", text}, Field{"target", string(target)})

	nonce, err := c.getNonce(target)
	if err != nil {
		return "", err
	}

	return c.request(text, nonce, string(target), false)
}

// TranslateWithCustomTarget sends the provided text to the AnythingTranslate service to be translated
// using a custom translator target ID.
//
// This method is useful when you have a translator ID that is not included in the predefined constants.
// Returns the translated text or an error if the translation failed.
func (c *AnythingTranslateClient) TranslateWithCustomTarget(text string, customTarget string) (string, error) {
	c.logger.Info("translating text with custom target", Field{"text", text}, Field{"custom target", customTarget})

	target := TranslatorTarget(customTarget)
	nonce, err := c.getNonce(target)
	if err != nil {
		return "", err
	}

	return c.request(text, nonce, string(target), false)
}

// getNonce retrieves an existing nonce for the given target or generates a new one if none exists.
// Nonces are used for authentication with the AnythingTranslate API.
func (c *AnythingTranslateClient) getNonce(target TranslatorTarget) (string, error) {
	if _, found := c.nonces[string(target)]; !found {
		newNonce, err := generateNonce(10)
		if err != nil {
			return "", err
		}
		c.nonces[string(target)] = newNonce
	}

	return c.nonces[string(target)], nil
}

// request sends a translation request to the AnythingTranslate API and processes the response.
// It handles authentication, retries, and error parsing.
//
// Parameters:
//   - input: the text to translate
//   - nonce: the authentication nonce for the request
//   - postId: the ID of the translator to use
//   - isRetry: whether this is a retry attempt to prevent infinite recursion
//
// Returns the translated text or an error if the translation failed.
func (c *AnythingTranslateClient) request(input string, nonce string, postId string, isRetry bool) (string, error) {
	apiURL := "https://anythingtranslate.com/wp-admin/admin-ajax.php"

	encodedText := url.QueryEscape(input)

	payload := "action=do_translation&translator_nonce=" + nonce + "&post_id=" + postId + "&to_translate=" + encodedText

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(payload))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	// if a cookie exists in the client, and its before the expiration, use it
	if c.cookie != "" && c.cookieExpiration != nil && time.Now().Before(*c.cookieExpiration) {
		req.Header.Set("Cookie", "f_t="+c.cookie)
	}
	req.Header.Set("Referer", "https://anythingtranslate.com/")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// if Set-Cookie header exists, always use it
	setCookieVal := resp.Header.Get("Set-Cookie")
	if setCookieVal != "" {
		// Extract cookie value safely
		ftParts := strings.Split(setCookieVal, "f_t=")
		if len(ftParts) < 2 {
			c.logger.Debug("cookie format unexpected", Field{"cookie", setCookieVal})
		} else {
			valueParts := strings.Split(ftParts[1], ";")
			if len(valueParts) > 0 {
				c.cookie = valueParts[0]
			}
		}

		// Extract expiration time safely
		expParts := strings.Split(setCookieVal, "expires=")
		if len(expParts) < 2 {
			c.logger.Debug("cookie expiration format unexpected", Field{"cookie", setCookieVal})
		} else {
			expValueParts := strings.Split(expParts[1], ";")
			if len(expValueParts) > 0 {
				cookieExp, err := parseCookieTimeIntoTime(expValueParts[0])
				if err != nil {
					c.logger.Debug("failed to parse cookie expiration time", Field{"error", err.Error()})
				} else {
					c.cookieExpiration = cookieExp
				}
			}
		}
	}

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	c.logger.Debug("raw api response", Field{"raw_response", string(resBody)})

	var res requestResponse
	if err = json.Unmarshal(resBody, &res); err != nil {
		return "", err
	}

	if !res.Success {
		// If we fail but have set the cookie, try again once
		// This avoids infinite recursion by checking isRetry flag
		if !isRetry {
			c.logger.Debug("translation failed, retrying with updated cookie")
			return c.request(input, nonce, postId, true)
		}

		// If we're already on a retry or the retry failed, return an appropriate error
		switch v := res.Data.(type) {
		case string:
			return "", fmt.Errorf("translation failed: %s", v)
		case map[string]interface{}:
			errorMsg := "unknown error"
			if msg, ok := v["message"].(string); ok {
				errorMsg = msg
			}
			errorCode := "unknown"
			if code, ok := v["error"].(string); ok {
				errorCode = code
			}
			return "", fmt.Errorf("translation failed: %s (code: %s)", errorMsg, errorCode)
		default:
			return "", fmt.Errorf("translation failed with unknown error format")
		}
	}

	translation, ok := res.Data.(string)
	if !ok {
		return "", fmt.Errorf("translation failed: unexpected data format in successful response")
	}

	return translation, nil
}

// parseCookieTimeIntoTime converts a cookie expiration time string into a time.Time object.
// The time string is expected to be in the format "Mon, 02-Jan-2006 15:04:05 MST".
func parseCookieTimeIntoTime(timeStr string) (*time.Time, error) {
	layout := "Mon, 02-Jan-2006 15:04:05 MST"
	parsedTime, err := time.Parse(layout, timeStr)
	if err != nil {
		return nil, err
	}

	return &parsedTime, nil
}

// generateNonce creates a random hexadecimal string of the specified length.
// The length must be even because each byte is represented by two hex characters.
func generateNonce(length int) (string, error) {
	if length%2 != 0 {
		return "", fmt.Errorf("length must be even, got %d", length)
	}
	b := make([]byte, length/2)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
