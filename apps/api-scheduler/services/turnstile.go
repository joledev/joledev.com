package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

var turnstileClient = &http.Client{Timeout: 10 * time.Second}

func VerifyTurnstile(token, remoteIP string) error {
	secret := os.Getenv("TURNSTILE_SECRET_KEY")
	if secret == "" {
		return nil // Skip in dev (no key configured)
	}
	if token == "" {
		return fmt.Errorf("CAPTCHA verification required")
	}
	resp, err := turnstileClient.PostForm(
		"https://challenges.cloudflare.com/turnstile/v0/siteverify",
		url.Values{
			"secret":   {secret},
			"response": {token},
			"remoteip": {remoteIP},
		},
	)
	if err != nil {
		return fmt.Errorf("CAPTCHA verification failed")
	}
	defer resp.Body.Close()
	var result struct {
		Success bool `json:"success"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil || !result.Success {
		return fmt.Errorf("CAPTCHA verification failed")
	}
	return nil
}
