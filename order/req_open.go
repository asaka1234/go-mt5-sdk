package order

import (
	"crypto/tls"
	"errors"
	"github.com/json-iterator/go"
	"safexapp.com/tradfi/go-mt5-sdk/utils"
)

// 开仓
func (cli *Client) OpenOrder(req OpenRequest) (*OpenResponse, error) {

	rawURL := cli.Params.OpenUrl

	//返回值会放到这里
	var result OpenResponse

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
	cli.logger.Infof("MT5#PlaceOrder#Open->%+v", string(restLog))

	if err != nil {
		return nil, err
	}

	//-----------错误处理------------------------
	if resp.StatusCode() != 201 {
		if result.Error != "" {
			return nil, errors.New(result.Error)
		}
		return nil, errors.New(result.RespMsg)
	}

	return &result, err
}
