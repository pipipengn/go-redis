package consistenthash

import (
	"hash/crc32"
	"sort"
)

type HashFunc func(data []byte) uint32

type NodeMap struct {
	hashFunc    HashFunc
	nodeHashs   []int
	nodeHashMap map[int]string
}

func NewEmptyNodeMap(f HashFunc) *NodeMap {
	m := &NodeMap{
		hashFunc:    f,
		nodeHashs:   []int{},
		nodeHashMap: map[int]string{},
	}
	if m.hashFunc == nil {
		m.hashFunc = crc32.ChecksumIEEE
	}
	return m
}

func (m *NodeMap) IsEmpty() bool {
	return len(m.nodeHashs) == 0
}

func (m *NodeMap) AddNodes(nodes ...string) {
	for _, node := range nodes {
		if node == "" {
			continue
		}
		hash := int(m.hashFunc([]byte(node)))
		m.nodeHashs = append(m.nodeHashs, hash)
		m.nodeHashMap[hash] = node
	}
	sort.Ints(m.nodeHashs)
}

func (m *NodeMap) PickNode(key string) string {
	if m.IsEmpty() {
		return ""
	}
	hash := int(m.hashFunc([]byte(key)))
	idx := sort.Search(len(m.nodeHashs), func(i int) bool {
		return m.nodeHashs[i] >= hash
	})
	if idx == len(m.nodeHashs) {
		idx = 0
	}
	return m.nodeHashMap[m.nodeHashs[idx]]
}
