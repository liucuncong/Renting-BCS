package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"

	"strings"
)

/*
	征信认证
*/

type Credit struct {
}

func (this *Credit) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (this *Credit) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, parameters := stub.GetFunctionAndParameters()

	if function == "check" {
		return this.check(stub, parameters)
	}else if function=="add"{
		return this.add(stub,parameters)
	}

	return shim.Error("Invalid Smart Contract function name")
}

//查询身份证号码是否与征信级别匹配
func (this *Credit) check(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// 参数传递顺序:身份证号  征信级别
	id := args[0]
	rank := args[1]

	// 依据身份证号码进行数据查询
	data, err := stub.GetState(id) // 设计:记录名字和不良记录(id:rank)

	if err != nil {
		return shim.Error(err.Error())
	}

	if data != nil {
		str := string(data[:])
		if str == rank{
			//校验成功，返回true
			return shim.Success([]byte("true"))
		}
	}
	//校验失败，返回false
	return shim.Success([]byte("false"))

}

func (this *Credit) add(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// 传递三个参数:身份证号  姓名 是否有不良的个人记录
	// 以身份证号码为key  姓名:是否有不良的个人记录作为value
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments.Expecting 2")
	}

	id := args[0]
	rank := args[1]

	err := stub.PutState(id, []byte(rank))
	if err != nil {
		shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(Credit))
	if err!=nil{
		fmt.Println("chaincode start error!")
	}
}