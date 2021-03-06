package ariago

import "context"

type Aria interface {
	AddURI(ctx context.Context, uri string, params map[string]interface{}) (*AddURIResponse, error)
	CheckStatus(ctx context.Context, gid string) (*CheckStatusResponse, error)
}

type AddURIResponse struct {
	GID     string `json:"result"`
	File    string `json:"id"`
	Url     string `json:"-"`
	Version string `json:"jsonrpc"`
	Error   struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type CheckStatusResponse struct {
	JsonRPC string `json:"jsonrpc"`
	Result  struct {
		GID             string `json:"gid"`
		Status          string `json:"status"`
		CompletedLength string `json:"completedLength"`
		TotalLength     string `json:"totalLength"`
		DownloadSpeed   string `json:"downloadSpeed"`
		ErrorCode       string `json:"errorCode"`
		ErrorMessage    string `json:"errorMessage"`
	} `json:"result"`
}

func (c CheckStatusResponse) Done() bool {
	return c.Result.Status == "complete"
}
