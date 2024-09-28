package envyclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
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

func (c *Client) GetVar(name string) (Var, error) {

	var v Var

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/getvar", c.baseUrl), nil)
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
	return v, err
}

func (c *Client) GetAllVars(name string) ([]Var, error) {
	var v []Var

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/getallvars", c.baseUrl), nil)
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
	return v, err
}
