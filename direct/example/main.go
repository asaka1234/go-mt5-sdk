package main

import (
	"fmt"
	"github.com/asaka1234/go-mt5-sdk/direct"
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

	//_, err := cli.ListSymbol()
	//resp, err := cli.TickReview()

	//req := GenUserCreateReqDemo()
	//_, err := cli.UserCreate(req)

	req := GenBalanceOperationReqDemo()
	_, err := cli.BalanceOperation(req)

	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
		return
	}
	//fmt.Printf("resp:%+v\n", resp)

}

func GenUserCreateReqDemo() direct.UserCreateReq {
	return direct.UserCreateReq{
		Uid:      123,
		Internal: 2,
	}
}

func GenBalanceOperationReqDemo() direct.BalanceOperationReq {
	return direct.BalanceOperationReq{
		Login:   123450051,
		Balance: -30,
	}
}
