/*****************************************************************************
 *  file name : Field.go
 *  author : Wu Yinghao
 *  email  : wyh817@gmail.com
 *
 *  file description : 字段的基本单元，一个field中可以包含一个倒排一个正排
 *
******************************************************************************/

package FalconIndex

import (
	"utils"
)

// Field struct description : 字段的基本单元，需要实现indexinterface的所有方法，这是对外的最基本单元
type Field struct {
	DimensionName string               `json:"dimensionname"` //维度名称
	IndexName     string               `json:"indexname"`     //索引名称
	FieldName     string               `json:"fieldname"`     //字段名称
	FullName      string               `json:"fullname"`      //全名，用来保存数据用
	SerialNum     uint32               `json:"serialname"`    //序列号
	StartDocId    uint32               `json:"startdocid"`    //起始docid
	DocIdLens     uint32               `json:"docidlens"`     //docid长度
	MaxDocId      uint32               `json:"maxdocid"`      //最大的docid
	FieldType     uint64               `json:"fieldtype"`     //字段类型
	IsMomery      bool                 `json:"ismomery"`      //是否还在内存中没有序列化到硬盘中
	Logger        *utils.Log4FE        `json:"-"`             //logger
	ivtert        utils.InvertInterface //一个倒排接口
	profile       utils.ProfileInterface //一个正排接口
}


// AddDocument function description : 增加一个doc文档
// params : docid docid的编号
//			contentstr string  文档内容
// return : error 成功返回Nil，否则返回相应的错误信息
func (this *Field) AddDocument(docid uint64, contentstr string) error {
	return nil
}

// Query function description : 给定一个查询词query，找出doc的列表（标准操作）
// params : key string 查询的key值
// return : docid结构体列表  bool 是否找到相应结果
func (this *Field) Query(key interface{}) ([]utils.DocIdNode, bool) {
	return nil, false
}

// Serialization function description : 序列化倒排索引（标准操作）
// params :
// return : error 正确返回Nil，否则返回错误类型
func (this *Field) Serialization() error {
	return nil
}

// Load function description : 序列化
// params :
// return : error 正确返回Nil，否则返回错误类型
func (this *Field) Load() error {
	return nil
}



// Destroy function description : 销毁字段
// params :
// return :
func (this *Field) Destroy() error {
	return nil
}


// Filter function description : 过滤
// params :
// return :
func (this *Field) Filter(docid uint64, filtertype uint64, start, end int64) bool {
	return false
}

// GetValue function description : 获取字符串类型的值
// params :
// return :
func (this *Field) GetStringValue(docid uint64) (string, bool) {
	return "", false
}

// GetIntValue function description : 获取数字类型的值
// params :
// return :
func (this *Field) GetNumberValue(docid uint64) (int64, bool) {
	return 0, false
}

