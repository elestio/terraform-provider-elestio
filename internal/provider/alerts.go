package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/elestio/elestio-go-api-client/v2"
)

// alertRule is one monitoring threshold, as the Elestio backend expects it.
type alertRule struct {
	Parameter string `json:"parameter"`
	Value     int    `json:"value"`
	Cycles    int    `json:"cycles"`
	Unit      string `json:"unit"`
}

// defaultAlertRules mirrors the backend's own defaultMonicAlertRules
// (helper/constant.js): the rules a new service gets. The backend builds its
// monit configuration from every one of them and fails with a generic 500
// ActionError when one is missing. elestio-go-api-client v2.2.0 sends only ten of
// the twelve (it lacks READ_RATE_OPERATIONS and WRITE_RATE_OPERATIONS), so its
// EnableAlerts can never succeed against the current backend and this provider
// sends the request itself.
var defaultAlertRules = []alertRule{
	{"CPU", 90, 15, "%"},
	{"MEMORY", 90, 15, "%"},
	{"SWAP", 90, 15, "%"},
	{"SPACE", 80, 15, "%"},
	{"INODE", 80, 15, "%"},
	{"READ_RATE", 20, 15, "MB/s"},
	{"READ_RATE_OPERATIONS", 500, 5, "operations/s"},
	{"WRITE_RATE", 20, 15, "MB/s"},
	{"WRITE_RATE_OPERATIONS", 500, 5, "operations/s"},
	{"SATURATION", 90, 15, "%"},
	{"DOWNLOAD", 25, 15, "MB/s"},
	{"UPLOAD", 25, 15, "MB/s"},
}

const defaultMonitCycleInSeconds = 60

// enableAlerts turns on monitoring alerts for a service with the default rules.
// The request goes through the API client's HTTP client, whose transport adds
// the session JWT (the empty "jwt" field below is replaced).
func enableAlerts(ctx context.Context, client *elestio.Client, serviceID string) error {
	rules, err := json.Marshal(defaultAlertRules)
	if err != nil {
		return err
	}
	body, err := json.Marshal(map[string]any{
		"jwt":                 "",
		"vmID":                serviceID,
		"action":              "enableAlerts",
		"monitCycleInSeconds": defaultMonitCycleInSeconds,
		"rules":               string(rules), // the backend expects the rules as a JSON string
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, client.BaseURL+"/api/servers/DoActionOnServer", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, respBody)
	}
	var out struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	if json.Unmarshal(respBody, &out) == nil && out.Status == "KO" {
		return fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, out.Message)
	}
	return nil
}
