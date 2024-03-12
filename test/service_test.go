package test

import (
	"ray.li/entrytaskserver/cache"
	"ray.li/entrytaskserver/conf"
	"ray.li/entrytaskserver/coreservice"
	"ray.li/entrytaskserver/dao"
	"ray.li/entrytaskserver/server"
	"ray.li/entrytaskserver/utils"
	"testing"
)

func initDependence() {
	cache.InitMyLocalCache(10)
	// must init before test
	globalConfig := conf.GetGlobalConfig()
	globalConfig.FromFile("test_config.xml")
	dbConfig := globalConfig.DBconf
	// init db
	param := dao.MysqlInitParam{
		dbConfig.DBUser, dbConfig.DBPwd, dbConfig.DBHost, dbConfig.DBName, dbConfig.MaxOpenConns, dbConfig.MaxIdleConns,
	}
	if err := dao.ConnectMysql(param); err != nil {
		server.ExitOnError(err)
	}
}

func TestLoginAuthenticate(t *testing.T) {
	// Table-Driven Test
	tests := []struct{ username, password string }{
		{"ray", "pwd"},
		{"ray", "pwd"},
		//{"ray2", "pwd"},
	}
	initDependence()
	utils.Sha256Encode([]byte("pwd"))
	for _, tt := range tests {
		if actual, _ := coreservice.LoginAuthenticate(tt.username, tt.password); !actual.Success {
			t.Errorf("wrong")
		} else {
			t.Log("success")
		}
	}
}

// benchmark test
func BenchmarkLoginAuthenticate(b *testing.B) {
	initDependence()
	for i := 0; i < b.N; i++ {
		actual, _ := coreservice.LoginAuthenticate("ray", "pwd")
		if !actual.Success {
			b.Errorf("false")
		}
	}
}
