package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// WebhookPayload represents the GitHub webhook payload structure
type WebhookPayload struct {
	Ref        string     `json:"ref"`
	Repository Repository `json:"repository"`
	Pusher     Pusher     `json:"pusher"`
	HeadCommit HeadCommit `json:"head_commit"`
	Deployment Deployment `json:"deployment,omitempty"`
}

type Repository struct {
	Name          string `json:"name"`
	FullName      string `json:"full_name"`
	HTMLURL       string `json:"html_url"`
	CloneURL      string `json:"clone_url"`
	DefaultBranch string `json:"default_branch"`
}

type Pusher struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type HeadCommit struct {
	ID        string `json:"id"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	URL       string `json:"url"`
	Author    Author `json:"author"`
}

type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Deployment struct {
	Environment   string `json:"environment"`
	ImageTag      string `json:"image_tag"`
	ContainerName string `json:"container_name"`
}

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("%s🚀 Webhook Test Tool for Vietnam Admin API%s\n", colorBlue, colorReset)
		fmt.Printf("Usage: %s <webhook_url> <webhook_secret> [environment]\n", os.Args[0])
		fmt.Printf("Example: %s https://webhook1.iceteadev.site/ my_secret_key production\n", os.Args[0])
		os.Exit(1)
	}

	webhookURL := os.Args[1]
	webhookSecret := os.Args[2]
	environment := "staging"

	if len(os.Args) > 3 {
		environment = os.Args[3]
	}

	fmt.Printf("%s🔧 Testing webhook deployment for Vietnam Admin API%s\n", colorCyan, colorReset)
	fmt.Printf("Webhook URL: %s\n", webhookURL)
	fmt.Printf("Environment: %s\n", environment)
	fmt.Printf("Time: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	// Create test payload
	payload := createTestPayload(environment)

	// Convert to JSON
	jsonPayload, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		fmt.Printf("%s❌ Error creating JSON payload: %v%s\n", colorRed, err, colorReset)
		os.Exit(1)
	}

	fmt.Printf("%s📤 Sending webhook payload:%s\n", colorPurple, colorReset)
	fmt.Printf("%s\n", string(jsonPayload))

	// Generate HMAC signature
	signature := generateSignature(jsonPayload, webhookSecret)

	// Send webhook request
	success := sendWebhookRequest(webhookURL, jsonPayload, signature)

	if success {
		fmt.Printf("\n%s✅ Webhook test completed successfully!%s\n", colorGreen, colorReset)
		fmt.Printf("%s💡 Check your deployment server logs for deployment progress%s\n", colorYellow, colorReset)
	} else {
		fmt.Printf("\n%s❌ Webhook test failed!%s\n", colorRed, colorReset)
		fmt.Printf("%s💡 Check webhook server configuration and try again%s\n", colorYellow, colorReset)
		os.Exit(1)
	}
}

func createTestPayload(environment string) WebhookPayload {
	now := time.Now()
	commitID := fmt.Sprintf("test-%d", now.Unix())

	var ref string
	var imageTag string

	if environment == "production" {
		ref = "refs/heads/main"
		imageTag = "latest"
	} else {
		ref = "refs/heads/develop"
		imageTag = fmt.Sprintf("develop-%s", commitID)
	}

	return WebhookPayload{
		Ref: ref,
		Repository: Repository{
			Name:          "vietnam-admin-api",
			FullName:      "owner/vietnam-admin-api",
			HTMLURL:       "https://github.com/owner/vietnam-admin-api",
			CloneURL:      "https://github.com/owner/vietnam-admin-api.git",
			DefaultBranch: "main",
		},
		Pusher: Pusher{
			Name:  "test-user",
			Email: "test@example.com",
		},
		HeadCommit: HeadCommit{
			ID:        commitID,
			Message:   fmt.Sprintf("Test deployment to %s", environment),
			Timestamp: now.Format(time.RFC3339),
			URL:       fmt.Sprintf("https://github.com/owner/vietnam-admin-api/commit/%s", commitID),
			Author: Author{
				Name:  "Test User",
				Email: "test@example.com",
			},
		},
		Deployment: Deployment{
			Environment:   environment,
			ImageTag:      imageTag,
			ContainerName: fmt.Sprintf("vietnam-admin-api-%s", environment),
		},
	}
}

func generateSignature(payload []byte, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(payload)
	return "sha256=" + hex.EncodeToString(h.Sum(nil))
}

func sendWebhookRequest(url string, payload []byte, signature string) bool {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Printf("%s❌ Error creating request: %v%s\n", colorRed, err, colorReset)
		return false
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Hub-Signature-256", signature)
	req.Header.Set("X-GitHub-Event", "push")
	req.Header.Set("X-GitHub-Delivery", fmt.Sprintf("test-%d", time.Now().Unix()))
	req.Header.Set("User-Agent", "GitHub-Hookshot/test")

	fmt.Printf("%s📡 Sending request to webhook server...%s\n", colorYellow, colorReset)

	// Print request details
	fmt.Printf("Headers:\n")
	for key, values := range req.Header {
		fmt.Printf("  %s: %s\n", key, values[0])
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("%s❌ Error sending request: %v%s\n", colorRed, err, colorReset)
		return false
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s❌ Error reading response: %v%s\n", colorRed, err, colorReset)
		return false
	}

	// Print response details
	fmt.Printf("\n%s📥 Response received:%s\n", colorPurple, colorReset)
	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	fmt.Printf("Headers:\n")
	for key, values := range resp.Header {
		fmt.Printf("  %s: %s\n", key, values[0])
	}
	fmt.Printf("Body:\n%s\n", string(body))

	// Check if request was successful
	if resp.StatusCode == 200 {
		fmt.Printf("%s✅ Webhook request successful!%s\n", colorGreen, colorReset)
		return true
	} else {
		fmt.Printf("%s❌ Webhook request failed with status: %d%s\n", colorRed, resp.StatusCode, colorReset)
		return false
	}
}
