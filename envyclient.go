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
	XEnvHeader      = "X-Env"
	XVarHeader      = "X-Var"
	XVarValueHeader = "X-Var-Value"
)

type Client struct {
	httpClient *http.Client
	baseUrl    string
	env        string
	authKey    string
}

type Var struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func NewClient(baseUrl string, env string, authKey string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		env:     env,
		baseUrl: baseUrl,
		authKey: authKey,
	}
}

func (c *Client) SetHTTPClient(cl *http.Client) {
	c.httpClient = cl
}

// GetVar fetches a variable in an environment, sets it as an environment variable (if no error occurred) and returns it.
func (c *Client) GetVar(name string) (Var, error) {

	var v Var

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/getvar", c.baseUrl), nil)
	req.Header.Set(authHeader, fmt.Sprintf("Bearer %s", c.authKey))
	req.Header.Set(XEnvHeader, c.env)
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

// GetAllVars fetches all variables of an environment, sets them as an environment variables (if no error occurred) and returns them.
func (c *Client) GetAllVars() ([]Var, error) {
	var v []Var

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/getallvars", c.baseUrl), nil)
	req.Header.Set(authHeader, fmt.Sprintf("Bearer %s", c.authKey))
	req.Header.Set(XEnvHeader, c.env)

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

func (c *Client) SetVar(name string, value string) error {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/setvar", c.baseUrl), nil)
	req.Header.Set(authHeader, fmt.Sprintf("Bearer %s", c.authKey))
	req.Header.Set(XEnvHeader, c.env)
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
