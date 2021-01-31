package client

import (
	"errors"
	"net/url"
)

// Option (s) are call back funcs, passed into the constructor to customise the client
type Option func(c *Client) error

// WithClient set the http client to use
func WithClient(client doer) Option {
	return func(c *Client) error {
		if client == nil {
			return errors.New("nil client passed into WithClient")
		}
		c.client = client
		return nil
	}
}

// WithDomain sets the domain in which request will be sent to
func WithDomain(domain string) Option {
	return func(c *Client) error {
		u, err := url.Parse(domain)
		if err != nil {
			return err
		}
		c.endpoint = u
		return nil
	}
}
