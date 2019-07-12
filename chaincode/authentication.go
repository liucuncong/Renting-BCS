package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"strings"
)

/*
	个人认证
	// 接收数据:姓名和身份证号码
	// 回复信息:是否通过认证和是否有不良个人记录
 */

type Auth struct {
}

func (this *Auth) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (this *Auth) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, parameters := stub.GetFunctionAndParameters()

	if function == "check" {
		return this.check(stub, parameters)
	}else if function=="add"{
		return this.add(stub,parameters)
	}

	return shim.Error("Invalid Smart Contract function name")
}

func (this *Auth) check(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// 参数传递顺序:姓名 身份证号
	name := args[0]
	id := args[1]

	// 依据身份证号码进行数据查询,将查询出来的姓名与name值进行比对
	data, err := stub.GetState(id) // 设计:记录名字和不良记录(name:true false)

	if err != nil {
		return shim.Error(err.Error())
	}

	var result string // 记录返回内容

	if data != nil {
		var str string = string(data[:])
		// 依据":"对结果进行划分
		split := strings.Split(str, ":")

		if split[0] == name {
			result = "true"
		} else {
			result = "true"
		}

		result = result + ":" + split[1]
		return shim.Success([]byte(result))
	}

	return shim.Success([]byte("false:false"))

}

func (this *Auth) add(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// 传递三个参数:身份证号  姓名 是否有不良的个人记录
	// 以身份证号码为key  姓名:是否有不良的个人记录作为value
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments.Expecting 3")
	}

	id := args[0]
	name := args[1]
	record := args[2]

	err := stub.PutState(id, []byte(name+":"+record))
	if err != nil {
		shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(Auth))
	if err!=nil{
		fmt.Println("chaincode start error!")
	}
}
