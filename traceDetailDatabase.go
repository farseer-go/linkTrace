package linkTrace

import (
	"fmt"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/trace"
	"strings"
)

type TraceDetailDatabase struct {
	trace.BaseTraceDetail
	DbName           string // 数据库名
	TableName        string // 表名
	Sql              string // SQL
	ConnectionString string // 连接字符串
	RowsAffected     int64  // 影响行数
}

func (receiver *TraceDetailDatabase) GetTraceDetail() *trace.BaseTraceDetail {
	return &receiver.BaseTraceDetail
}

func (receiver *TraceDetailDatabase) ToString() string {
	if receiver.Sql != "" {
		sql := receiver.Sql
		if len(sql) > 1000 {
			sql = sql[:1000] + "......"
		}
		sql = flog.ReplaceBlues(sql, "SELECT ", "UPDATE ", "DELETE ", "INSERT INTO ", " FROM ", " WHERE ", " LIMIT ", " SET ", " ORDER BY ", " VALUES ", " and ", " or ", "`")
		sql = strings.ReplaceAll(sql, receiver.TableName, flog.Green(receiver.TableName))
		return fmt.Sprintf("[%s]耗时：%s，[影响%d行]%s", flog.Yellow(receiver.CallType.ToString()), flog.Red(receiver.UseTs.String()), receiver.RowsAffected, sql)
	} else if receiver.ConnectionString != "" {
		return fmt.Sprintf("[%s]耗时：%s， 连接数据库：%s", flog.Yellow(receiver.CallType.ToString()), flog.Red(receiver.UseTs.String()), receiver.ConnectionString)
	}
	return ""
}

func (receiver *TraceDetailDatabase) SetSql(connectionString string, DbName string, tableName string, sql string, rowsAffected int64) {
	receiver.ConnectionString = connectionString
	receiver.DbName = DbName
	receiver.TableName = tableName
	receiver.Sql = sql
	receiver.RowsAffected = rowsAffected
}

func (receiver *TraceDetailDatabase) Desc() (caption string, desc string) {
	if receiver.TableName == "" && receiver.Sql == "" {
		caption = fmt.Sprintf("打开数据库 => %s %s", receiver.DbName, receiver.ConnectionString)
	} else {
		caption = fmt.Sprintf("执行数据库 => %s %s 影响%v行", receiver.DbName, receiver.TableName, receiver.RowsAffected)
	}
	desc = receiver.Sql
	return
}
