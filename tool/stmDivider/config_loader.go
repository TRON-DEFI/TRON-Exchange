package main

import (
	"encoding/base64"
	"github.com/adschain/log"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
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

	//initSmart
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

	//logrus.SetReportCaller(true)
	//formatter := &logrus.TextFormatter{
	//	TimestampFormat: "2006-01-02T15:04:05.000000-0700 MST",
	//	FullTimestamp:   true,
	//	CallerPrettyfier: func(f *runtime.Frame) (string, string) {
	//		return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", path.Base(f.File), f.Line)
	//	},
	//}
	//logrus.SetFormatter(formatter)
	//logFile := &lumberjack.Logger{
	//	Filename:   viper.GetString("log.logger_file"),
	//	MaxSize:    viper.GetInt("log.log_rotate_size"), // megabytes
	//	MaxBackups: viper.GetInt("log.log_backup_count"),
	//	MaxAge:     30,   //days
	//	Compress:   true, // disabled by default
	//	LocalTime:  true,
	//}
	//mw := io.MultiWriter(os.Stdout, logFile)
	//logrus.SetOutput(mw)
	//level, err := logrus.ParseLevel(viper.GetString("log.logger_level"))
	//if nil != err {
	//	level = logrus.InfoLevel
	//}
	//logrus.SetLevel(level)
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
	ContractAddress = viper.GetString("smart.contractAddress")
	OwnerAddr = viper.GetString("smart.ownerAddress")
	encryptPrivateKey := viper.GetString("smart.privateKey")
	log.Infof("SmartContractAddr:%v, SmartOwnerAddr:%v, SmartPrivateKey:%v", ContractAddress, OwnerAddr, "hehe")
	if err := InitSmart(); err != nil {
		return err
	}
	key := []byte("dblyztradestmvid")
	bytesPass, err := base64.StdEncoding.DecodeString(encryptPrivateKey)
	if err != nil {
		return err
	}
	priKey, err := util.AesDecrypt(bytesPass, key)
	if err != nil {
		return err
	}
	PrivateKey = string(priKey)

	PoolAddress = viper.GetString("smart.poolAddress")
	encryptPoolPrivateKey := viper.GetString("smart.poolPrivateKey")
	bytesPoolPass, err := base64.StdEncoding.DecodeString(encryptPoolPrivateKey)
	if err != nil {
		return err
	}
	priPoolKey, err := util.AesDecrypt(bytesPoolPass, key)
	if err != nil {
		return err
	}
	PoolPrivateKey = string(priPoolKey)
	log.Infof("PoolAddress:%v, PoolPrivateKey:%v", PoolAddress, "hehe")
	return nil
}
