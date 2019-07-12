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

type CertificationController struct {
	beego.Controller
}

// 认证房屋
func (this *CertificationController)Check()  {

	houseId := this.GetString("houseId")
	id := this.GetString("id")
	if houseId == "" || id == ""{
		handleResponse(this.Ctx.ResponseWriter,400,"Request parameter houseId(or id) can't be empty")
		return
	}
	beego.Info(houseId+":"+id)

	var (
		channelId = beego.AppConfig.String("channel_id_fgj")
		chaincodeId = beego.AppConfig.String("chaincode_id_house")
		userID = beego.AppConfig.String("user_id")
		conf = beego.AppConfig.String("CORE_OFGJ_CONFIG_FILE")
	)

	ccs, err := models.Initialize(channelId, chaincodeId, userID, conf)
	if err != nil {
		handleResponse(this.Ctx.ResponseWriter,500,err.Error())
		return
	}
	defer ccs.Close()

	args := [][]byte{[]byte(houseId),[]byte(id)}
	response,err := ccs.ChainCodeQuery("check",args)
	if err != nil {
		handleResponse(this.Ctx.ResponseWriter,500,err.Error())
		return
	}
	handleResponse(this.Ctx.ResponseWriter,200,response)

}


// 上传房屋记录信息
func (this *CertificationController)RecordHouse() {
	// 上传房屋记录信息
	var key = "house"
	_, header, err := this.GetFile(key)
	if err != nil {
		handleResponse(this.Ctx.ResponseWriter,400,err.Error())
		return
	}

	fileName := header.Filename
	beego.Info("文件名称：",fileName)

	err = this.SaveToFile(key, path.Join("static/upload", fileName))
	if err != nil {
		handleResponse(this.Ctx.ResponseWriter,500,err.Error())
		return
	}

	//开启任务
	var myTask = "tk1"

	t := time.Now().Add(5*time.Second)
	second := t.Second()
	minute := t.Minute()
	hour := t.Hour()
	spec := fmt.Sprintf("%d %d %d * * *",second,minute,hour)

	task := toolbox.NewTask(myTask,spec, func() error {
		defer toolbox.StopTask()


		return myTaskHouse(fileName)
	})

	toolbox.AddTask(myTask,task)
	toolbox.StartTask()

	handleResponse(this.Ctx.ResponseWriter,200,"ok")
}

func myTaskHouse(fileName string) error {
	// 读文件
	// 写数据
	// 异常处理

	var (
		channelID = beego.AppConfig.String("channel_id_fgj")
		chaincodeID = beego.AppConfig.String("chaincode_id_house")
		userID = beego.AppConfig.String("user_id")
		conf = beego.AppConfig.String("CORE_OFGJ_CONFIG_FILE")
	)
	ccs, err := models.Initialize(channelID, chaincodeID, userID,conf)
	if err != nil {

		beego.Error(err.Error())
		return err
	}
	defer ccs.Close()


	file, _ := os.Open(path.Join("static/upload", fileName))
	reader := csv.NewReader(file)


	var line = 0
	var lines []string
	for {
		line++
		lineStr := strconv.Itoa(line)

		record,err := reader.Read()
		if err == io.EOF{
			break
		}
		if err != nil {
			lines = append(lines, lineStr)
			continue
		}
		if len(record) != 3{
			lines = append(lines, lineStr)
			continue
		}

		var args [][]byte
		for _, value := range record {
			// 参数
			args = append(args, []byte(value))

		}

		_, err = ccs.ChainCodeUpdate("add", args)
		if err != nil {
			lines = append(lines, lineStr)
		}

	}

	if len(lines) > 0 {
		beego.Error("Error lines:",strings.Join(lines,","))
	} else {
		beego.Info("write data success")
	}

	return nil
}



