package judge0

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
    baseURL    string
    apiKey     string
    httpClient *http.Client
}

type SubmissionRequest struct {
    SourceCode    string   `json:"source_code"`
    LanguageID    string   `json:"language_id"`
    Stdin         string   `json:"stdin"`
    ExpectedOutput string  `json:"expected_output"`
    TimeLimit     float64  `json:"time_limit"`
    MemoryLimit   int      `json:"memory_limit"`
}

type SubmissionResponse struct {
    Token string `json:"token"`
}

type SubmissionStatus struct {
    Status    Status `json:"status"`
    Stdout    string `json:"stdout"`
    Time      string `json:"time"`
    Memory    int    `json:"memory"`
    stderr    string `json:"stderr"`
    Message   string `json:"message"`
    ExitCode  int    `json:"exit_code"`
}

type Status struct {
    ID          int    `json:"id"`
    Description string `json:"description"`
}

func NewClient(baseURL, apiKey string) *Client {
    return &Client{
        baseURL: baseURL,
        apiKey:  apiKey,
        httpClient: &http.Client{
            Timeout: time.Second * 10,
        },
    }
}

func (c *Client) SubmitCode(req SubmissionRequest) (string, error) {
    url := fmt.Sprintf("%s/submissions?base64_encoded=false", c.baseURL)
    
    body, err := json.Marshal(req)
    if err != nil {
        return "", fmt.Errorf("error marshaling request: %v", err)
    }

    request, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
    if err != nil {
        return "", fmt.Errorf("error creating request: %v", err)
    }


    response, err := c.httpClient.Do(request)
    if err != nil {
        return "", fmt.Errorf("error making request: %v", err)
    }
    defer response.Body.Close()

    if response.StatusCode != http.StatusCreated {
        return "", fmt.Errorf("unexpected status code: %d", response.StatusCode)
    }

    var result SubmissionResponse
    if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
        return "", fmt.Errorf("error decoding response: %v", err)
    }

    return result.Token, nil
}

func (c *Client) GetSubmissionStatus(token string) (*SubmissionStatus, error) {
    url := fmt.Sprintf("%s/submissions/%s?base64_encoded=false", c.baseURL, token)

    request, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("error creating request: %v", err)
    }

    response, err := c.httpClient.Do(request)
    if err != nil {
        return nil, fmt.Errorf("error making request: %v", err)
    }
    defer response.Body.Close()

    if response.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
    }

    var status SubmissionStatus
    if err := json.NewDecoder(response.Body).Decode(&status); err != nil {
        return nil, fmt.Errorf("error decoding response: %v", err)
    }

    return &status, nil
} 