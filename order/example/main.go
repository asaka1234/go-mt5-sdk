package main

import (
	"fmt"
	"github.com/asaka1234/go-mt5-sdk/order"
)

type VLog struct {
}

func (l VLog) Debugf(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}
func (l VLog) Infof(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}
func (l VLog) Warnf(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}
func (l VLog) Errorf(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

func main() {
	vlog := VLog{}

	ADDR := "http://127.0.0.1:8352"

	//构造client
	cli := order.NewClient(vlog, &order.InitParams{ADDR}) //
	cli.SetDebugModel(true)

	//---->开仓-------------
	//resp, err := cli.OpenPosition(GenOpenPositionRequestDemo())
	//---->修改pos-------------
	resp, err := cli.ModifyPosition(GenModifyPositionRequestDemo())
	//---->close-------------
	//resp, err := cli.ClosePosition(GenClosePositionRequestDemo())
	//---->close all-------------
	//resp, err := cli.CloseAllPositions(GenCloseAllPositionRequestDemo())
	//---->挂单-------------
	//resp, err := cli.PlacePendingOrder(GenPlacePendingOrderRequestDemo())
	//---->修改挂单-------------
	//resp, err := cli.ModifyPendingOrder(GenModifyPendingOrderRequestDemo())
	//---->取消挂单-------------
	//resp, err := cli.RemovePendingOrder(GenRemovePendingOrderRequestDemo())
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
		return
	}
	fmt.Printf("resp:%+v\n", resp)
}

func GenOpenPositionRequestDemo() order.OpenPositionRequest {
	return order.OpenPositionRequest{
		Login:   123450079,
		Lots:    "0.5",
		Symbol:  "XAUUSD.s",
		Type:    0,
		Comment: "new",
	}
}

func GenModifyPositionRequestDemo() order.ModifyPositionRequest {
	return order.ModifyPositionRequest{
		Ticket: 467068,
		Tp:     "2.1",
	}
}

func GenClosePositionRequestDemo() order.ClosePositionRequest {
	return order.ClosePositionRequest{
		Lots:    "0.1", //部分平仓
		Ticket:  642331,
		Comment: "uid-1234",
	}
}

func GenCloseAllPositionRequestDemo() order.CloseAllPositionsRequest {
	return order.CloseAllPositionsRequest{
		Login:   123450079,
		Comment: "all close",
	}
}

func GenPlacePendingOrderRequestDemo() order.PlacePendingOrderRequest {
	return order.PlacePendingOrderRequest{
		Login:   123450079,
		Lots:    "0.5",
		Symbol:  "XAUUSD.s",
		Type:    3, // 2-OP_BUY_LIMIT, 3-OP_SELL_LIMIT, 4-OP_BUY_STOP, 5-OP_SELL_STOP，6-OP_BUY_STOP_LIMIT，7-OP_SELL_STOP_LIMIT
		Price:   4237.1,
		Comment: "GenPlacePendingOrderRequestDemo",
	}
}

func GenModifyPendingOrderRequestDemo() order.ModifyPendingOrderRequest {
	return order.ModifyPendingOrderRequest{
		Ticket:  654016,
		Price:   "4225",
		Comment: "ModifyPendingOrderRequest",
	}
}

func GenRemovePendingOrderRequestDemo() order.RemovePendingOrderRequest {
	return order.RemovePendingOrderRequest{
		Ticket:  654016,
		Comment: "RemovePendingOrderRequest",
	}
}
