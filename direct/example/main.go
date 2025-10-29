package main

import (
	"fmt"
	"github.com/asaka1234/go-mt5-sdk/direct"
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

	ADDR := "http://127.0.0.1:8351"

	//构造client
	cli := direct.NewClient(vlog, &direct.InitParams{ADDR})
	cli.SetDebugModel(true)

	//---->symbol list-------------
	resp, err := cli.ListSymbol()
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
		return
	}
	fmt.Printf("resp:%+v\n", resp)

}

func GenOpenPositionRequestDemo() order.OpenPositionRequest {
	return order.OpenPositionRequest{
		Login:  88000047,
		Lots:   0.5,
		Symbol: "XAUUSD",
		Type:   0,
	}
}
