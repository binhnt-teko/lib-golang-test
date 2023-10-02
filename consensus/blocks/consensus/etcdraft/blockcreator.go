/*
Copyright IBM Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package etcdraft

import (
	"github.com/blcvn/lib-golang-test/consensus/blocks/consensus/protoutil"
	cb "github.com/blcvn/lib-golang-test/consensus/blocks/types/common"
	"github.com/blcvn/lib-golang-test/log/flogging"
	"github.com/golang/protobuf/proto"
)

// blockCreator holds number and hash of latest block
// so that next block will be created based on it.
type blockCreator struct {
	hash   []byte
	number uint64

	logger *flogging.FabricLogger
}

func (bc *blockCreator) createNextBlock(envs []*cb.Envelope) *cb.Block {
	data := &cb.BlockData{
		Data: make([][]byte, len(envs)),
	}

	var err error
	for i, env := range envs {
		data.Data[i], err = proto.Marshal(env)
		if err != nil {
			bc.logger.Panicf("Could not marshal envelope: %s", err)
		}
	}

	bc.number++

	block := protoutil.NewBlock(bc.number, bc.hash)
	block.Header.DataHash = protoutil.BlockDataHash(data)
	block.Data = data

	bc.hash = protoutil.BlockHeaderHash(block.Header)
	return block
}
