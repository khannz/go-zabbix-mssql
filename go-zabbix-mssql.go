package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/AlekSi/zabbix-sender"
	_ "github.com/denisenkom/go-mssqldb"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net"
	"path/filepath"
	"regexp"
)

func main() {
	var (
		userid   = flag.String("U", "", "login for MS SQL connection")
		password = flag.String("P", "", "password for MS SQL connection")
		server   = flag.String("S", "", "server_name[\\instance_name]")
		zServer  = flag.String("Z", "", "zabbix-server name with proper trapper items configured")
		zHost    = flag.String("H", "", "hostname as it presented in Zabbix")
		confFile = flag.String("F", "", "full path to config YAML file")
	)
	flag.Parse()

	dsn := "server=" + *server + ";user id=" + *userid + ";password=" + *password

	db, err := sql.Open("mssql", dsn)
	if err != nil {
		fmt.Println("Cannot connect: ", err.Error())
		return
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Cannot connect: ", err.Error())
		return
	}

	defer db.Close()

	// Reading config from YAML
	type Config struct {
		Metrics []string `yaml:"sql_metrics"`
	}

	filename, _ := filepath.Abs(*confFile)
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var config Config

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	concatQuery := ``
	longQueriesResult := make(map[string]interface{})
	for i := 0; i < len(config.Metrics); i++ {
		if isML(config.Metrics[i]) == false {
			if concatQuery != `` {
				concatQuery += "\nUNION ALL\n"
			}
			concatQuery += config.Metrics[i]
		} else {
			for k, v := range dbQuery(db, config.Metrics[i]) {
				longQueriesResult[k] = v
			}
		}
	}

	resultQuery := dbQuery(db, concatQuery)
	for k, v := range longQueriesResult {
		resultQuery[k] = v
	}

	data := resultQuery
	di := zabbix_sender.MakeDataItems(data, *zHost)
	addr, _ := net.ResolveTCPAddr("tcp", *zServer+":10051")
	res, _ := zabbix_sender.Send(addr, di)
	fmt.Println(res.Failed) //
}

func dbQuery(db *sql.DB, cmd string) map[string]interface{} {
	var (
		name string
		val  float64
	)

	rows, err := db.Query(cmd)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	pipisa := make(map[string]interface{})

	for rows.Next() {
		err = rows.Scan(&name, &val)
		if err != nil {
			log.Fatal(err)
		}
		pipisa[name] = val
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return pipisa
}

func isML(val string) bool { // returns TRUE if *val (supposed to be query string from config) contains multiple lines
	match, _ := regexp.MatchString("[dD][eE][cC][lL][aA][rR][eE]\\s", val) // used regex pattern for DECLARE
	return match
}
