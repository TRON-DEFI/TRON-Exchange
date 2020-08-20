package mysql

import (
	"errors"
	"fmt"
	"github.com/adschain/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/wlcy/tradehome-service/util"
	"strconv"
	"strings"
)

var dbReadHost = ""
var dbReadPort = ""
var dbReadSchema = ""
var dbReadName = ""
var dbReadPass = ""

//数据库的连接配置
type DBParam struct {
	Mode         string
	ConnSQL      string
	MaxOpenconns int
	MaxIdleConns int
}

//连接DB的实例对象
var dbInstance *TronDB

// Initialize 初始化
func InitializeRead(host, port, schema, username, password string) {
	dbReadHost = strings.TrimSpace(host)
	dbReadPort = strings.TrimSpace(port)
	dbReadSchema = strings.TrimSpace(schema)
	dbReadName = strings.TrimSpace(username)
	dbReadPass = strings.TrimSpace(password)
	fmt.Printf("InitializeRead dbReadHost:%v, dbReadPort:%v, dbReadSchema:%v, dbReadName:%v, dbReadPass:%v\n",
		dbReadHost, dbReadPort, dbReadSchema, dbReadName, dbReadPass)
}

//GetReadDatabase Get一个连接的数据库对象
func GetReadDatabase() (*TronDB, error) {
	return retrieveReadDatabase()
}

//retrieveReadDatabase 刷新DB的连接
func retrieveReadDatabase() (*TronDB, error) {
	defer util.CatchError()
	if nil == dbInstance {
		//连接数据库的参数
		para := GetMysqlReadConnectionInfo()
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
		dbInstance = dbPtr
	}
	//测试一下是否是连接成功的
	if err := dbInstance.Ping(); err != nil {
		//dbInstance.Close()
		dbInstance = nil
		return nil, err
	}
	return dbInstance, nil
}

//GetMysqlConnectionInfo 获取连接mysql的相关信息
func GetMysqlReadConnectionInfo() DBParam {
	dbConfig := DBParam{
		Mode:         string("mysql"),
		ConnSQL:      fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8", dbReadName, dbReadPass, dbReadHost, dbReadPort, dbReadSchema),
		MaxOpenconns: 10,
		MaxIdleConns: 10,
	}
	log.Infof("GetMysqlReadConnectionInfo ConnSQL:%v", dbConfig.ConnSQL)
	return dbConfig
}

//QueryTableData 查询数据库数据
func QueryTableData(strSQL string) (*TronDBRows, error) {
	//获取数据库对象
	var dbPtr *TronDB
	var err error
	if dbPtr, err = GetReadDatabase(); err != nil {
		return nil, errors.New("get database error")
	}

	//查询数据集
	rows, err := dbPtr.Select(strSQL)
	if err != nil {
		return rows, err
	}
	return rows, err
}

//QueryTableDataCount 返回某个表的记录个数
func QueryTableDataCount(tableName string) (int64, error) {
	rowCount := int64(0) //返回的数据库表行数

	//判断输入参数
	if len(tableName) < 0 {
		return 0, errors.New("parameter invalid")
	}

	//获取数据库对象
	var dbPtr *TronDB
	var err error
	if dbPtr, err = GetReadDatabase(); err != nil {
		return 0, errors.New("db not connected")
	}

	strSQL := "select count(*) as rowcounts from " + tableName
	var data *TronDBRows
	if data, err = dbPtr.Select(strSQL); err != nil {
		return 0, errors.New("query table data count error")
	}

	if data.NextT() {
		strValue := data.GetField("rowcounts")
		if count, err := strconv.ParseInt(strValue, 10, 64); err != nil {
			return 0, errors.New("query table data count error")
		} else {
			rowCount = count //set the count
		}
	}

	return rowCount, nil
}

//QueryTableDataCount 返回某个表的记录个数
func QuerySQLViewCount(strSQLView string) (int64, error) {
	rowCount := int64(0) //返回的数据库表行数

	//判断输入参数
	if len(strSQLView) < 0 {
		return 0, errors.New("parameter invalid")
	}

	//获取数据库对象
	var dbPtr *TronDB
	var err error
	if dbPtr, err = GetReadDatabase(); err != nil {
		return 0, errors.New("db not connected")
	}

	strSQL := fmt.Sprintf(`
		select count(*) as rowcounts from (
			%s
		) newtableName`, strSQLView)
	var data *TronDBRows
	if data, err = dbPtr.Select(strSQL); err != nil {
		return 0, errors.New("query sql view count error")
	}

	if data.NextT() {
		strValue := data.GetField("rowcounts")
		if count, err := strconv.ParseInt(strValue, 10, 64); err != nil {
			return 0, errors.New("query sql view count error")
		} else {
			rowCount = count //set the count
		}
	}

	return rowCount, nil
}
