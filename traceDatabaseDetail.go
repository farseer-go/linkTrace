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
	sql = strings.ReplaceAll(sql, receiver.TableName, flog.Green(receiver.TableName))
	return fmt.Sprintf("%s：[%s]耗时：%s， %s", flog.Blue(index), flog.Yellow(receiver.CallType.ToString()), flog.Red(receiver.UseTs.String()), sql)
}

func TraceDatabase() *TraceDatabaseDetail {
	detail := &TraceDatabaseDetail{
		TraceDetail: TraceDetail{
			//CallStackTrace: CallStackTrace{},
			CallMethod: "",
			CallType:   eumCallType.Database,
			StartTs:    time.Now().UnixMicro(),
		},
	}

	if trace := GetCurTrace(); trace != nil && defConfig.Enable {
		trace.List.Add(detail)
	}
	return detail
}
