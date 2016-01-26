package utils

type DocIdNode uint64

// 索引类型说明
const (
	IDX_TYPE_STRING      = 1 //字符型索引[全词匹配]
	IDX_TYPE_STRING_SEG  = 2 //字符型索引[切词匹配，全文索引]
	IDX_TYPE_STRING_LIST = 3 //字符型索引[列表类型，直接切分]

	IDX_TYPE_NUMBER = 11 //数字型索引，只支持整数，数字型索引只建立正排
	IDX_TYPE_DATE   = 12 //日期型索引 '2015-11-11 00:11:12'，日期型只建立正排

	IDX_TYPE_PK = 21 //主键类型，倒排正排都需要
	//GATHER_TYPE = 22 //汇总类型，倒排正排都需要[后续使用]

)

// 过滤类型，对应filtertype
const (
	FILT_EQ    uint64 = 1 //等于
	FILT_OVER  uint64 = 2 //大于
	FILT_LESS  uint64 = 3 //小于
	FILT_RANGE uint64 = 4 //范围内
)

/*************************************************************************
索引类的接口
************************************************************************/
// IndexInterface interface description : 索引的接口描述，所有索引（正排，倒排）都需要实现此接口
type IndexInterface interface {
	AddDocument(docid uint64, contentstr string) error             //添加文档
	Query(key interface{}) ([]DocIdNode, bool)                     //关键词查询[倒排]
	Serialization() error                                          //序列化数据
	Filter(docid uint64, filtertype uint64, start, end int64) bool //过滤操作[正排]【>,<,==,!=,<>》】
	GetValue(docid uint64) (string, bool)                          //对于字符型正排，获取值
	GetIntValue(docid uint64) (int64, bool)                        //对于数字型正排，获取值
	Destroy() error                                                //销毁这个字段的内容
}
