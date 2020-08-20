package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/adschain/log"
	"github.com/wlcy/tradehome-service/util"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

var dbWriteHost = ""
var dbWritePort = ""
var dbWriteSchema = ""
var dbWriteName = ""
var dbWritePass = ""

//连接DB的实例对象
var dbWriteInstance *TronDB

// InitializeWriter 初始化
func InitializeWriter(host, port, schema, username, password string) {
	dbWriteHost = strings.TrimSpace(host)
	dbWritePort = strings.TrimSpace(port)
	dbWriteSchema = strings.TrimSpace(schema)
	dbWriteName = strings.TrimSpace(username)
	dbWritePass = strings.TrimSpace(password)
	fmt.Printf("InitializeWriter dbWriteHost:%v, dbWritePort:%v, dbWriteSchema:%v, dbWriteName:%v, dbWritePass:%v\n",
		dbWriteHost, dbWritePort, dbWriteSchema, dbWriteName, dbWritePass)
}

//GetWriteDatabase Get一个连接的数据库对象
func GetWriteDatabase() (*TronDB, error) {
	return retrieveWriteDatabase()
}

//retrieveWriteDatabase 刷新DB的连接
func retrieveWriteDatabase() (*TronDB, error) {
	defer util.CatchError()

	if nil == dbWriteInstance {
		//连接数据库的参数
		para := GetMysqlWriteConnectionInfo()

		//打开这个DB对象
		dbPtr, err := OpenDB(para.Mode, para.ConnSQL)
		if err != nil {
			return nil, err
		}
		if dbPtr == nil {
			return nil, errors.New("db not connected")
		}

		//设置连接池信息
		dbPtr.SetConnsParam(para.MaxOpenconns, para.MaxIdleConns)
		dbWriteInstance = dbPtr
	}

	//测试一下是否是连接成功的
	if err := dbWriteInstance.Ping(); err != nil {
		dbWriteInstance = nil
		return nil, err
	}

	return dbWriteInstance, nil
}

//GetMysqlWriteConnectionInfo 获取连接mysql的相关信息
func GetMysqlWriteConnectionInfo() DBParam {
	dbConfig := DBParam{
		Mode:         string("mysql"),
		ConnSQL:      fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8", dbWriteName, dbWritePass, dbWriteHost, dbWritePort, dbWriteSchema),
		MaxOpenconns: 10,
		MaxIdleConns: 10,
	}
	log.Infof("GetMysqlWriteConnectionInfo ConnSQL:%v", dbConfig.ConnSQL)
	return dbConfig
}

// OpenDataBaseWriteTransaction 开启一个数据库事物
func OpenDataBaseWriteTransaction() (*sql.Tx, error) {
	dataPtr, err := GetWriteDatabase()
	if err != nil {
		return nil, err
	}

	sqlTx, err := dataPtr.Begin()
	if err != nil {
		return nil, err
	}
	return sqlTx, nil
}

// OpenDBPrepare statement
func OpenDBWritePrepare(query string) (*sql.Tx, *sql.Stmt, error) {
	dataPtr, err := GetWriteDatabase()
	if err != nil {
		return nil, nil, err
	}

	sqlTx, err := dataPtr.Begin()
	if err != nil {
		return nil, nil, err
	}

	stmt, err := sqlTx.Prepare(query)
	if err != nil {
		return nil, nil, err
	}

	return sqlTx, stmt, nil
}

//ExecuteSQLCommand 执行insert update, delete操作,依次返回 插入消息的主键，影响的条数，错误对象
func ExecuteSQLCommand(strSQL string) (int64, int64, error) {
	var key int64
	var rows int64
	var err error
	var dbPtr *TronDB

	//获取数据库对象
	if dbPtr, err = GetWriteDatabase(); err != nil {
		return 0, 0, errors.New("db not connected")
	}

	//执行语句
	if key, rows, err = dbPtr.Execute(strSQL); err != nil {
		return key, rows, errors.New("execute sql command error")
	}

	return key, rows, err
}
