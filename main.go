/**
 * Created by Adwind.
 * User: liuyunlong
 * Date: 11/7/18
 * Time: 16:36
 */

package main

import (
	"strconv"
	"time"
	"strings"
	"log"
	"github.com/wanglei-ok/logfile"

	"eth_scan/config"
	"eth_scan/database"
	"eth_scan/etherscan_api"
)


func init() {
	logfile.Setup()
}

func main() {

	if err := config.ReadConfig(); err != nil {
		log.Println("Error readConfig:", err)
		return
	}


	//打开数据库
	if err := database.OpenDatabase(config.Config1.Mysql.DSN); err != nil {
		log.Println("Error OpenDatabase:", err.Error())
		return
	}
	//程序结束前关闭数据库
	defer database.CloseDatabase()

	////用订单数据更新监控地址列表
	//if err := UpdateAddressLogByCrowdOrder(); err != nil {
	//	log.Println("Error UpdateAddressLogByCrowdOrder:", err.Error() )
	//	return
	//}
	//用网站配置更新监控地址列表
	if err := database.UpdateAddressLogByCrowdOrder(); err != nil {
		log.Println("Error UpdateAddressLogBySetting:", err.Error() )
		return
	}

	//for  {
	//取得地址和最后一次更新的块号
	eais, err := database.GetEthAddress()
	if err != nil {
		log.Println("Error GetEthAddress:", err.Error())
		return
	}

	if len(eais) == 0 {
		log.Println("Error not found ethereum address.")
		return
	}
	for _,eai := range(eais) {
		retrieve(eai.Address,eai.LastBlock)
	}
	}
//}

//取回指定地址和开始块号
func retrieve(addr string, startBlock int) {
	proc := 0
	inc := 0
	skip := 0

	t1 := time.Now()
	log.Printf("Begin retrieve.(Address: %s%s, StartBlock:%d)", config.Config1.EtherscanApi.ApiAddress, addr, startBlock)
	defer func (){
		elapsed := time.Since(t1)
		log.Printf("End retrieve.(Elapsed:%v, Process:%d, Increase:%d, Skip:%d)", elapsed, proc, inc, skip)
	}()

	maxBlock := startBlock

	//使用Etherscan API检索交易列表
	txlistJson, err := etherscan_api.Retrieve(addr, startBlock, true)

	//检索失败处理
	if err != nil {
		log.Println("Error retrieve:", err.Error())
		return
	}

	//API返回错误处理
	if txlistJson.Status != "1" {
		log.Println("Etherscan api:", txlistJson.Message)
		return
	}

	//没有交易记录
	if txlistJson.Result == nil {
		return
	}

	//遍历交易记录插入数据库
	//开启事务
	trans, err := database.TxBegin()
	if err != nil {
		log.Println("Error TxBegin:", err.Error())
		return
	}
	for _ , tx := range txlistJson.Result {
		proc++
		err := trans.InsertTx(&tx)
		if err != nil {
			//插入失败显示日志
			//txString, _ := json.Marshal(tx)
			errString := err.Error()
			if strings.Contains(errString, "Duplicate entry") {
				log.Printf("Skip Duplicate tx: %s%s", config.Config1.EtherscanApi.ApiTx,tx.Hash)
			} else {
				log.Printf("Error insertTx:%v, %s%s", errString, config.Config1.EtherscanApi.ApiTx, tx.Hash)
			}
			skip++
		}else{
			inc++
			log.Printf("Increase tx: %s%s", config.Config1.EtherscanApi.ApiTx,tx.Hash)
		}


		//最后块编号
		b, ok := strconv.Atoi(tx.BlockNumber )
		if ok == nil {
			if b > maxBlock {
				maxBlock = b
			}
		}
	}

	//保存最后块编号
	if maxBlock > startBlock {
		if trans.UpdateLastBlock(addr, maxBlock) == 0 {
			trans.Rollback()
		}
	}

	trans.Commit()
}
