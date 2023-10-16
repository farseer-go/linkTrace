package linkTrace

import (
	"fmt"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/linkTrace/eumCallType"
	"strings"
	"time"
)

type TraceDatabaseDetail struct {
	TraceDetail
	DbName    string // 数据库名
	TableName string // 表名
	Sql       string // SQL
}

func (receiver *TraceDatabaseDetail) ToString(index int) string {
	sql := flog.ReplaceBlues(receiver.Sql, "SELECT ", "UPDATE ", "DELETE ", " FROM ", " WHERE ", " LIMIT ", " SET ", " ORDER BY ", " and ", " or ")
	sql = strings.ReplaceAll(sql, receiver.TableName, flog.Red(receiver.TableName))
	return fmt.Sprintf("%s：[%s] %s", flog.Blue(index), flog.Yellow(receiver.CallType.ToString()), sql)
}

func TraceDatabase(dbName, tableName, sql string) *TraceDatabaseDetail {
	detail := &TraceDatabaseDetail{
		TraceDetail: TraceDetail{
			CallStackTrace: CallStackTrace{},
			CallMethod:     "",
			CallType:       eumCallType.Database,
			StartTs:        time.Now().UnixMicro(),
		},
		DbName:    dbName,
		TableName: tableName,
		Sql:       sql,
	}

	if trace := GetCurTrace(); trace != nil && defConfig.Enable {
		trace.List.Add(detail)
	}
	return detail
}
