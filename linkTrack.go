package linkTrack

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type linkData struct {
	moduleName string    // 模块名称
	startTime  time.Time // 执行时间
	endTime    time.Time // 执行结束时间
	useTime    int64
	link       string // 链路路径
	remoteAddr string
	userAgent  string
	requestId  string
}

type linkTrack struct {
	chlQueue chan linkData // 通道队列，打开通道后，不关系，放回此队列
	config   linkTrackConfig
}

// Complete 添加记录
func (linkData *linkData) Complete() {
	if !linkData.startTime.IsZero() {
		linkData.endTime = time.Now()
		linkData.useTime = linkData.endTime.Unix() - linkData.startTime.Unix()
		// @todo 根据配置插入es或者数据库等
	}
}

// Create一个请求
func Create(moduleNmae string, link string, r *http.Request) *linkData {
	requestId := strconv.Itoa(int(time.Now().Unix())) + strconv.Itoa(rand.Intn(1000000))
	r.Header.Set("rq_id", requestId)
	return &linkData{
		moduleName: moduleNmae,
		startTime:  time.Now(),
		link:       link,
		remoteAddr: getIP(r),
		userAgent:  r.Header.Get("User-Agent"),
		requestId:  requestId,
	}
}

// New新建一个链路追踪
func New(r *http.Request, moduleNmae string, link string) *linkData {
	return &linkData{
		moduleName: moduleNmae,
		startTime:  time.Now(),
		link:       link,
		requestId:  r.Header.Get("rq_id"),
	}
}

func getIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if net.ParseIP(ip) != nil {
		return ip
	}
	ip = r.Header.Get("X-Forward-For")
	for _, i := range strings.Split(ip, ",") {
		if net.ParseIP(i) != nil {
			return i
		}
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}
	if net.ParseIP(ip) != nil {
		return ip
	}
	return ""
}

// Init 初始化
func (receiver *linkTrack) Init() {
	// 首次使用
	if receiver.chlQueue == nil {
		receiver.chlQueue = make(chan linkData, receiver.config.ChannelSize)
	}
	go func() {
		var datas []linkData
		for v := range receiver.chlQueue { // 读取所有，直到通道关闭
			datas = append(datas, v)
		}
		s, _ := strconv.Atoi(fmt.Sprintf("%d", time.Second))
		time.Sleep(time.Duration(receiver.config.TimeInterval * s))
	}()
}
