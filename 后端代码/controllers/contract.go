package controllers

import (
	"github.com/astaxie/beego"
	"strings"
	"crypto/sha256"
	"io"
	"encoding/hex"
	"hkzf/models"
)

type ContractController struct {
	beego.Controller
}

func (this *ContractController) SetValue ()  {
	// 接收合同图片
	// 对图片进行sha256 --value
	// 约定图片的名称：key

	file, header, err := this.GetFile("contract")
	if err != nil {
		handleResponse(this.Ctx.ResponseWriter,400,err.Error())
		return
	}
	defer file.Close()

	fileName := header.Filename
	split := strings.Split(fileName, ".")
	key := split[0]

	beego.Info("key",key)
	// 获取图片的sha256
	hash := sha256.New()
	//hash.Write([]byte)
	// 文件操作
	_,err = io.Copy(hash,file)
	if err != nil {
		handleResponse(this.Ctx.ResponseWriter,400,err.Error())
		return
	}
	sum := hash.Sum(nil)

	value := hex.EncodeToString(sum)

	//初始化ccs
	channelId := beego.AppConfig.String("channel_id_union")
	chaincodeId := beego.AppConfig.String("chaincode_id_contract")
	userId := beego.AppConfig.String("user_id")
	conf := beego.AppConfig.String("CORE_OUNION_CONFIG_FILE")

	ccs, err := models.Initialize(channelId, chaincodeId, userId, conf)
	if err != nil {
		handleResponse(this.Ctx.ResponseWriter,500,err.Error())
		return
	}
	defer ccs.Close()

	var args [][]byte
	args = append(args, []byte(key))
	args = append(args, []byte(value))

	response, err := ccs.ChainCodeUpdate("set", args)
	if err != nil {
		handleResponse(this.Ctx.ResponseWriter,500,err.Error())
		return
	}

	handleResponse(this.Ctx.ResponseWriter,200,response)
}


func (this *ContractController) GetValue ()  {
	// 接收合同图片
	// 对图片进行sha256 --value
	// 约定图片的名称：key

	key := this.GetString("contractId")
	if key == "" {
		handleResponse(this.Ctx.ResponseWriter,400,"Request paramter contractId can't be empty")
		return
	}
	beego.Info("key",key)

	//初始化ccs
	channelId := beego.AppConfig.String("channel_id_union")
	chaincodeId := beego.AppConfig.String("chaincode_id_transaction")
	userId := beego.AppConfig.String("user_id")
	conf := beego.AppConfig.String("CORE_OUNION_CONFIG_FILE")

	ccs, err := models.Initialize(channelId, chaincodeId, userId, conf)
	if err != nil {
		handleResponse(this.Ctx.ResponseWriter,500,err.Error())
		return
	}
	defer ccs.Close()

	var args [][]byte
	args = append(args, []byte(key))

	response, err := ccs.ChainCodeQuery("get", args)
	if err != nil {
		handleResponse(this.Ctx.ResponseWriter,500,err.Error())
		return
	}

	handleResponse(this.Ctx.ResponseWriter,200,response)
}


