package main

import (
	"github.com/KeepMovingLr/taskserver/cache"
	"github.com/KeepMovingLr/taskserver/conf"
	"github.com/KeepMovingLr/taskserver/dao"
	"github.com/KeepMovingLr/taskserver/server"
	"os"
	"runtime/pprof"
)

func main() {
	// add cpu profile monitor
	cpuProfile, _ := os.Create("cpu_profile")
	defer cpuProfile.Close()
	// do not deal with err
	_ = pprof.StartCPUProfile(cpuProfile)
	defer pprof.StopCPUProfile()

	// init global config variable
	if err := conf.InitializeConfig(); err != nil {
		server.ExitOnError(err)
	}
	globalConfig := conf.GetGlobalConfig()

	// init DB
	dbConfig := globalConfig.DBconf
	param := dao.MysqlInitParam{
		dbConfig.DBUser, dbConfig.DBPwd, dbConfig.DBHost, dbConfig.DBName, dbConfig.MaxOpenConns, dbConfig.MaxIdleConns,
	}
	if err := dao.ConnectMysql(param); err != nil {
		server.ExitOnError(err)
	}

	// init Local cache
	localCacheConfig := globalConfig.LocalCacheConfig
	if err := cache.InitMyLocalCache(localCacheConfig.LocalCacheSize); err != nil {
		server.ExitOnError(err)
	}

	// open listener
	TCPConfig := globalConfig.TCPConfig
	server.RegisterServerListener(TCPConfig.Port, TCPConfig.ReadWaitTime, TCPConfig.WriteWaitTime)

}
