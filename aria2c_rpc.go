package ariago

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/czasg/snow"
	"io/ioutil"
	"net/http"
)

var _ Aria = (*Aria2c)(nil)

type Aria2c struct {
	RPCUrl string `default:"http://localhost:6800/jsonrpc"`
}

func (a *Aria2c) AddURI(ctx context.Context, uri string, params map[string]string) (*AddURIResponse, error) {
	if params == nil {
		params = map[string]string{}
	}
	if params["out"] == "" {
		params["out"] = fmt.Sprintf("%d", snow.Next())
	}
	body, err := json.Marshal(map[string]interface{}{
		"id":      params["out"],
		"jsonrpc": "2.0",
		"method":  "aria2.addUri",
		"params": []interface{}{
			[]string{uri},
			params,
		},
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.RPCUrl, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	response := AddURIResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error.Code != 0 {
		err = errors.New(response.Error.Message)
		return nil, err
	}
	return &response, nil
}

func (a *Aria2c) CheckStatus(ctx context.Context, gid string) (*CheckStatusResponse, error) {
	body, err := json.Marshal(map[string]interface{}{
		"id":      "",
		"jsonrpc": "2.0",
		"method":  "aria2.tellStatus",
		"params":  []string{gid},
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.RPCUrl, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	response := CheckStatusResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}