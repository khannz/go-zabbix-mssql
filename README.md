[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/khannz/go-zabbix-mssql/master/LICENSE) [![Build Status](https://travis-ci.org/khannz/go-zabbix-mssql.svg?branch=master)](https://travis-ci.org/khannz/go-zabbix-mssql) [![Go Report Card](https://goreportcard.com/badge/github.com/khannz/go-zabbix-mssql)](https://goreportcard.com/report/github.com/khannz/go-zabbix-mssql)

# go-zabbix-mssql

## Intro

First of all, this is the utility to collect tricky metrics from `MS SQL Server` instance and send it to `Zabbix trapper` items using [Go implementation of zabbix-sender](https://github.com/AlekSi/zabbix-sender). By 'tricky' I mean things that must be collected from inside of Instance with proper queries.

> Note: you can use go-zabbix-mssql to grab any data from DBs of Instance with key:value query result. Just check out queries from YAML, to get get exact idea.

In case there would be developed any universal templates that can be easily imported into Zabbix and/or Grafana, it would be also added to repository. Anyway, I got plans to describe everything later in some kind of guide.

One of the goals is to achieve something like following view with `zabbix` + `windows performance counters` + `go-zabbix-mssql` + `grafana`
![telegraf-sqlserver-full](https://cloud.githubusercontent.com/assets/16494280/12591426/fa2b17b4-c467-11e5-9c00-929f4c4aea57.png)

## How to cook

With just one config YAML file, go-zabbix-mssql allows you to send (almost?) endless list of queries that is joining with `UNION ALL` before SQL engine get it. So technically any number of lines in YAML still produces just one query, amirite? Be careful about response time for final big query!

### Running

You **must** use command line flags to tune few details for zabbix-sender part:
* **-U** -- login for MS SQL connection
* **-P** -- password for MS SQL connection
* **-S** -- connection string in `server_name[\instance_name]` format
* **-Z** -- zabbix-server FQDN where proper trapper items live
* **-H** -- MS SQL Instance hostname, as it named in zabbix-server
* **-F** -- full path to config YAML file, like `go-zabbix-mssql.config.yaml`

### Setting up Zabbix Agent

go-zabbix-mssql can (and intended to) be easily used as part of zabbix-agent:

* add compiled binary to your `${zabbix-agent\scripts}` or whatever path you use to store extensions in your zabbix-agent distributive
* make sure, you have `UnsafeUserParameters=1` in `zabbix_agentd.conf`
* also add something like `UserParameter=mssql.metrics[*],${zabbix-agent\scripts}\go-zabbix-mssql.exe -U=$1 -P=$2 -S=$3 -Z=$4 -H=$5 -F="${full path to YAML config}"`

### Setting up Zabbix Server

* configure proper Macros details (in Zabbix console) to every monitored host (read official guide for `UnsafeUserParameters`) which gonna replace all those `$1`, `$2` and so on... In my case, part of Macroses belongs to 'MS SQL 2012 Template' template and part of it defined manually for each Instance.
* add items to your template/host with type `Zabbix trapper`. Keys must be equal to metric name you await to receive
* add and use item of `Zabbix agent` type with key `mssql.metrics` on each MS SQL Instance, configured as `Host` or in general template for MS SQL Instances 
     * that's how you will control frequency of metrics gathering
     * this configuration also sends values of Macroses for every host, so you get most control of how your monitoring configured on zabbix-server side 
     * execution of binary returns decimal number -- it shows how much of all collected metrics was not received correctly by zabbix-server. So perfect return must be `0` and that what would be self metric for go-zabbix-mssql (in case of troubles, answer will become text of error)

### Getting metrics for exact DBs 

TL;DR: You can notice that some part of metrics (last one monstrous query for now) from examples got also dynamic parts that returns metrics for each DB, not for whole Instance. If it is a goal for you (it is for me), it can be used with help of [popular method of MS SQL monitoring](https://share.zabbix.com/databases/microsoft-sql-server/template-app-ms-sql-default-installation-lld) by [Stephen Fritz](https://share.zabbix.com/owner/g_111769865974589121086).

This solution collects data by LLD via WMI with PowerShell script. You receive an array `{#DBS}` with names of DBs as result. **Important to notice** - attached setup guide suggests to cut off system DBs from that LLD with *regular expression filter*. In case you planning to use go-zabbix-mssql for metrics of every table, do not implement this step.    

## Notes

I love [Prometheus](https://github.com/prometheus/prometheus) project and its ideas so that's why metrics from example YAML named so.

## TODOs

All ideas respectively placed in [Projects](https://github.com/khannz/go-zabbix-mssql/projects) and [Milestones](https://github.com/khannz/go-zabbix-mssql/milestones)

## Thanks

- Whole [Golang Project team](https://golang.org/project/), involved developers & engineers of [Google](https://google.com) as well as [Project Contributors](https://golang.org/CONTRIBUTORS)
- [JetBrains s.r.o.](https://www.jetbrains.com) for [Gogland IDE](https://www.jetbrains.com/go/) which is damn awesome and feature-rich (for my very tiny needs) right from start of EAP program
- Some guys from [Telegraf](https://github.com/influxdata/telegraf) contributors for R&D of [detailed and sophisticated metrics collection queries](https://github.com/influxdata/telegraf/blob/master/plugins/inputs/sqlserver/sqlserver.go#L191), which was perfectly combined for Grafana dashboard (shown at illustration above)
- [AlekSi](https://github.com/AlekSi) for [Go implementation of zabbix-sender](https://github.com/AlekSi/zabbix-sender)
- [denisenkom](https://github.com/denisenkom) for [Microsoft SQL driver written in Go](https://github.com/denisenkom/go-mssqldb)
- Team of contributors, who developed [YAML.v2](https://github.com/go-yaml/yaml/tree/v2) Go package
- [@mmcgrana](https://twitter.com/mmcgrana) for [Go by Example](https://gobyexample.com) project which helped me to dive into Golang with minimal doubts

Hope it all is not too pretentious for 123 lines of code :trollface:
