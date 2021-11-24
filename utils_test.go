package ariago

import (
	"context"
	"errors"
	"strconv"
	"time"
)

func WatchWrap(aria Aria, tolerateStandard int) func(ctx context.Context, timeout int64, uri string, params map[string]interface{}) (*AddURIResponse, *CheckStatusResponse, error) {
	return func(ctx context.Context, timeout int64, uri string, params map[string]interface{}) (*AddURIResponse, *CheckStatusResponse, error) {
		addResp, err := aria.AddURI(ctx, uri, params)
		if err != nil {
			return nil, nil, err
		}
		tolerate := 0
		for {
			time.Sleep(time.Second)
			timeout--
			checkResp, err := aria.CheckStatus(ctx, addResp.GID)
			if err != nil {
				return nil, nil, err
			}
			if checkResp.Done() {
				return addResp, checkResp, nil
			}
			total, _ := strconv.ParseInt(checkResp.Result.TotalLength, 0, 64)
			completed, _ := strconv.ParseInt(checkResp.Result.CompletedLength, 0, 64)
			speed, _ := strconv.ParseInt(checkResp.Result.DownloadSpeed, 0, 64)
			if speed < 1 || (total != 0 && speed*timeout < (total-completed)) {
				tolerate++
			}
			if tolerate >= tolerateStandard {
				return nil, nil, errors.New("timeout")
			}
		}
	}
}
