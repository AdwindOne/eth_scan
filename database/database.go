/*
Copyright 2018 The go-eam Authors
This file is part of the go-eam library.

database
封装数据库相关操作


wanglei.ok@foxmail.com

1.0
版本时间：2018年4月13日18:32:12

*/

package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"fmt"
	"eth_scan/etherscan_api"
)

//数据库操作对象
var db *sql.DB

const (
	//连接池属性
	POOL_MAXOPENCONNS = 10	//最大连接数
	POOL_MAXIDLECONNS = 2	//空闲连接数
)

//打开数据库
func OpenDatabase(dsn string) error {
	db1, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	//连接池
	db1.SetMaxOpenConns(POOL_MAXOPENCONNS)
	db1.SetMaxIdleConns(POOL_MAXIDLECONNS)
	//连接
	if  err = db1.Ping(); err != nil {
		return err
	}
	db = db1
	return nil
}

//关闭数据库
func CloseDatabase() {
	db.Close()
}

//自定义事务结构
type MyTx struct {
	Tx *sql.Tx
}

//开始事务并返回事务对象
func TxBegin() (*MyTx, error) {
	tx, err := db.Begin()
	return &MyTx{tx}, err
}

//提交事务
func (x *MyTx)Commit() error {
	return x.Tx.Commit()
}

//回滚事务
func (x *MyTx)Rollback() error {
	return x.Tx.Rollback()
}

//使用事务插入一条交易记录
func (x *MyTx) InsertTx(tx *etherscan_api.TxJson) error {
	stmt, err := x.Tx.Prepare("INSERT ec_ethdata SET block_number=?,time_stamp=?,tx_hash=?, nonce=?, block_hash=?, tx_index=?, from_addr=?, to_addr=?, contract_addr=?, amount=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(tx.BlockNumber, tx.TimeStamp, tx.Hash, tx.Nonce, tx.BlockHash, tx.TransactionIndex, tx.From, tx.To, tx.ContractAddress, tx.Value )
	if err != nil {
		return err
	}

	return err
}

//使用事务更新最后一个块计数
//返回收到影响的记录数
func (x *MyTx) UpdateLastBlock(addr string, block int) (affect int64) {

	affect = 0

	stmt, err := x.Tx.Prepare("update ec_address_log set last_block = ? where address = ?")
	if err != nil {
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(strconv.Itoa(block), addr)
	if err != nil {
		return
	}

	affect, err = res.RowsAffected()

	return
}

type EthAddressInfo struct {
	Address string
	LastBlock int
}

//获取以太坊地址及最后块计数器
func GetEthAddress() ([]EthAddressInfo, error){

	eais := make([]EthAddressInfo,0)

	//查询数据
	rows, err := db.Query("select address, last_block from ec_address_log where type = 'eth'")
	fmt.Printf("d%",rows)
	if err != nil {
		return eais, err
	}

	for rows.Next() {
		a := ""
		b := 0
		err = rows.Scan(&a, &b)
		if err != nil {
			return eais, err
		}
		eais = append(eais, EthAddressInfo{a,b})
	}

	return eais, nil
}


func UpdateAddressLogByCrowdOrder() error {
	//-- 用订单数据更新监控地址列表
	//-- 1.订单列表中不存在的地址,不再监控
	//UPDATE ec_address_log SET state = '1' WHERE type='eth' AND address NOT IN ( SELECT DISTINCT pay_url AS address FROM `ec_crowd_order` WHERE order_type='eth' AND is_delete=0 AND  pay_status IN (0,2) );
	//
	//-- 2. 订单列表中的地址,开始监控
	//UPDATE ec_address_log SET state = '0' WHERE type='eth' AND address IN ( SELECT DISTINCT pay_url AS address FROM `ec_crowd_order` WHERE order_type='eth' AND is_delete=0 AND  pay_status IN (0,2) );
	//
	//-- 3. 插入新的监控地址
	//INSERT INTO ec_address_log (address,type,state,last_block) SELECT DISTINCT pay_url AS address, 'eth' AS type, '0' AS state, 0 AS last_block FROM `ec_crowd_order` WHERE order_type='eth' AND is_delete=0 AND  pay_status IN (0,2) AND NOT EXISTS (SELECT address FROM ec_address_log WHERE address = ec_crowd_order.pay_url);


	var sqls = [...]string {
		"UPDATE ec_address_log SET state = '1' WHERE type='eth' AND address NOT IN ( SELECT DISTINCT pay_url AS address FROM `ec_crowd_order` WHERE order_type='eth' AND is_delete=0 AND  pay_status IN (0,2) );",
		"UPDATE ec_address_log SET state = '0' WHERE type='eth' AND address IN ( SELECT DISTINCT pay_url AS address FROM `ec_crowd_order` WHERE order_type='eth' AND is_delete=0 AND  pay_status IN (0,2) );",
		"INSERT INTO ec_address_log (address,type,state,last_block) SELECT DISTINCT pay_url AS address, 'eth' AS type, '0' AS state, 0 AS last_block FROM `ec_crowd_order` WHERE order_type='eth' AND is_delete=0 AND  pay_status IN (0,2) AND NOT EXISTS (SELECT address FROM ec_address_log WHERE address = ec_crowd_order.pay_url);",
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for _, sql := range sqls {
		_, err := tx.Exec(sql)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}


func UpdateAddressLogBySetting() error {
	//-- 用设置更新监控地址列表
	//-- 1.设置中不存在地址,不再监控
	//UPDATE ec_address_log SET state = '1' WHERE type='eth' AND LOWER(address) NOT IN ( SELECT DISTINCT LOWER(set_value) FROM ec_setting WHERE set_key = 'eth_acceptAddress' );
	//
	//-- 2. 设置中的地址,开始监控
	//UPDATE ec_address_log SET state = '0' WHERE type='eth' AND LOWER(address) IN ( SELECT DISTINCT LOWER(set_value) FROM ec_setting WHERE set_key = 'eth_acceptAddress' );
	//
	//-- 3. 插入新的监控地址
	//INSERT INTO ec_address_log (address,type,state,last_block) SELECT DISTINCT set_value AS address, 'eth' AS type, '0' AS state, 0 AS last_block FROM ec_setting WHERE set_key = 'eth_acceptAddress' AND NOT EXISTS (SELECT LOWER(address) FROM ec_address_log WHERE LOWER(address) = LOWER(ec_setting.set_value)) LIMIT 1;

	var sqls = [...]string {
		"UPDATE ec_address_log SET state = '1' WHERE type='eth' AND LOWER(address) NOT IN ( SELECT DISTINCT LOWER(set_value) FROM ec_setting WHERE set_key = 'eth_acceptAddress' );",
		"UPDATE ec_address_log SET state = '0' WHERE type='eth' AND LOWER(address) IN ( SELECT DISTINCT LOWER(set_value) FROM ec_setting WHERE set_key = 'eth_acceptAddress' );",
		"INSERT INTO ec_address_log (address,type,state,last_block) SELECT DISTINCT set_value AS address, 'eth' AS type, '0' AS state, 0 AS last_block FROM ec_setting WHERE set_key = 'eth_acceptAddress' AND NOT EXISTS (SELECT LOWER(address) FROM ec_address_log WHERE LOWER(address) = LOWER(ec_setting.set_value)) LIMIT 1;",
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for _, sql := range sqls {
		_, err := tx.Exec(sql)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}