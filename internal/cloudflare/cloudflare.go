// Cloudflare API integration
package cloudflare

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type CloudflareManager struct {
	ApiToken string
	ZoneID   string
	DB       *sql.DB
}

func NewCloudflareManager(apiToken, zoneID string, db *sql.DB) *CloudflareManager {
	return &CloudflareManager{
		ApiToken: apiToken,
		ZoneID:   zoneID,
		DB:       db,
	}
}

type LockdownRule struct {
	ID             string           `json:"id,omitempty"` // Include for updates, omit for creates
	Description    string           `json:"description"`
	URLs           []string         `json:"urls"`
	Configurations []*Configuration `json:"configurations"`
	Paused         bool             `json:"paused"`
}

type Configuration struct {
	Target string `json:"target"` // e.g., "ip" for IP addresses
	Value  string `json:"value"`  // The actual IP address or range
}

type ApiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ApiResponse struct {
	Success  bool         `json:"success"`
	Errors   []ApiError   `json:"errors"`
	Messages []string     `json:"messages"`
	Result   LockdownRule `json:"result"`
}

func (manager *CloudflareManager) FetchLockdownRule(ruleID string) (LockdownRule, error) {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/firewall/lockdowns/%s", manager.ZoneID, ruleID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+manager.ApiToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return LockdownRule{}, err
	}
	defer resp.Body.Close()

	var response ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return LockdownRule{}, err
	}

	if !response.Success {
		for _, apiError := range response.Errors {
			log.Printf("API Error: %+v\n", apiError)
		}

		return LockdownRule{}, fmt.Errorf("failed to update lockdown rule: %v", response.Errors)
	}

	return response.Result, nil
}

func (manager *CloudflareManager) UpdateLockdownRule(rule LockdownRule) error {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/firewall/lockdowns/%s", manager.ZoneID, rule.ID)

	payload, _ := json.Marshal(rule)
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer "+manager.ApiToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var response ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	log.Println(response)

	if !response.Success {
		for _, apiError := range response.Errors {
			log.Printf("API Error: %+v\n", apiError)
		}

		return fmt.Errorf("failed to update lockdown rule: %v", response.Errors)
	}

	return nil
}
