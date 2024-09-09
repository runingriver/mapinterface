package conf

import "github.com/runingriver/mapinterface/logx"

var (
	CONF *Conf
)

type Conf struct {
	// 将对象转成string时是否支持使用对象实现的String方法
	CvtStrUseStringMethod bool
	// 参考ToMapType接口方法,将对象转成map[x]y类型时,是否忽略转换失败的情况,如:转成map[int]int类型时,9个转成功,1个转失败时,是返回含9个元素的map还是返回错误;
	SkipCvtFailForToMapType bool
	// 参考ToArrayType接口方法,将对象转成[]int类型时,是否忽略转换失败的情况,如:[]int类型时,9个转成功,1个转失败时,是返回含9个元素的list还是返回错误;
	SkipCvtFailForToArrayType bool
}

func init() {
	CONF = &Conf{
		CvtStrUseStringMethod:     true,
		SkipCvtFailForToMapType:   false,
		SkipCvtFailForToArrayType: false,
	}
}

func (c *Conf) SetLogger(log logx.LogItf) *Conf {
	logx.SetLogger(log)
	return c
}

// SetCvtStrUseStringMethod 将结构体转换成str的时候是否尝试调用String方法,默认:true
func (c *Conf) SetCvtStrUseStringMethod(b bool) *Conf {
	c.CvtStrUseStringMethod = b
	return c
}

// SetSkipCvtFailForToMapType ToMapType接口转换中,是否跳过部分转换失败的情况,默认:false
func (c *Conf) SetSkipCvtFailForToMapType(b bool) *Conf {
	c.SkipCvtFailForToMapType = b
	return c
}

// SetSkipCvtFailForToArrayType ToArrayType接口转换中,是否跳过部分转换失败的情况,默认:false
func (c *Conf) SetSkipCvtFailForToArrayType(b bool) *Conf {
	c.SkipCvtFailForToArrayType = b
	return c
}
