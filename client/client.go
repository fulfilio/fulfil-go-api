package client

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "io/ioutil"
    "net/http"
    "log"
)

type Model struct {
    name string
}

type Report struct {
    name string
}

type Client struct {
    subDomain string
    apiKey string
    context []string
    logger *log.Logger
}

// NewClient creates a new go API client instance.
func NewClient(subDomain string, apiKey string) *Client {
    return &Client{
        subDomain: subDomain,
        apiKey: apiKey,
    }
}

// SetTraceLogger sets a new trace logger and returns the old logger.
// If logger is not nil, Client will output all internal messages to the logger.
func (c *Client) SetTraceLogger(logger *log.Logger) *log.Logger {
    oldLogger := c.logger
    c.logger = logger
    return oldLogger
}

func (c *Client) Execute(method, path string, input, output interface{}) error {
    url := fmt.Sprintf("https://%s.fulfil.io/%s", c.subDomain, path)

    if c.logger != nil {
        c.logger.Printf("Request Endpoint: %s", url)
    }

    req, err := c.createRequest(method, url, input)
    if err != nil {
        return fmt.Errorf("Error creating request object: %s", err.Error())
    }
    c.executeRequest(req, output)
    return nil
}

func (c *Client) RefreshContext(method, path string, input, output interface{}) error {
    url := fmt.Sprintf("https://%s.fulfil.io/%s", c.subDomain, path)

    if c.logger != nil {
        c.logger.Printf("Request Endpoint: %s", url)
    }

    req, err := c.createRequest(method, url, input)
    if err != nil {
        return fmt.Errorf("Error creating request object: %s", err.Error())
    }
    c.executeRequest(req, output)
    return nil
}

func (c *Client) createRequest(method, url string, bodyObject interface{}) (req *http.Request, err error) {
    var reqBody io.Reader
    if bodyObject != nil {
        data, err := json.Marshal(bodyObject)
        if err != nil {
            return nil, fmt.Errorf("Error marshaling body object: %s", err.Error())
        }
        reqBody = bytes.NewBuffer(data)
    }
    req, err = http.NewRequest(method, url, reqBody)
    if err != nil {
	return nil, fmt.Errorf("Error creating HTTP request: %s", err.Error())
    }

    q := req.URL.Query()
    q.Add("context", "{}")

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")
    req.Header.Set("User-Agent", "Golang Client")
    req.Header.Set("x-api-key", c.apiKey)

    req.Header.Set("Connection", "close")
    req.Close = true

    return req, nil
}

func (c *Client) executeRequest(req *http.Request, output interface{}) (err error) {
    httpClient := http.Client{}
    res, err := httpClient.Do(req)
    if err != nil {
        return fmt.Errorf("Error making HTTP request: %s", err.Error())
    }
    defer res.Body.Close()

    resData, err := ioutil.ReadAll(res.Body)
    if err != nil {
	return fmt.Errorf("Error reading response body data: %s", err.Error())
    }

    if c.logger != nil {
	c.logger.Printf("Client.executeRequest() response: status=%q, body=%q", res.Status, string(resData))
    }
    json.Unmarshal(resData, output)
    return nil
}
