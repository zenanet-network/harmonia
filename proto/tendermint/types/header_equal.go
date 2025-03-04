package types

import (
	"bytes"
)

// Equal 메서드는 두 Header가 동일한지 비교합니다.
func (h *Header) Equal(that interface{}) bool {
	if that == nil {
		return h == nil
	}

	that1, ok := that.(*Header)
	if !ok {
		that2, ok := that.(Header)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return h == nil
	} else if h == nil {
		return false
	}
	
	// Version 비교
	if h.Version.Block != that1.Version.Block || h.Version.App != that1.Version.App {
		return false
	}
	
	// 기본 필드 비교
	if h.ChainID != that1.ChainID || h.Height != that1.Height {
		return false
	}
	
	// Time 비교 (나노초 단위까지)
	if !h.Time.Equal(that1.Time) {
		return false
	}
	
	// LastBlockId 비교
	if !bytes.Equal(h.LastBlockId.Hash, that1.LastBlockId.Hash) {
		return false
	}
	if h.LastBlockId.PartSetHeader.Total != that1.LastBlockId.PartSetHeader.Total {
		return false
	}
	if !bytes.Equal(h.LastBlockId.PartSetHeader.Hash, that1.LastBlockId.PartSetHeader.Hash) {
		return false
	}
	
	// 해시 비교
	if !bytes.Equal(h.LastCommitHash, that1.LastCommitHash) {
		return false
	}
	if !bytes.Equal(h.DataHash, that1.DataHash) {
		return false
	}
	if !bytes.Equal(h.ValidatorsHash, that1.ValidatorsHash) {
		return false
	}
	if !bytes.Equal(h.NextValidatorsHash, that1.NextValidatorsHash) {
		return false
	}
	if !bytes.Equal(h.ConsensusHash, that1.ConsensusHash) {
		return false
	}
	if !bytes.Equal(h.AppHash, that1.AppHash) {
		return false
	}
	if !bytes.Equal(h.LastResultsHash, that1.LastResultsHash) {
		return false
	}
	if !bytes.Equal(h.EvidenceHash, that1.EvidenceHash) {
		return false
	}
	if !bytes.Equal(h.ProposerAddress, that1.ProposerAddress) {
		return false
	}
	
	return true
} 