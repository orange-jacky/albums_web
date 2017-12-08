package util

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"sync"
)

//Gin server
type GinServer struct {
	Mode            string `xml:"mode"`
	Url             string `xml:"url"`
	Port            string `xml:"port"`
	Timeout_read_s  int    `xml:"timeout_read_s"`
	Timeout_write_s int    `xml:"timeout_write_s"`
}

//backend
type Backend struct {
	Host string `xml:"host"`
	Url  struct {
		Signup   string `xml:"signup"`
		Login    string `xml:"login"`
		Upload   string `xml:"upload"`
		Download string `xml:"download"`
		Delete   string `xml:"delete"`
		Search   string `xml:"search"`
		Albummgt struct {
			Insert string `xml:"insert"`
			Delete string `xml:"delete"`
			Get    string `xml:"get"`
		} `xml:"albummgt"`
		Deeplearning       string `xml:"deeplearning"`
		Objectdetection_dl string `xml:"objectdetection_dl"`
	} `xml:"url"`
}

//configure
type configure struct {
	XMLName   xml.Name  `xml:"configure"`
	GinServer GinServer `xml:"gin_server"`
	Backend   Backend   `xml:"backend"`
}

var (
	conf      *configure
	conf_once sync.Once
)

//Configure 载入xml配置文件
func Configure(file string) *configure {
	conf_once.Do(func() {
		conf = &configure{}
		if err := conf.init(file); err != nil {
			log.Fatalln(err)
		}
	})
	return conf
}

//init 载入xml配置文件
func (c *configure) init(file string) error {
	fd, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("error open file %s fail,%v", file, err)
	}
	defer fd.Close()
	/*
		//处理不了字段中包含html字符的,比如&
		content, err := ioutil.ReadAll(fd)
		if err != nil {
			return err
		}
		return xml.Unmarshal(content, c)
	*/
	//使用decoder处理包含html字符的内容
	d := xml.NewDecoder(fd)
	d.Strict = false
	d.AutoClose = xml.HTMLAutoClose
	d.Entity = xml.HTMLEntity
	return d.Decode(c)
}

func (c *configure) String() string {
	return fmt.Sprintf("%#v", c)
}

//得到配置实例
func GetConfigure() *configure {
	return conf
}
