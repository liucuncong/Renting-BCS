package controllers

import (
	"github.com/astaxie/beego"
	"hkzf/models"
	"path"
	"time"
	"github.com/astaxie/beego/toolbox"
	"fmt"
	"os"
	"encoding/csv"
	"io"
	"strconv"
	"strings"
)

type CerditController struct {
	beego.Controller
}

// 征信认证
func (this *CerditController)Check()  {
	// 1.接收数据
	id := this.GetString("id")
	rank := this.GetString("rank")
	// 2.检验数据
	if rank == "" || id == ""{
		handleResponse(this.Ctx.ResponseWriter,400,"Request parameter rank(or id) can't be empty")
		return
	}

	// 3.区块链查询数据
	//初始化ccs
	channelId := beego.AppConfig.String("channel_id_zxzx")
	chaincodeId := beego.AppConfig.String("chaincode_id_credit")
	userId := beego.AppConfig.String("user_id")
	conf := beego.AppConfig.String("CORE_OZXZX_CONFIG_FILE")

	ccs, err := models.Initialize(channelId, chaincodeId, userId, conf)
	if err != nil {
		handleResponse(this.Ctx.ResponseWriter,500,err.Error())
		return
	}
	defer ccs.Close()
	//调用查询函数
	args := [][]byte{[]byte(id),[]byte(rank)}
	response, err := ccs.ChainCodeQuery("check", args)
	if err != nil {
		handleResponse(this.Ctx.ResponseWriter,500,err.Error())
		return
	}

	// 4.返回结果
	handleResponse(this.Ctx.ResponseWriter,200,response)
}

// 征信认证记录上传
func (this *CerditController)RecordCredit()  {
	// 1.接收文件
	var key = "credit"
	_, header, err := this.GetFile(key)
	if err != nil {
		handleResponse(this.Ctx.ResponseWriter,400,err.Error())
		return
	}
	fileName := header.Filename
	beego.Info("fileName:",fileName)

	// 2.存储文件
	err = this.SaveToFile(key,path.Join("static/upload",fileName))
	if err != nil {
		handleResponse(this.Ctx.ResponseWriter,500,err.Error())
		return
	}
	// 3.5秒后开启任务
	var myTask = "tk2"
	t := time.Now().Add(5*time.Second)
	second := t.Second()
	minute := t.Minute()
	hour := t.Hour()
	spec := fmt.Sprintf("%d %d %d * * *",second,minute,hour)

	// 新建任务
	task := toolbox.NewTask(myTask,spec, func() error {
		// 停止任务
		defer toolbox.StopTask()
		return myTaskCredit(fileName)
	})
	// 添加任务
	toolbox.AddTask(myTask,task)
	// 开启任务
	toolbox.StartTask()

}

func myTaskCredit(fileName string) error {
	// 1.读取文件
	// 2.写入区块链
	// 3.异常处理


	//初始化ccs
	channelId := beego.AppConfig.String("channel_id_zxzx")
	chaincodeId := beego.AppConfig.String("chaincode_id_credit")
	userId := beego.AppConfig.String("user_id")
	conf := beego.AppConfig.String("CORE_OZXZX_CONFIG_FILE")

	ccs, err := models.Initialize(channelId, chaincodeId, userId, conf)
	if err != nil {
		beego.Error(err.Error())
		return err
	}
	defer ccs.Close()
	file, _ := os.Open(path.Join("static/upload", fileName))
	reader := csv.NewReader(file)

	line := 0
	lines := []string{}

	for {
		line++
		lineStr := strconv.Itoa(line)
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			lines = append(lines, lineStr)
			continue
		}
		if len(record) != 2{
			lines = append(lines, lineStr)
			continue
		}

		args := [][]byte{}
		for _, value := range record {
			args = append(args, []byte(value))

		}
		_,err = ccs.ChainCodeUpdate("add",args)
		if err != nil {
			lines = append(lines, lineStr)
		}
	}

	if len(lines) > 0{
		beego.Info("Error lines:",strings.Join(lines,","))
	} else {
		beego.Info("write data success")
	}

	return nil
}