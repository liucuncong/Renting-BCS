package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

//处理数据回复
func handleResponse(response *context.Response,code int,msg interface{})  {
	if code >= 400{
		beego.Error(msg)
	}else {
		beego.Info(msg)
	}
	response.WriteHeader(code)

	b,ok := msg.([]byte)
	if ok {
		response.Write(b)
	}else {
		s := msg.(string)
		response.Write([]byte(s))
	}
}







