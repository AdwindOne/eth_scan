/**
 * Created by Adwind.
 * User: liuyunlong
 * Date: 11/7/18
 * Time: 16:36
 */
package config

import (
	"gopkg.in/gcfg.v1"
)

var Config1 Config

type (
	Config struct{
		Mysql struct{
			DSN string
		}

		EtherscanApi struct {
			ApiTxlist string
			ApiAddress string
			ApiTx string
		}
	}
)

//func iniFileName() string {
//	exePath := os.Args[0]
//	fmt.Println(exePath)
//	base := filepath.Base(exePath)
//	fmt.Println(base)
//	suffix := filepath.Ext(exePath)
//	return strings.TrimSuffix(base, suffix) + ".ini"
//}

func ReadConfig() error {
	//return gcfg.ReadFileInto(&config, iniFileName())
	return gcfg.ReadFileInto(&Config1, "eth_scan.ini")
}
