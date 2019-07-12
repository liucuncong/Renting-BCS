package controllers

import (
	"github.com/astaxie/beego"
	"hkzf/models"
	"path"
	"github.com/astaxie/beego/toolbox"
	"time"
	"fmt"
	"os"
	"encoding/csv"
	"io"
	"strconv"
	"strings"
)

type AuthController struct {
	beego.Controller
}

func (this *AuthController) Check ()  {
	// 发送数据，用户姓名和用户身份证号码
	name := this.GetString("name")
	id := this.GetString("id")

	if name == "" || id == ""{
		handleResponse(this.Ctx.ResponseWriter,400,"Request parameter name(or id) can't be empty")
		return
	}
	beego.Info(name+":"+id)

	// 利用models中的ChainCodeQuery查询当前用户的匹配结果和是否有个人不良记录的结果
	var (
		channelID = beego.AppConfig.String("channel_id_gaj")
		chaincodeID = beego.AppConfig.String("chaincode_id_auth")
		userID = beego.AppConfig.String("user_id")
		conf = beego.AppConfig.String("CORE_OGAJ_CONFIG_FILE")
	)
	ccs, err := models.Initialize(channelID, chaincodeID, userID,conf)
	if err != nil {
		handleResponse(this.Ctx.ResponseWriter,500,err.Error())
		return
	}
	defer ccs.Close()

	args := [][]byte{[]byte(name),[]byte(id)}
	response, err := ccs.ChainCodeQuery("check", args)
	if err != nil {
		handleResponse(this.Ctx.ResponseWriter,500,err.Error())
		return
	}

	handleResponse(this.Ctx.ResponseWriter,200,response)

}

func (this *AuthController) RecordAuth () {

	// 接收上传的文件
	beego.Info("receive file")

	// 定义读取文件的key
	var key = "auth"
	// auth是一个key值，需要与前端约定
	file, header, err := this.GetFile(key)
	if err != nil {
		// 如果出错，没有读到文件
		handleResponse(this.Ctx.ResponseWriter,400,err.Error())
		return
	}
	defer file.Close()
	// 保存文件，将文件保存在static/upload目录下面
	fileName := header.Filename
	beego.Info("文件名称：",fileName)

	// 关于保存文件的路径，static前不能添加/.(这里指在当前项目的static/upload文件夹下)
	err = this.SaveToFile(key, path.Join("static/upload", fileName))
	if err != nil {
		handleResponse(this.Ctx.ResponseWriter,500,err.Error())
		return
	}
	// 回复用户已经接收上传的文件


	// 开启任务：toolbox.NewTask()
	// 指定任务的开启时间，保存完文件后的5秒钟执行写数据的任务
	// 任务是被放到容器中进行管理的，容器的开启和关闭
	taskName := "t1"
	t := time.Now().Add(5*time.Second)
	second := t.Second()
	minute := t.Minute()
	hour := t.Hour()
	spec := fmt.Sprintf("%d %d %d * * *", second, minute, hour)

	task := toolbox.NewTask(taskName,spec, func() error {
		beego.Info("task start")
		// 当任务执行完成后，停止
		defer toolbox.StopTask()
		return myTask(fileName)
		
	})
	
	// 注意:
	// task.Run() 立即执行
	// 将任务添加到容器，容器中可以有多个任务
	toolbox.AddTask(taskName,task)
	// 开启任务执行
	toolbox.StartTask()
	handleResponse(this.Ctx.ResponseWriter,200,"ok")
}

// 耗时操作
func myTask(fileName string) error {
	// 初始化ccs
	channelId := beego.AppConfig.String("channel_id_gaj")
	chaincodeId := beego.AppConfig.String("chaincode_id_auth")
	userId := beego.AppConfig.String("user_id")
	conf := beego.AppConfig.String("CORE_OGAJ_CONFIG_FILE")

	ccs, err := models.Initialize(channelId, chaincodeId, userId,conf)
	if err != nil {
		beego.Error(err.Error())
	}
	defer ccs.Close()
	// 读文件
	file, _ := os.Open(path.Join("static/upload",fileName))
	defer file.Close()
	reader := csv.NewReader(file)

	// 并没有终止，原因:可能某一行数据有问题,并不是所有的都有问题，只需要记录有问题的航信息就可以了
	// 还有一种可能:某行数据写入区块链出错
	// 记录行数：可能有多个行数据出问题，不会使用行的字符串拼接方式
	// 会定义一个字符串的数组，最终生成字符串时，可以使用","进行分割行数
	var line = 0
	var lines []string

	for {
		line++
		lineStr := strconv.Itoa(line)
		record, err := reader.Read()
		if err == io.EOF{
			break
		}
		if err != nil{
			//  有异常需要
			lines = append(lines, lineStr)
			continue

		}

		if len(record) != 3 {
			//  有异常需要处理
			lines = append(lines, lineStr)
			continue
		}

		var args [][]byte


		for _, value := range record {
			// fmt.Print(value,"\t")
			args = append(args, []byte(value))
		}
		//fmt.Println()
		_, err = ccs.ChainCodeUpdate("add", args)
		if err != nil{
			// 有异常需要处理
			lines = append(lines, lineStr)
		}

	}
	if len(lines) > 0 {
		beego.Error("Error lines:",strings.Join(lines,","))
	} else {
		// 执行的一切顺利
		beego.Info("write data success")

	}
	return nil
	// 将读取的每条交易记录信息都写到区块中
	
}





