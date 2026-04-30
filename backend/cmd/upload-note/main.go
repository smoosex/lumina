package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/smoosex/lumina/backend/internal/config"
	"github.com/smoosex/lumina/backend/internal/markdown"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	filePath := flag.String("file", "", "markdown file to upload")
	slugFlag := flag.String("slug", "", "override note slug")
	baseURLFlag := flag.String("base-url", "", "API base URL")
	flag.Parse()

	if strings.TrimSpace(*filePath) == "" {
		return errors.New("-file is required")
	}

	_ = godotenv.Load(".env", "backend/.env")

	secret := strings.TrimSpace(os.Getenv("UPLOAD_SIGNATURE_SECRET"))
	if secret == "" {
		return errors.New("UPLOAD_SIGNATURE_SECRET is required")
	}

	baseURL := firstNonEmpty(*baseURLFlag, os.Getenv("UPLOAD_API_BASE_URL"))
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	body, err := os.ReadFile(*filePath)
	if err != nil {
		return fmt.Errorf("read markdown file: %w", err)
	}

	slug := strings.TrimSpace(*slugFlag)
	if slug == "" {
		parser := markdown.New(config.Config{
			Markdown: config.MarkdownConfig{
				RenderHTML:   false,
				SanitizeHTML: false,
			},
		})
		doc, err := parser.Parse(body)
		if err != nil {
			return fmt.Errorf("parse markdown: %w", err)
		}
		slug = doc.Meta.Slug
	}

	endpoint, err := url.JoinPath(baseURL, "/api/admin/reading-notes", slug)
	if err != nil {
		return fmt.Errorf("build upload URL: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create upload request: %w", err)
	}
	req.Header.Set("Content-Type", "text/markdown; charset=utf-8")
	signRequest(req, body, secret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("upload note: %w", err)
	}
	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return fmt.Errorf("upload failed: %s\n%s", resp.Status, strings.TrimSpace(string(responseBody)))
	}

	fmt.Printf("uploaded %s to %s\n", slug, endpoint)
	return nil
}

func signRequest(req *http.Request, body []byte, secret string) {
	bodySum := sha256.Sum256(body)
	bodyHash := hex.EncodeToString(bodySum[:])
	timestamp := fmt.Sprint(time.Now().Unix())

	payload := fmt.Sprintf(
		"%s\n%s\n%s\n%s",
		req.Method,
		req.URL.Path,
		timestamp,
		bodyHash,
	)

	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(payload))

	req.Header.Set("X-Lumina-Timestamp", timestamp)
	req.Header.Set("X-Lumina-Body-SHA256", bodyHash)
	req.Header.Set("X-Lumina-Signature", hex.EncodeToString(mac.Sum(nil)))
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}
