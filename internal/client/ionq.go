package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const ionqBaseURL = "https://api.ionq.co/v0.3"

type IonQClient struct {
	APIKey     string
	HTTPClient *http.Client
}

func NewIonQClient(apiKey string) *IonQClient {
	return &IonQClient{
		APIKey:     apiKey,
		HTTPClient: &http.Client{},
	}
}

type JobInput struct {
	Circuit json.RawMessage `json:"circuit"`
	Shots   int64           `json:"shots"`
	Target  string          `json:"target"`
}

type JobResponse struct {
	ID      string          `json:"id"`
	Status  string          `json:"status"`
	Target  string          `json:"target"`
	Results json.RawMessage `json:"results,omitempty"`
}

type GateCircuit struct {
	Qubits int    `json:"qubits"`
	Gates  []Gate `json:"circuit"`
}

type Gate struct {
	Gate     string  `json:"gate"`
	Target   int     `json:"target,omitempty"`
	Control  int     `json:"control,omitempty"`
	Targets  []int   `json:"targets,omitempty"`
	Rotation float64 `json:"rotation,omitempty"`
}

func (c *IonQClient) CreateJob(ctx context.Context, input JobInput) (*JobResponse, error) {
	body := map[string]interface{}{
		"target": input.Target,
		"shots":  input.Shots,
		"input":  input.Circuit,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshaling job request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, "POST", ionqBaseURL+"/jobs", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "apiKey "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("submitting job: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("IonQ API error (status %d): %s", resp.StatusCode, string(respBody))
	}
	var jobResp JobResponse
	if err := json.NewDecoder(resp.Body).Decode(&jobResp); err != nil {
		return nil, fmt.Errorf("decoding job response: %w", err)
	}
	return &jobResp, nil
}

func (c *IonQClient) GetJob(ctx context.Context, jobID string) (*JobResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", ionqBaseURL+"/jobs/"+jobID, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "apiKey "+c.APIKey)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("getting job: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("IonQ API error (status %d): %s", resp.StatusCode, string(respBody))
	}
	var jobResp JobResponse
	if err := json.NewDecoder(resp.Body).Decode(&jobResp); err != nil {
		return nil, fmt.Errorf("decoding job response: %w", err)
	}
	return &jobResp, nil
}

func (c *IonQClient) CancelJob(ctx context.Context, jobID string) error {
	req, err := http.NewRequestWithContext(ctx, "PUT", ionqBaseURL+"/jobs/"+jobID+"/status/cancel", nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "apiKey "+c.APIKey)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("canceling job: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("IonQ API error (status %d): %s", resp.StatusCode, string(respBody))
	}
	return nil
}

type BackendResponse struct {
	Name        string `json:"backend"`
	Status      string `json:"status"`
	Qubits      int64  `json:"qubits"`
	AverageT1   int64  `json:"average_t1"`
	AverageT2   int64  `json:"average_t2"`
	Description string `json:"characterization_url"`
}

type BackendsListResponse []BackendResponse

func (c *IonQClient) GetBackend(ctx context.Context, name string) (*BackendResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", ionqBaseURL+"/backends/"+name, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "apiKey "+c.APIKey)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("getting backend: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("IonQ API error (status %d): %s", resp.StatusCode, string(respBody))
	}
	var backend BackendResponse
	if err := json.NewDecoder(resp.Body).Decode(&backend); err != nil {
		return nil, fmt.Errorf("decoding backend response: %w", err)
	}
	return &backend, nil
}
