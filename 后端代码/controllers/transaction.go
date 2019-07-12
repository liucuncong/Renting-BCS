package controllers

import (
	"github.com/astaxie/beego"
	"hkzf/models"
)

type TransactionController struct {
	beego.Controller
}

func (this *TransactionController) SetValue () {
	// key  订单编号:期数
	// value  from:to:金额:是否逾期:类型:备注
	orderId := this.GetString("orderId")
	issue := this.GetString("issue")
	if orderId == "" || issue == ""{
		handleResponse(this.Ctx.ResponseWriter,400,"Request paramter orderId(issue) can't be empty")
		return
	}

	from := this.GetString("from")
	to := this.GetString("to")
	rent := this.GetString("rent")
	overdue := this.GetString("overdue")
	types := this.GetString("types")
	desc := this.GetString("desc")

	//初始化ccs
	channelId := beego.AppConfig.String("channel_id_union")
	chaincodeId := beego.AppConfig.String("chaincode_id_transaction")
	userId := beego.AppConfig.String("user_id")
	conf := beego.AppConfig.String("CORE_TRAN_CONFIG_FILE")

	ccs, err := models.Initialize(channelId, chaincodeId, userId, conf)
	if err != nil {
		handleResponse(this.Ctx.ResponseWriter,500,err.Error())
		return
	}
	defer ccs.Close()

	var args [][]byte
	args = append(args, []byte(orderId))
	args = append(args, []byte(issue))
	args = append(args, []byte(from))
	args = append(args, []byte(to))
	args = append(args, []byte(rent))
	args = append(args, []byte(overdue))
	args = append(args, []byte(types))
	args = append(args, []byte(desc))

	response, err := ccs.ChainCodeUpdate("set", args)
	if err != nil {
		handleResponse(this.Ctx.ResponseWriter,500,err.Error())
		return
	}

	handleResponse(this.Ctx.ResponseWriter,200,response)

}


func (this *TransactionController) GetValue ()  {
	orderId := this.GetString("orderId")
	issue := this.GetString("issue")
	if orderId == "" || issue == ""{
		handleResponse(this.Ctx.ResponseWriter,400,"Request paramter orderId(issue) can't be empty")
		return
	}

	//初始化ccs
	channelId := beego.AppConfig.String("channel_id_union")
	chaincodeId := beego.AppConfig.String("chaincode_id_transaction")
	userId := beego.AppConfig.String("user_id")
	conf := beego.AppConfig.String("CORE_TRAN_CONFIG_FILE")

	ccs, err := models.Initialize(channelId, chaincodeId, userId, conf)
	if err != nil {
		handleResponse(this.Ctx.ResponseWriter,500,err.Error())
		return
	}
	defer ccs.Close()

	var args [][]byte
	args = append(args, []byte(orderId))
	args = append(args, []byte(issue))

	response, err := ccs.ChainCodeQuery("get", args)
	if err != nil {
		handleResponse(this.Ctx.ResponseWriter,500,err.Error())
		return
	}

	handleResponse(this.Ctx.ResponseWriter,200,response)
}
