package routers

import (
	"hkzf/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/auth", &controllers.AuthController{},"get:Check;post:RecordAuth")
    beego.Router("/house", &controllers.CertificationController{},"get:Check;post:RecordHouse")
    beego.Router("/credit", &controllers.CerditController{},"get:Check;post:RecordCredit")
    beego.Router("/contract", &controllers.ContractController{},"get:GetValue;post:SetValue")
    beego.Router("/transaction", &controllers.TransactionController{},"get:GetValue;post:SetValue")
}


