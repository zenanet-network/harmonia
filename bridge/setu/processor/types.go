package processor

import (
	"fmt"
	"math/big"
)

// HeaderBlock header block
type HeaderBlock struct {
	start  uint64
	end    uint64
	number *big.Int
}

// ContractCheckpoint contract checkpoint
type ContractCheckpoint struct {
	newStart           uint64
	newEnd             uint64
	currentHeaderBlock *HeaderBlock
}

func (c ContractCheckpoint) String() string {
	return fmt.Sprintf("newStart: %v, newEnd %v, contractStart: %v, contractEnd %v, contractHeaderNumber %v",
		c.newStart, c.newEnd, c.currentHeaderBlock.start, c.currentHeaderBlock.end, c.currentHeaderBlock.number)
}

// HarmoniaCheckpoint harmonia checkpoint
type HarmoniaCheckpoint struct {
	start uint64
	end   uint64
}

// NewContractCheckpoint creates contract checkpoint
func NewContractCheckpoint(_newStart uint64, _newEnd uint64, _currentHeaderBlock *HeaderBlock) *ContractCheckpoint {
	return &ContractCheckpoint{
		newStart:           _newStart,
		newEnd:             _newEnd,
		currentHeaderBlock: _currentHeaderBlock,
	}
}

// NewHarmoniaCheckpoint creates new harmonia checkpoint object
func NewHarmoniaCheckpoint(_start uint64, _end uint64) *HarmoniaCheckpoint {
	return &HarmoniaCheckpoint{
		start: _start,
		end:   _end,
	}
}
