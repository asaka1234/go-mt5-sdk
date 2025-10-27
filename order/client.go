package order

import (
	"github.com/go-resty/resty/v2"
	"safexinternal.com/tradfi/go-mt5-sdk/utils"
)

type Client struct {
	Params *InitParams

	ryClient  *resty.Client
	debugMode bool
	logger    utils.Logger
}

func NewClient(logger utils.Logger, params *InitParams) *Client {
	return &Client{
		Params: params,

		ryClient:  resty.New(), //client实例
		debugMode: false,
		logger:    logger,
	}
}

func (cli *Client) SetDebugModel(debugModel bool) {
	cli.debugMode = debugModel
}
