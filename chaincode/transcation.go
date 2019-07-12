package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"fmt"
)

/*
	交易记录
	key :订单编号:期数
	value:From:To:金额:是否逾期:类型:备注
 */

type Transcation struct {
}

func (this *Transcation) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (this *Transcation) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, parameters := stub.GetFunctionAndParameters()

	if function == "set" {
		return this.set(stub, parameters)
	} else if function == "get" {
		return this.get(stub, parameters)
	}

	return shim.Error("Invalid Smart Contract function name")
}

func (this *Transcation) get(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	orderId := args[0]
	issue := args[1]

	id := fmt.Sprintf("%s:%s", orderId, issue)
	data, err := stub.GetState(id)

	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(data)
}

func (this *Transcation) set(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 8 {
		return shim.Error("Incorrect number of arguments.Expecting 8")
	}

	orderId := args[0]
	issue := args[1]
	from := args[2]
	to := args[3]
	rent := args[4]
	overdue := args[5]
	types := args[6]
	desc := args[7]

	key := fmt.Sprintf("%s:%s", orderId, issue)
	value := fmt.Sprintf("%s:%s:%s:%s:%s:%s", from, to, rent, overdue, types, desc)

	err := stub.PutState(key, []byte(value))
	if err != nil {
		shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(Transcation))
	if err!=nil{
		fmt.Println("chaincode start error!")
	}
}
