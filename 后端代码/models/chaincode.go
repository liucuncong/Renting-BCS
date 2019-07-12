package models

import (
	"github.com/hyperledger/fabric-sdk-go/api/apitxn"
		"github.com/hyperledger/fabric-sdk-go/def/fabapi"
	"github.com/astaxie/beego"
	"fmt"
)

type ChainCodeSpec struct {
	client apitxn.ChannelClient
	chaincodeID string
}

func Initialize(channelID,chaincodeID,userID string,config string) (*ChainCodeSpec,error) {
	//config := beego.AppConfig.String("CORE_OGAJ_CONFIG_FILE")
	fmt.Println("!!!!!!!:",config)

	fabricSDK, err := getSDK(config)
	if err != nil {
		return nil,err
	}

	client, err := fabricSDK.NewChannelClient(channelID, userID)
	if err != nil {
		return nil,err
	}
	return &ChainCodeSpec{client,chaincodeID}, nil

}

func (this *ChainCodeSpec)ChainCodeQuery(fun string,args [][]byte) (response []byte,err error) {
	queryRequest := apitxn.QueryRequest{this.chaincodeID, fun, args}
	return this.client.Query(queryRequest)
}

func (this *ChainCodeSpec)ChainCodeUpdate(fun string,args [][]byte) (response []byte,err error) {
	request := apitxn.ExecuteTxRequest{ChaincodeID: this.chaincodeID, Fcn: fun, Args: args}
	id, err := this.client.ExecuteTx(request)
	return []byte(id.ID),err
}
func (this *ChainCodeSpec) Close() {
	this.client.Close()
}

func getSDK(config string) (*fabapi.FabricSDK,error) {
	options := fabapi.Options{ConfigFile: config}

	fabricSDK, err := fabapi.NewSDK(options)
	if err != nil {
		beego.Error(err)
		return nil,err
	}
	return fabricSDK,nil
}