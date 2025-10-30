package direct

import (
	"crypto/tls"
	"fmt"
	"github.com/asaka1234/go-mt5-sdk/utils"
	"github.com/json-iterator/go"
)

// 全量获取支持的所有symbol

func (cli *Client) TickReview() (*TickReviewResp, error) {

	rawURL := cli.Params.Address + "/v1/tick/review"

	//返回值会放到这里
	var result TickReviewResp

	resp, err := cli.ryClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		SetCloseConnection(true).
		R().
		SetHeaders(getHeaders()).
		SetDebug(cli.debugMode).
		SetResult(&result).
		SetError(&result).
		Get(rawURL)

	//print log
	restLog, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(utils.GetRestyLog(resp))
	cli.logger.Infof("MT5#TickReview#List->%+v", string(restLog))

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		//反序列化错误会在此捕捉
		return nil, fmt.Errorf("status code: %d", resp.StatusCode())
	}

	if resp.Error() != nil {
		//反序列化错误会在此捕捉
		return nil, fmt.Errorf("%v, body:%s", resp.Error(), resp.Body())
	}

	return &result, err
}
