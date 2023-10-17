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

func (receiver *TraceDatabaseDetail) ToString() string {
	sql := flog.ReplaceBlues(receiver.Sql, "SELECT ", "UPDATE ", "DELETE ", " FROM ", " WHERE ", " LIMIT ", " SET ", " ORDER BY ", " and ", " or ")
	sql = strings.ReplaceAll(sql, receiver.TableName, flog.Green(receiver.TableName))
	return fmt.Sprintf("[%s]耗时：%s， %s", flog.Yellow(receiver.CallType.ToString()), flog.Red(receiver.UseTs.String()), sql)
}

// TraceDatabase 数据库埋点
func TraceDatabase() *TraceDatabaseDetail {
	detail := &TraceDatabaseDetail{
		TraceDetail: newTraceDetail(eumCallType.Database),
	}

	if trace := GetCurTrace(); trace != nil && defConfig.Enable {
		// 时间轴：上下文入口起点时间到本次开始时间
		detail.Timeline = time.Duration(detail.StartTs-trace.StartTs) * time.Microsecond
		if trace.List.Count() > 0 {
			detail.UnTraceTs = time.Duration(detail.StartTs-trace.List.Last().GetEndTs()) * time.Microsecond
		} else {
			detail.UnTraceTs = time.Duration(detail.StartTs-trace.StartTs) * time.Microsecond
		}
		trace.List.Add(detail)
	}
	return detail
}
