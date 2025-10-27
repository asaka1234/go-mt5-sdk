package order

import (
	"crypto/tls"
	"fmt"
	"github.com/json-iterator/go"
	"gitlab.safexinternal.com/tradefi/go-mt5-sdk/utils"
)

// 挂单
func (cli *Client) RemovePendingOrder(req RemovePendingOrderRequest) (*CommonResp, error) {

	rawURL := cli.Params.RemovePendingUrl

	//返回值会放到这里
	var result CommonResp

	resp, err := cli.ryClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		SetCloseConnection(true).
		R().
		SetBody(req).
		SetHeaders(getHeaders()).
		SetDebug(cli.debugMode).
		SetResult(&result).
		SetError(&result).
		Post(rawURL)

	//print log
	restLog, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(utils.GetRestyLog(resp))
	cli.logger.Infof("MT5#RemovePendingOrder->%+v", string(restLog))

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
