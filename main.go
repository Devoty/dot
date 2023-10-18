package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {

	filePath := "config.yml"

	var err error

	err = InitConnection(filePath)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected!")
	http.HandleFunc("/dot", DotServer)
	err = http.ListenAndServe(":8091", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}

func DotServer(_ http.ResponseWriter, req *http.Request) {

	query := req.URL.RawQuery

	querys := strings.Split(query, "&")

	var dotTag string

	for _, qu := range querys {
		params := strings.Split(qu, "=")
		if params[0] == "tag" {
			dotTag = params[1]
		}
	}

	_, err := recordDot(dotTag)
	if err != nil {
		return
	}
}

func recordDot(dotTag string) (int64, error) {

	stmt, err := Db.Prepare("insert into dot (tag_code, dot_time) VALUES( ?, ?)")

	if _, err := stmt.Exec(dotTag, time.Now().Format("2006-01-02 15:04:05")); err != nil {
		log.Fatal(err)
	}

	return 0, err
}

var Db *sql.DB

type DataSource struct {
	Mysql struct {
		User   string `yaml:"user"`
		Passwd string `yaml:"password"`
		Host   string `yaml:"host"`
		Port   int    `yaml:"port"`
		Name   string `yaml:"dbname"`
	}
}

func ReadYaml(initFile string) (*os.File, error) {
	dbConfig, err := os.Open(initFile)

	if err != nil {
		return nil, err
	}

	return dbConfig, nil
}

func InitConnection(initFile string) error {
	file, err := ReadYaml(initFile)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("dec.Decode() failed with `%s` \n", err)
			return
		}
	}(file)

	decoder := yaml.NewDecoder(file)

	var ds DataSource

	err = decoder.Decode(&ds)
	if err != nil {
		log.Fatalf("dec.Decode() failed with `%s` \n", err)
		return err
	}

	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True", ds.Mysql.User, ds.Mysql.Passwd, ds.Mysql.Host, ds.Mysql.Port, ds.Mysql.Name)

	Db, err = sql.Open("mysql", url)

	if err != nil {
		log.Fatalf("配置文件出错，原因为 %s\n", err)
		return err
	}
	err = Db.Ping()
	if err != nil {
		log.Fatalf("数据库连接失败，错误为%s\n", err)
		return err
	}
	return nil

}
