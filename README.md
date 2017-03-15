# go-zabbix-mssql

## Intro

First of all, this is the utility to collect tricky metrics from `MS SQL Server` instance and send it to Zabbix trapper items using [Go implementation of zabbix-sender](https://github.com/AlekSi/zabbix-sender).

In case there would be developed any universal templates that can be easily imported into Zabbix or Grafana, it would be also added to repository. Anyway, I got plans to describe everything in some kind of guide.

This development inspired by [Grafana dashboard](https://grafana.net/dashboards/409) which utilizes list of [MS SQL metrics by Telegraf](https://github.com/influxdata/telegraf/tree/master/plugins/inputs/sqlserver) team.

One of the goals is to achieve something like following view with `zabbix` + `windows performance counters` + `go-zabbix-mssql` + `grafana`
![telegraf-sqlserver-full](https://cloud.githubusercontent.com/assets/16494280/12591426/fa2b17b4-c467-11e5-9c00-929f4c4aea57.png)

## Features

1. With just one config YAML file, this util allows to send almost endless list of queries that is joining with `UNION ALL` before SQL engine get it. So technically any number of lines in YAML still produces just one query, amirite? Be careful about response time for final big query!
1. You **must** use command line flags to tune few details for zabbix-sender part^
    * **-U** -- login for MS SQL connection
    * **-P** -- password for MS SQL connection
    * **-S** -- connection string in `server_name[\instance_name]` format
    * **-Z** -- zabbix-server FQDN where proper trapper items live
    * **-H** -- MS SQL Instance hostname, as it named in zabbix-server
    * **-F** -- full path to config YAML file, like `go-zabbix-mssql.config.yaml`
1. Util can be easily used as part of zabbix-agent:
    * make sure, you have `UnsafeUserParameters=1` in `zabbix_agentd.conf`
    * also add something like `UserParameter=mssql.metrics[*],${zabbix-agent\scripts}\go-zabbix-mssql.exe -U=$1 -P=$2 -S=$3 -Z=$4 -H=$5 -F="${full path to YAML config}"`
    * configure proper Macros details (in Zabbix console) to every monitored host (read official guide for `UnsafeUserParameters`) which gonna replace all those `$1`, `$2` and so on...
    * add items to your template/host with type `Zabbix trapper`
1. Execution of code returns decimal number -- it shows how much of all collected metrics was not received correctly by zabbix-server. So perfect return must be `0`.

## Notes

1. I love [Prometheus](https://github.com/prometheus/prometheus) project and its ideas so that's why metrics from example YAML named so.

## TODOs
1. Need to finish list of metrics for YAML
1. Add all metrics to zabbix-server as proper `Zabbix trapper` items
1. Finish Grafana dashboard tuning
1. Add necessary exports to current repo
1. Describe everything in documentation