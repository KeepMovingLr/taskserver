package conf

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
)

type Config struct {
	XMLName          xml.Name `xml:"config"`
	DBconf           DBConfig
	LocalCacheConfig LocalCacheConfig
	TCPConfig        TCPConfig
}

type DBConfig struct {
	XMLName      xml.Name `xml:"DBConfig"`
	MaxOpenConns int      `xml:"MaxOpenConns"`
	MaxIdleConns int      `xml:"MaxIdleConns"`
	DBName       string   `xml:"DBName"` // User
	DBUser       string   `xml:"DBUser"` // root
	DBPwd        string   `xml:"DBPwd"`  // rootroot
	DBHost       string   `xml:"DBHost"` // 127.0.0.1:3306
}

type LocalCacheConfig struct {
	XMLName        xml.Name `xml:"LocalCacheConfig"`
	LocalCacheSize int      `xml:"LocalCacheSize"`
}

type TCPConfig struct {
	XMLName       xml.Name `xml:"TCPConfig"`
	Port          string   `xml:"Port"`
	ReadWaitTime  int      `xml:"ReadWaitTime"`
	WriteWaitTime int      `xml:"WriteWaitTime"`
}

func (config *Config) FromFile(filename string) (err error) {
	configFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	err = xml.Unmarshal(configFile, config)
	return
}

func InitializeConfig() (err error) {
	if err := globalCfg.FromFile("./conf/config.xml"); err != nil {
		return errors.New("initialize config error" + err.Error())
	}
	return nil
}
