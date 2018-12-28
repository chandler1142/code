package core

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

//Block keeps block header
type Block struct {
	Timestamp     int64  //区块创建时间戳
	Data          []byte //区块包含的数据
	PrevBlockHash []byte //前一个节点的HASH值
	Hash          []byte //区块自身的HASH值，用于校验数据有效
}

//新建一个区块
func NewBlock(data string, prevBlockHash []byte) *Block {
	//todo 对象创建，返回值等
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}}
	block.SetHash()
	return block
}

func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	//todo 看一下这个Join语法
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)
	//todo
	b.Hash = hash[:]
}

//创世纪块
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}
