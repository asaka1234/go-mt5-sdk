package main

import (
	"fmt"
	"safexapp.com/tradfi/go-mt5-sdk/order"
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

	OPEN_URL := "http://127.0.0.1/v1/order/open"
	CLOSE_URL := "http://127.0.0.1/v1/order/close"

	//构造client
	cli := order.NewClient(vlog, &order.InitParams{OPEN_URL, CLOSE_URL})
	cli.SetDebugModel(true)

	//---->open-------------
	resp, err := cli.OpenPosition(GenOpenPositionRequestDemo())
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
		return
	}
	fmt.Printf("resp:%+v\n", resp)

	//---->close-------------
	resp, err = cli.ClosePosition(GenClosePositionRequestDemo())
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
		return
	}
	fmt.Printf("resp:%+v\n", resp)
}

func GenOpenPositionRequestDemo() order.OpenPositionRequest {
	return order.OpenPositionRequest{
		Login:  700,
		Lots:   0.1,
		Symbol: "XAUUSD",
		Type:   0,
	}
}

func GenClosePositionRequestDemo() order.ClosePositionRequest {
	return order.ClosePositionRequest{
		Lots:   0.1,
		Ticket: 1028028,
	}
}
