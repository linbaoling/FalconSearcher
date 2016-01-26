package HashTable

import (
	"fmt"
)

type Entity struct {
	//Key    [512]byte
	Len    uint64
	Offset uint64
	Index  uint64
	HashA  uint64
	HashB  uint64
}

type DealMmap struct {
	table   []Entity
	initLen int // XXX 初始化table的长度, 一般为 length*1.2
	length  int // 启动服务初始化的实际idx数量
	buckets int
	size    int // 加载到slice中得entity数量 实际和length的值一致，赋值一个在初始化数据前，一个在之后
	mpqhash []MPQTable
}





func NewDealMap(maxsize int) *DealMmap {
	var buckets int

	if maxsize != 0 {
		buckets = custMinBuckets(maxsize)
		for _, size := range BUCKETS_LIST {
			if buckets < size {
				buckets = size
				break
			}
		}
	}
	// 初始化的table的len， 方便后续的更新操作
	initLen := maxsize //+ int(math.Ceil(float64(maxsize) * 0.3))

	dealMap := &DealMmap{
		size:    0,
		initLen: initLen,
		length:  maxsize,
		buckets: buckets,
	}
	// cgo malloc 内存地址
	dealMap.table = mallocEntity(initLen)
	dealMap.mpqhash = NewMPQTable(buckets)

	return dealMap
}

// 查找合适的桶的个数
func custMinBuckets(size int) int {
	var buckets int
	v := size
	// round 到最近的2的倍数
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	// size * 4 /3
	b := size * 4 / 3
	if b > v {
		buckets = b
	} else {
		buckets = v
	}
	return buckets
}

func (this *DealMmap) Display() {
	for _, i := range this.table {
		fmt.Printf("%v \n", i)
	}

	for _, i := range this.mpqhash {
		fmt.Printf("%v \n", i)
	}
}

// 获取初始化的信息
func (this *DealMmap) GetInfo() string {
	return fmt.Sprintf("DealMap has tables: %d, buckets, mpqHashTable: %d", len(this.table), this.buckets)
}

// 返回Entity的数量
func (this *DealMmap) GetSize() int {
	return this.size
}

func (this *DealMmap) SetSize(size int) bool {
	this.size = size
	return true
}

// 返回table的长度
func (this *DealMmap) GetLength() int {
	return this.length
}

// 返回table的cap
func (this *DealMmap) GetCap() int {
	return this.initLen
}

// 返回桶的数量
func (this *DealMmap) GetBuckets() int {
	return this.buckets
}

// 返回桶
func (this *DealMmap) GetMpqTable() []MPQTable {
	return this.mpqhash
}

// 返回索引列
func (this *DealMmap) GetEntityTable() []Entity {
	return this.table
}


func (this *DealMmap)Free() error {
    Free(this.mpqhash)
    Free(this.table)
    return nil
}

// update data
func (this *DealMmap) UpdateTableIndx(entity *Entity, vlen uint64) {
	beforeEntity := this.table[this.length-1]
	offset := beforeEntity.Offset + beforeEntity.Len

	entity.Offset = offset
	entity.Len = vlen

	// 增加冗余数据
	tableIdx := int64(this.length)
	if this.length+1 < this.initLen {
		this.table[tableIdx].Len = vlen
		this.table[tableIdx].Offset = offset
	} else {
		object := Entity{
			Len:    vlen,
			Offset: offset,
		}
		this.table = append(this.table, object)
	}
	this.length += 1
	this.size += 1

	return
}

// add data
func (this *DealMmap) InsertTableIndx(key string, vlen uint64) error {
	beforeEntity := this.table[this.length-1]
	offset := beforeEntity.Offset + beforeEntity.Len

	hashKey := HashKey(key, HASH_O)
	hasha := HashKey(key, HASH_A)
	hashb := HashKey(key, HASH_B)

	start := hashKey % uint64(this.buckets)
	pos := start
	for this.mpqhash[pos].GetPos() != DefaultPos {
		pos = (pos + 1) % uint64(this.buckets)
		if pos == start {
			return fmt.Errorf("the pos had entity. need to add buckets.")
		}
	}
	tableIdx := int64(this.length)

	if this.length+1 < this.initLen {
		// table step2
		this.table[tableIdx].Len = vlen
		this.table[tableIdx].Offset = offset
		this.table[tableIdx].Index = uint64(pos)
		this.table[tableIdx].HashA = hasha
		this.table[tableIdx].HashB = hashb
	} else {
		entity := Entity{
			Len:    vlen,
			Offset: offset,
			Index:  uint64(pos),
			HashA:  hasha,
			HashB:  hashb,
		}
		this.table = append(this.table, entity)
	}
	// buckets step1
	this.mpqhash[pos].nPosInt = tableIdx
	this.mpqhash[pos].nHashA = hasha
	this.mpqhash[pos].nHashB = hashb

	this.length += 1
	this.size += 1

	return nil
}

// 初始化DealMap.mpqhash
func (this *DealMmap) PushEntity() error {
	if this.buckets == 0 {
		return fmt.Errorf("The dealmap had not init.")
	}

	var idx uint64
	for pos, entity := range this.table {
		idx = entity.Index
		this.mpqhash[idx].nHashA = entity.HashA
		this.mpqhash[idx].nHashB = entity.HashB
		this.mpqhash[idx].nPosInt = int64(pos)
		this.size += 1
	}

	return nil
}

// get entity from mpqhash table
func (this *DealMmap) GetEntity(key string) (*Entity, error) {
	if this.size == 0 {
		return nil, fmt.Errorf("The mpqhash table had none.")
	}
	// 计算hash值
	hashValue := HashKey(key, HASH_O)
	hasha := HashKey(key, HASH_A)
	hashb := HashKey(key, HASH_B)

	start := hashValue % uint64(this.buckets)
	hashPos := start

	var pos int64 = DefaultPos

	for DefaultPos != this.mpqhash[hashPos].nPosInt {
		if this.mpqhash[hashPos].nHashA == hasha && this.mpqhash[hashPos].nHashB == hashb {
			pos = this.mpqhash[hashPos].nPosInt
			break
		}

		hashPos = (hashPos + 1) % uint64(this.buckets)
		if hashPos == start {
			break
		}
	}

	if pos == DefaultPos {
		return nil, nil
	}

	return &this.table[pos], nil
}
