package envyclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	authHeader      = "Authorization"
	XScopeHeader    = "X-Scope"
	XVarHeader      = "X-Var"
	XVarValueHeader = "X-Var-Value"
)

// Client represents the bridge between the app and an envy instance.
type Client struct {
	httpClient *http.Client
	baseUrl    string
	scope      string
	authKey    string
}

// Var represents variable a name/value pair.
type Var struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// NewClient returns a new client operating on the envy instance identified by the base URL,
// authorized with the auth key (if enabled) and the named scope.
//
// A scope is an arbitrary name used for grouping variables.
// Use a new client for a different scope on the same server.
func NewClient(baseUrl string, scope string, authKey string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		scope:   scope,
		baseUrl: baseUrl,
		authKey: authKey,
	}
}

// SetHTTPClient sets an +http.Client to fit custom needs.
func (c *Client) SetHTTPClient(cl *http.Client) {
	c.httpClient = cl
}

// GetVar fetches a variable by name, sets it as an environment variable (if no error occurred) and returns it.
func (c *Client) GetVar(name string) (Var, error) {

	var v Var

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/getvar", c.baseUrl), nil)
	req.Header.Set(authHeader, fmt.Sprintf("Bearer %s", c.authKey))
	req.Header.Set(XScopeHeader, c.scope)
	req.Header.Set(XVarHeader, name)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return v, err
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return v, err
	}

	err = json.Unmarshal(body, &v)
	if err != nil {
		return v, err
	}

	_ = os.Setenv(v.Name, v.Value)

	return v, nil
}

// GetAllVars fetches all variables, sets them as an environment variables (if no error occurred) and returns them.
func (c *Client) GetAllVars() ([]Var, error) {
	var v []Var

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/getallvars", c.baseUrl), nil)
	req.Header.Set(authHeader, fmt.Sprintf("Bearer %s", c.authKey))
	req.Header.Set(XScopeHeader, c.scope)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return v, err
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return v, err
	}

	err = json.Unmarshal(body, &v)
	if err != nil {
		return nil, err
	}

	for _, e := range v {
		_ = os.Setenv(e.Name, e.Value)
	}

	return v, nil
}

// SetVar permanently sets a variable on the envy instance for later retrieval. Existing variables are overwritten.
func (c *Client) SetVar(name string, value string) error {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/setvar", c.baseUrl), nil)
	req.Header.Set(authHeader, fmt.Sprintf("Bearer %s", c.authKey))
	req.Header.Set(XScopeHeader, c.scope)
	req.Header.Set(XVarHeader, name)
	req.Header.Set(XVarValueHeader, value)

	resp, err := c.httpClient.Do(req)
	_ = resp.Body.Close()
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return fmt.Errorf("expected status code to be 200, got %d", resp.StatusCode)
}
