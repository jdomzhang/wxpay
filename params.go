package wxpay

import (
	"fmt"
	"strconv"
)

type Params map[string]string

// map本来已经是引用类型了，所以不需要 *Params
func (p Params) SetString(k, s string) Params {
	p[k] = s
	return p
}

func (p Params) GetString(k string) string {
	s, _ := p[k]
	return s
}

func (p Params) SetInt64(k string, i int64) Params {
	p[k] = strconv.FormatInt(i, 10)
	return p
}

func (p Params) GetInt64(k string) int64 {
	i, _ := strconv.ParseInt(p.GetString(k), 10, 64)
	return i
}

// 判断key是否存在
func (p Params) ContainsKey(key string) bool {
	_, ok := p[key]
	return ok
}

func (p Params) Remove(key string) {
	if p.ContainsKey(key) {
		delete(p, key)
	}
}

func (obj Params) toError() error {

	if obj.GetString("return_code") != SUCCESS {
		return fmt.Errorf(obj.GetString("return_msg"))
	}

	if obj.GetString("result_code") != SUCCESS {
		return fmt.Errorf(obj.GetString("err_code_des"))
	}

	return nil
}
