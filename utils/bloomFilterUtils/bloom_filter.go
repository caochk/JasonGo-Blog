package bloomFilterUtils

import (
	"fmt"
	"github.com/bits-and-blooms/bitset"
	"github.com/spaolacci/murmur3"
)

type BloomFilter struct {
	m      uint64         // 位图的大小
	k      uint32         // 哈希函数个数
	bitset *bitset.BitSet // 布隆过滤器的基础位图
}

// NewBloomFilter New 新建布隆过滤器
func NewBloomFilter(m uint64, k uint32) *BloomFilter {
	return &BloomFilter{m: m, k: k, bitset: bitset.New(uint(m))}
}

// Add 往布隆过滤器中增加元素
func (f BloomFilter) Add(data string) {
	for i := uint32(0); i < f.k; i++ {
		idx := f.getLocation(data, i)
		fmt.Println("1", idx)
		f.bitset.Set(f.getLocation(data, i)) // 将位数组中某索引（有哈希函数计算）对应值置一
	}
}

// Exist 布隆过滤器中查找元素
func (f BloomFilter) Exist(data string) bool {
	for i := uint32(0); i < f.k; i++ {
		idx := f.getLocation(data, i)
		fmt.Println("2", idx)
		if !f.bitset.Test(f.getLocation(data, i)) {
			return false // 只要位数组中有一个索引不是1，就表明该key不在位数组中（即不在缓存或数据库中）
		}
	}
	return true
}

// 根据哈希函数计算哈希值，再对位数组长度取模后之值用作位数组索引值
func (f BloomFilter) getLocation(data string, i uint32) uint {
	return getHashValue(data, i) % uint(f.m)
}

// 计算得到哈希值
func getHashValue(data string, seed uint32) uint {
	m := murmur3.New64WithSeed(seed)
	_, _ = m.Write([]byte(data))
	return uint(m.Sum64())
}
