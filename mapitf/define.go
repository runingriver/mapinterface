package mapitf

import (
	"context"
	"github.com/runingriver/mapinterface/api"
	"github.com/runingriver/mapinterface/conf"
)

/*
类库的入口的接口定义文件, 注:api.mapitf是实现的功能方法的接口定义
*/

type EntranceFuncDefine interface {
	// From 等价于Fr,入口方法
	From(itf interface{}) api.MapInterface
	// Fr 等价于From,携带Context的版本
	Fr(ctx context.Context, itf interface{}) api.MapInterface
}

func Config() *conf.Conf {
	return conf.CONF
}
