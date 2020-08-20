package config

import (
	"encoding/base64"
	"github.com/adschain/log"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/wlcy/tradehome-service/common/contract"
	"github.com/wlcy/tradehome-service/common/mysql"
	"github.com/wlcy/tradehome-service/common/redis"
	"github.com/wlcy/tradehome-service/util"
	"strings"
)

type Config struct {
	Name string
}

func Init(cfg string) error {
	c := Config{
		Name: cfg,
	}

	// 初始化配置文件
	if err := c.initConfig(); err != nil {
		return err
	}

	// 初始化日志包
	c.initLog()

	// 监控配置文件变化并热加载程序
	c.watchConfig()

	// initDB
	c.initDB()

	// initRedis
	c.initRedis()

	// initSmart
	if err := c.initSmart(); err != nil {
		return err
	}

	return nil
}

func (c *Config) initConfig() error {
	if c.Name != "" {
		viper.SetConfigFile(c.Name) // 如果指定了配置文件，则解析指定的配置文件
	} else {
		viper.AddConfigPath("conf") // 如果没有指定配置文件，则解析默认的配置文件
		viper.SetConfigName("config")
	}
	viper.SetConfigType("yaml") // 设置配置文件格式为YAML
	viper.AutomaticEnv()        // 读取匹配的环境变量
	//viper.SetEnvPrefix("DAPPHOUSE") // 读取环境变量的前缀DAPPHOUSE
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	if err := viper.ReadInConfig(); err != nil { // viper解析配置文件
		return err
	}

	return nil
}

func (c *Config) initLog() {
	passLagerCfg := log.PassLagerCfg{
		Writers:        viper.GetString("log.writers"),
		LoggerLevel:    viper.GetString("log.logger_level"),
		LoggerFile:     viper.GetString("log.logger_file"),
		LogFormatText:  viper.GetBool("log.log_format_text"),
		RollingPolicy:  viper.GetString("log.rollingPolicy"),
		LogRotateDate:  viper.GetInt("log.log_rotate_date"),
		LogRotateSize:  viper.GetInt("log.log_rotate_size"),
		LogBackupCount: viper.GetInt("log.log_backup_count"),
	}

	log.InitWithConfig(&passLagerCfg)
}

// 监控配置文件变化并热加载程序
func (c *Config) watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("Config file changed: %s", e.Name)
	})
}

// initDB
func (c *Config) initDB() {
	writeHost := viper.GetString("db_write.host")
	writePort := viper.GetString("db_write.port")
	writeSchema := viper.GetString("db_write.schema")
	writeUsername := viper.GetString("db_write.username")
	writePassword := viper.GetString("db_write.password")
	log.Infof("InitializeWriter dbWriteHost:%v, dbWritePort:%v, dbWriteSchema:%v, dbWriteName:%v, dbWritePass:%v",
		writeHost, writePort, writeSchema, writeUsername, writePassword)
	mysql.InitializeWriter(writeHost, writePort, writeSchema, writeUsername, writePassword)

	readHost := viper.GetString("db_read.host")
	readPort := viper.GetString("db_read.port")
	readSchema := viper.GetString("db_read.schema")
	readUsername := viper.GetString("db_read.username")
	readPassword := viper.GetString("db_read.password")
	log.Infof("InitializeRead dbReadHost:%v, dbReadPort:%v, dbReadSchema:%v, dbReadName:%v, dbReadPass:%v",
		readHost, readPort, readSchema, readUsername, readPassword)
	mysql.InitializeRead(readHost, readPort, readSchema, readUsername, readPassword)

}

// initRedis
func (c *Config) initRedis() {
	redisInfo := struct {
		Addr     string
		Password string
		Db       int
		PoolSize int
	}{}
	redisInfo.Addr = viper.GetString("redis.host")
	redisInfo.Password = viper.GetString("redis.pass")
	redisInfo.Db = int(util.ToInt64(viper.GetString("redis.index")))
	redisInfo.PoolSize = int(util.ToInt64(viper.GetString("redis.poolSize")))
	log.Infof("redis init:Addr:%v, Password:%v, Db:%v, PoolSize:%v", redisInfo.Addr, redisInfo.Password, redisInfo.Db, redisInfo.PoolSize)
	redis.RedisCli = redis.NewClient(redisInfo.Addr, redisInfo.Password, redisInfo.Db, redisInfo.PoolSize)
}

// initSmart
func (c *Config) initSmart() error {
	contract.TestNet = viper.GetInt("smart.testNet")
	log.Infof("TestNet:%v", contract.TestNet)

	contract.EventLogURL = viper.GetString("smart.eventLogURL")
	contract.SmartNode = viper.GetString("smart.smartNode")
	contract.SmartContractAddr = viper.GetString("smart.contractAddress")
	contract.SmartOwnerAddr = viper.GetString("smart.ownerAddress")
	contract.SmartPrivateKey = viper.GetString("smart.privateKey")
	contract.SmartContractAddr10 = viper.GetString("smart.contractAddress10")
	contract.SmartOwnerAddr10 = viper.GetString("smart.ownerAddress10")
	contract.SmartPrivateKey10 = viper.GetString("smart.privateKey10")
	log.Infof("EventLogURL:%v, SmartNode:%v", contract.EventLogURL, contract.SmartNode)
	log.Infof("SmartContractAddr:%v, SmartOwnerAddr:%v, SmartPrivateKey:%v", contract.SmartContractAddr, contract.SmartOwnerAddr, contract.SmartPrivateKey)
	log.Infof("SmartContractAddr10:%v, SmartOwnerAddr10:%v, SmartPrivateKey10:%v", contract.SmartContractAddr10, contract.SmartOwnerAddr10, contract.SmartPrivateKey10)
	if err := contract.InitTradeSmart(); err != nil {
		return err
	}

	contract.MineContractAddress = viper.GetString("mine.contractAddress")
	contract.MineRate = viper.GetInt64("mine.rate")
	contract.MineOwnerAddr = viper.GetString("mine.ownerAddress")
	key := []byte("321423u9y8d2tron")
	encryptPrivateKey := viper.GetString("mine.privateKey")
	bytesPass, err := base64.StdEncoding.DecodeString(encryptPrivateKey)
	if err != nil {
		return err
	}
	priKey, err := util.AesDecrypt(bytesPass, key)
	if err != nil {
		return err
	}
	contract.MinePrivateKey = string(priKey)
	log.Infof("MineContractAddress:%v, MineRate:%v, MineOwnerAddr:%v, MinePrivateKey:%v", contract.MineContractAddress, contract.MineRate, contract.MineOwnerAddr, contract.MinePrivateKey)
	/*if err := contract.InitMineSmart(); err != nil {
		return err
	}*/

	return nil
}
