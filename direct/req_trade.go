package direct

import (
	"crypto/tls"
	"fmt"
	"github.com/asaka1234/go-mt5-sdk/utils"
	"github.com/json-iterator/go"
	"github.com/spf13/cast"
)

// 获取指定login的当前持仓列表

func (cli *Client) ListPosition(login uint64) (*TickReviewResp, error) {

	rawURL := cli.Params.Address + "/v1/position/list"

	//返回值会放到这里
	var result TickReviewResp

	resp, err := cli.ryClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		SetCloseConnection(true).
		R().
		SetHeaders(getHeaders()).
		SetQueryParam("login", cast.ToString(login)).
		SetDebug(cli.debugMode).
		SetResult(&result).
		SetError(&result).
		Get(rawURL)

	//print log
	restLog, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(utils.GetRestyLog(resp))
	cli.logger.Infof("MT5#ListPosition->%+v", string(restLog))

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

// 获取当前的挂单列表

func (cli *Client) ListPendingOrder(login uint64) (*ListPendingOrderResp, error) {

	rawURL := cli.Params.Address + "/v1/pendingOrder/list"

	//返回值会放到这里
	var result ListPendingOrderResp

	resp, err := cli.ryClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		SetCloseConnection(true).
		R().
		SetHeaders(getHeaders()).
		SetQueryParam("login", cast.ToString(login)).
		SetDebug(cli.debugMode).
		SetResult(&result).
		SetError(&result).
		Get(rawURL)

	//print log
	restLog, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(utils.GetRestyLog(resp))
	cli.logger.Infof("MT5#ListPendingOrder->%+v", string(restLog))

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
