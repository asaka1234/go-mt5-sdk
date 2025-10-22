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
	//发请求
	resp, err := cli.OpenOrder(GenOpenOrderRequestDemo())
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
		return
	}
	fmt.Printf("resp:%+v\n", resp)
}

func GenOpenOrderRequestDemo() order.OpenRequest {
	return order.OpenRequest{
		Amount: 700,
		Ref1:   "9438239651111", //Unique customer id in your system. 业务系统里的唯一客户id
		Ref3:   "uname",
	}
}
