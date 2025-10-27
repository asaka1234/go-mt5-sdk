package order

import (
	"git.safexinternal.com/tradfi/go-mt5-sdk/utils"
	"github.com/go-resty/resty/v2"
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
