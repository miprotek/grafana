+++
title = "What's New in Grafana v5.1"
description = "Feature & improvement highlights for Grafana v5.1"
keywords = ["grafana", "new", "documentation", "5.1"]
type = "docs"
[menu.docs]
name = "Version 5.1"
identifier = "v5.1"
parent = "whatsnew"
weight = -7
+++

# What's New in Grafana v5.1

Grafana v5.1 brings many enhancements to Prometheus, Cloudwatch, PostgreSQL, MySQL, alerting, graph and table panels. It also adds support for Microsft SQL Server as metric & table data source!

### Prometheus


### Cloudwatch


### PostgreSQL
New enhancement includes support for filling in values for missing intervals and better precision handling and data type support for time columns. More inforomation in the [PostgreSQL data source docs](/features/datasources/postgres/#time-series-queries).

### MySQL
New enhancements includes support for filling in values for missing intervals, better precision handling and data type support for time columns and any column except time and metric is now treated as a value column. More inforomation in the [MySQL data source docs](/features/datasources/mysql/#time-series-queries).

#### Microsoft SQL Server

Grafana v5.1 now ships with a built-in Microsoft SQL Server (MSSQL) data source plugin that allows you to query and visualize data from any Microsoft SQL Server 2005 or newer, including Microsoft Azure SQL Database. Have logs or metric data in MSSQL? You can now visualize that data and
define alert rules on it like any of our other data sources.

Same enhancements as for PostgreSQL and MySQL have been incorporated for MSSQL as well. More inforomation in the [Microsoft SQL Server data source docs](/features/datasources/mssql/).

{{< docs-imagebox img="/img/docs/v51/mssql_query_editor.png"  max-width= "800px" >}}

### Provisioning


### Alerting



### Graph Panel


### Table Panel


## Changelog

### New Features

* **MSSQL**: New Microsoft SQL Server data source [#10093](https://github.com/grafana/grafana/pull/10093), [#11298](https://github.com/grafana/grafana/pull/11298), thx [@linuxchips](https://github.com/linuxchips)
* **Prometheus**: The heatmap panel now support Prometheus histograms [#10009](https://github.com/grafana/grafana/issues/10009)
* **Postgres/MySQL**: Ability to insert 0s or nulls for missing intervals [#9487](https://github.com/grafana/grafana/issues/9487), thanks [@svenklemm](https://github.com/svenklemm)
* **Postgres/MySQL/MSSQL**: Fix precision for the time column in table mode [#11306](https://github.com/grafana/grafana/issues/11306)
* **Graph**: Align left and right Y-axes to one level [#1271](https://github.com/grafana/grafana/issues/1271) &  [#2740](https://github.com/grafana/grafana/issues/2740) thx [@ilgizar](https://github.com/ilgizar)
* **Graph**: Thresholds for Right Y axis [#7107](https://github.com/grafana/grafana/issues/7107), thx [@ilgizar](https://github.com/ilgizar)
* **Graph**: Support multiple series stacking in histogram mode [#8151](https://github.com/grafana/grafana/issues/8151), thx [@mtanda](https://github.com/mtanda)
* **Alerting**: Pausing/un alerts now updates new_state_date [#10942](https://github.com/grafana/grafana/pull/10942)
* **Alerting**: Support Pagerduty notification channel using Pagerduty V2 API [#10531](https://github.com/grafana/grafana/issues/10531), thx [@jbaublitz](https://github.com/jbaublitz)
* **Templating**: Add comma templating format [#10632](https://github.com/grafana/grafana/issues/10632), thx [@mtanda](https://github.com/mtanda)
* **Prometheus**: Show template variable candidate in query editor [#9210](https://github.com/grafana/grafana/issues/9210), thx [@mtanda](https://github.com/mtanda)
* **Prometheus**: Support POST for query and query_range [#9859](https://github.com/grafana/grafana/pull/9859), thx [@mtanda](https://github.com/mtanda)
* **Alerting**: Add support for retries on alert queries [#5855](https://github.com/grafana/grafana/issues/5855), thx [@Thib17](https://github.com/Thib17)
* **Table**: Table plugin value mappings [#7119](https://github.com/grafana/grafana/issues/7119), thx [infernix](https://github.com/infernix)
* **IE11**: IE 11 compatibility [#11165](https://github.com/grafana/grafana/issues/11165)

### Minor Changes

* **OpsGenie**: Add triggered alerts as description [#11046](https://github.com/grafana/grafana/pull/11046), thx [@llamashoes](https://github.com/llamashoes)
* **Cloudwatch**: Support high resolution metrics [#10925](https://github.com/grafana/grafana/pull/10925), thx [@mtanda](https://github.com/mtanda)
* **Cloudwatch**: Add dimension filtering to CloudWatch `dimension_values()` [#10029](https://github.com/grafana/grafana/issues/10029), thx [@willyhutw](https://github.com/willyhutw)
* **Units**: Second to HH:mm:ss formatter [#11107](https://github.com/grafana/grafana/issues/11107), thx [@gladdiologist](https://github.com/gladdiologist)
* **Singlestat**: Add color to prefix and postfix in singlestat panel [#11143](https://github.com/grafana/grafana/pull/11143), thx [@ApsOps](https://github.com/ApsOps)
* **Dashboards**: Version cleanup fails on old databases with many entries [#11278](https://github.com/grafana/grafana/issues/11278)
* **Server**: Adjust permissions of unix socket [#11343](https://github.com/grafana/grafana/pull/11343), thx [@corny](https://github.com/corny)
* **Shortcuts**: Add shortcut for duplicate panel [#11102](https://github.com/grafana/grafana/issues/11102)
* **AuthProxy**: Support IPv6 in Auth proxy white list [#11330](https://github.com/grafana/grafana/pull/11330), thx [@corny](https://github.com/corny)
* **SMTP**: Don't connect to STMP server using TLS unless configured. [#7189](https://github.com/grafana/grafana/issues/7189)
* **Prometheus**: Escape backslash in labels correctly. [#10555](https://github.com/grafana/grafana/issues/10555), thx [@roidelapluie](https://github.com/roidelapluie)
* **Variables** Case-insensitive sorting for template values [#11128](https://github.com/grafana/grafana/issues/11128) thx [@cross](https://github.com/cross)
