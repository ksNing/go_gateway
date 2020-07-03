package dao

import (
	"errors"
	"fmt"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/go_gateway/dto"
	"github.com/go_gateway/public"
	"net/http/httptest"
	"strings"
	"sync"
)

type ServiceDetail struct {
	Info          *ServiceInfo   `json:"info" description:"基本信息"`
	HTTPRule      *HttpRule      `json:"http_rule" description:"http_rule"`
	TCPRule       *TcpRule       `json:"tcp_rule" description:"tcp_rule"`
	GRPCRule      *GrpcRule      `json:"grpc_rule" description:"grpc_rule"`
	LoadBalance   *LoadBalance   `json:"load_balance" description:"load_balance"`
	AccessControl *AccessControl `json:"access_control" description:"access_control"`
}

var ServiceManagerHandler *ServiceManager

func init() {
	ServiceManagerHandler = NewServiceManager()
}

type ServiceManager struct {
	ServiceMap map[string] *ServiceDetail
	ServiceSlice []*ServiceDetail
	Locker sync.Mutex
	init sync.Once
	err error
}

func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		ServiceMap:map[string] *ServiceDetail{},
		ServiceSlice:[] *ServiceDetail{},
		Locker: sync.Mutex{},
		init: sync.Once{},
	}
}
//将所有服务加载到内存中
func (s *ServiceManager) LoadOnce() error {
	s.init.Do(func() {
		tx, err := lib.GetGormPool("default")
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		if err != nil {
			s.err = err
			return
		}
		params := &dto.ServiceListInput{PageSize: 99999, PageNo:1}

		//从db中分页读取服务信息
		serviceInfo := &ServiceInfo{}
		list, _, err := serviceInfo.PageList(c,tx,params)
		if err != nil {
			s.err = err
			return
		}
		//加锁
		s.Locker.Lock()
		defer s.Locker.Unlock()
		for _,listItem := range list {
			tempList := listItem
			serviceDetail, err := tempList.ServiceDetail(c, tx, &tempList)
			if err != nil {
				s.err = err
				return
			}
			s.ServiceMap[listItem.ServiceName] = serviceDetail
			s.ServiceSlice = append(s.ServiceSlice, serviceDetail)
		}
		fmt.Printf("s",public.Obj2Json(s))

	})
	return s.err

}

//服务接入层判断
func (s *ServiceManager) HTTPAccessMode(c *gin.Context) (*ServiceDetail, error) {
	//1 前缀匹配 /abc => serviceSlice.rule
	//2 域名匹配 www.test.com ==> serviceSlice.rule

	//domain c.Request.host
	//path c.Request.URL.path

	host := c.Request.Host
	host = host[0:strings.Index(host,":")]
	fmt.Printf("host",host)
	path := c.Request.URL.Path
	fmt.Printf("path",path)

	for _, serviceItem := range s.ServiceSlice {
		if serviceItem.Info.LoadType != public.LoadTypeHTTP {
			continue
		}

	    if serviceItem.HTTPRule.RuleType == public.HTTPRuleTypeDomain {
			if serviceItem.HTTPRule.Rule != "" && serviceItem.HTTPRule.Rule == host {
				return serviceItem, nil
			}
		}

		if serviceItem.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL {
			if strings.HasPrefix(path, serviceItem.HTTPRule.Rule) {
				return serviceItem, nil
			}
 		}

	}
	return nil, errors.New("no match services")


}









