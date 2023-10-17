package linkTrace

import (
	"fmt"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/linkTrace/eumCallType"
	"strings"
)

type TraceDatabaseDetail struct {
	TraceDetail
	DbName           string // 数据库名
	TableName        string // 表名
	Sql              string // SQL
	ConnectionString string // 连接字符串
}

func (receiver *TraceDatabaseDetail) GetTraceDetail() *TraceDetail {
	return &receiver.TraceDetail
}

func (receiver *TraceDatabaseDetail) ToString() string {
	if receiver.Sql != "" {
		sql := flog.ReplaceBlues(receiver.Sql, "SELECT ", "UPDATE ", "DELETE ", " FROM ", " WHERE ", " LIMIT ", " SET ", " ORDER BY ", " and ", " or ")
		sql = strings.ReplaceAll(sql, receiver.TableName, flog.Green(receiver.TableName))
		return fmt.Sprintf("[%s]耗时：%s， %s", flog.Yellow(receiver.CallType.ToString()), flog.Red(receiver.UseTs.String()), sql)
	} else if receiver.ConnectionString != "" {
		return fmt.Sprintf("[%s]耗时：%s， 连接数据库：%s", flog.Yellow(receiver.CallType.ToString()), flog.Red(receiver.UseTs.String()), receiver.ConnectionString)
	}
	return ""
}

// TraceDatabase 数据库埋点
func TraceDatabase() *TraceDatabaseDetail {
	detail := &TraceDatabaseDetail{
		TraceDetail: newTraceDetail(eumCallType.Database),
	}
	add(detail)
	return detail
}
