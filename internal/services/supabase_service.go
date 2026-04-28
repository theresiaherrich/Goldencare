package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/theresiaherrich/Goldencare/internal/config"
)

type SupabaseService struct {
	url        string
	anonKey    string
	serviceKey string
}

func NewSupabaseService(cfg *config.Config) *SupabaseService {
	return &SupabaseService{
		url:        cfg.SupabaseURL,
		anonKey:    cfg.SupabaseAnonKey,
		serviceKey: cfg.SupabaseServiceKey,
	}
}

type SignedURLRequest struct {
	ExpiresIn int `json:"expiresIn"`
}

type SignedURLResponse struct {
	SignedURL string `json:"signedURL"`
}

func (s *SupabaseService) GenerateSignedURL(bucket, filename string, expirySeconds int) (string, error) {
	if expirySeconds == 0 {
		expirySeconds = 3600
	}

	pathName := fmt.Sprintf("%s/%s", bucket, filename)
	encodedPath := url.QueryEscape(pathName)

	endpoint := fmt.Sprintf("%s/storage/v1/object/sign/%s", s.url, encodedPath)

	body := SignedURLRequest{
		ExpiresIn: expirySeconds,
	}
	jsonBody, _ := json.Marshal(body)

	req := &http.Request{
		Method: "POST",
		Header: http.Header{
			"Authorization": []string{fmt.Sprintf("Bearer %s", s.anonKey)},
			"Content-Type":  []string{"application/json"},
		},
		Body: io.NopCloser(bytes.NewReader(jsonBody)),
	}
	req.URL, _ = url.Parse(endpoint)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result SignedURLResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%s", s.url, result.SignedURL), nil
}

func (s *SupabaseService) UploadFile(bucket, path string, fileData []byte, contentType string) (string, error) {
	endpoint := fmt.Sprintf("%s/storage/v1/object/%s/%s", s.url, bucket, path)

	req := &http.Request{
		Method: "POST",
		Header: http.Header{
			"Authorization": []string{fmt.Sprintf("Bearer %s", s.serviceKey)},
			"Content-Type":  []string{contentType},
		},
		Body: io.NopCloser(bytes.NewReader(fileData)),
	}
	req.URL, _ = url.Parse(endpoint)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload failed: %d %s", resp.StatusCode, string(body))
	}

	return s.GetPublicURL(bucket, path), nil
}

func (s *SupabaseService) DeleteFile(bucket, filename string) error {
	pathName := fmt.Sprintf("%s/%s", bucket, filename)
	endpoint := fmt.Sprintf("%s/storage/v1/object/%s", s.url, pathName)

	req := &http.Request{
		Method: "DELETE",
		Header: http.Header{
			"Authorization": []string{fmt.Sprintf("Bearer %s", s.serviceKey)},
		},
	}
	req.URL, _ = url.Parse(endpoint)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("failed to delete file: %d", resp.StatusCode)
	}

	return nil
}

func (s *SupabaseService) GetPublicURL(bucket, filename string) string {
	return fmt.Sprintf("%s/storage/v1/object/public/%s/%s", s.url, bucket, filename)
}

func (s *SupabaseService) VerifySupabaseSignature(token string) (bool, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false, fmt.Errorf("invalid token format")
	}

	message := strings.Join(parts[:2], ".")
	signature := parts[2]

	decodedKey, err := base64.RawURLEncoding.DecodeString(s.serviceKey)
	if err != nil {
		return false, err
	}

	h := hmac.New(sha256.New, decodedKey)
	h.Write([]byte(message))
	expectedSignature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	return signature == expectedSignature, nil
}

type SupabaseToken struct {
	Sub string `json:"sub"`
	Aud string `json:"aud"`
	Exp int64  `json:"exp"`
}

func (s *SupabaseService) ParseToken(token string) (*SupabaseToken, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	var claims SupabaseToken
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, err
	}

	if claims.Exp < time.Now().Unix() {
		return nil, fmt.Errorf("token expired")
	}

	return &claims, nil
}
