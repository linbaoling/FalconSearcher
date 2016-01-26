/*****************************************************************************
 *  file name : HashTable.go
 *  author : Wu Yinghao
 *  email  : wyh817@gmail.com
 *
 *  file description : 字符串的Hash实现
 *
******************************************************************************/

package HashTable

const (
	TABLE_LEN  = 0x500
	HASH_O     = 0
	HASH_A     = 1
	HASH_B     = 2
	MaxASCII   = '\u007F'
	DefaultPos = -1
)

var (
	cryptTable   [TABLE_LEN]uint64
	BUCKETS_LIST = []int{17, 37, 79, 163, 331, 673, 1361, 2729, 5471, 10949, 21911, 43853, 87719, 175447, 350899, 701819, 1403641, 2807303, 5614657, 11229331, 22458671, 44917381, 89834777, 179669557, 359339171, 718678369, 1437356741, 2147483647}
)

type MPQTable struct {
	nHashA  uint64
	nHashB  uint64
	nPosInt int64
}

// new mpqtable
func NewMPQTable(size int) []MPQTable {
	table := mallocMPQTable(size)
	for i := 0; i < size; i++ {
		table[i].nHashA, table[i].nHashB = 0, 0
		table[i].nPosInt = DefaultPos
	}
	// XXX init crypTable, 计算hash值需要使用
	initCryptTable()

	return table
}

func (this *MPQTable) GetPos() int64 {
	return this.nPosInt
}

func (this *MPQTable) SetPos(value int64) {
	this.nPosInt = value
}

func (this *MPQTable) SetHashA(value uint64) {
	this.nHashA = value
}

func (this *MPQTable) SetHashB(value uint64) {
	this.nHashB = value
}

// 初始化hash计算需要的基础map table
func initCryptTable() {
	var seed, index1, index2 uint64 = 0x00100001, 0, 0
	i := 0
	for index1 = 0; index1 < 0x100; index1 += 1 {
		for index2, i = index1, 0; i < 5; index2 += 0x100 {
			seed = (seed*125 + 3) % 0x2aaaab
			temp1 := (seed & 0xffff) << 0x10
			seed = (seed*125 + 3) % 0x2aaaab
			temp2 := seed & 0xffff
			cryptTable[index2] = temp1 | temp2
			i += 1
		}
	}
}

// hash, 以及相关校验hash值
func HashKey(lpszString string, dwHashType int) uint64 {
	i, ch := 0, 0
	var seed1, seed2 uint64 = 0x7FED7FED, 0xEEEEEEEE
	var key uint8
	strLen := len(lpszString)
	for i < strLen {
		key = lpszString[i]
		ch = int(toUpper(rune(key)))
		i += 1
		seed1 = cryptTable[(dwHashType<<8)+ch] ^ (seed1 + seed2)
		seed2 = uint64(ch) + seed1 + seed2 + (seed2 << 5) + 3
	}
	return uint64(seed1)
}

func toUpper(r rune) rune {
	if r <= MaxASCII {
		if 'a' <= r && r <= 'z' {
			r -= 'a' - 'A'
		}
		return r
	}
	return r
}
