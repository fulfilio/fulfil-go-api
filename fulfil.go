package fulfil

import "github.com/fulfilio/fulfil-go-api/client"

// NewClient creates a new fulfil client.
func NewClient(subDomain string, apiKey string) *client.Client {
    return client.NewClient(subDomain, apiKey)
}
