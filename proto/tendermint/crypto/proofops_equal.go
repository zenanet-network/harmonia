package crypto

import (
	"bytes"
)

// Equal 메서드는 두 ProofOps가 동일한지 비교합니다.
func (po *ProofOps) Equal(that interface{}) bool {
	if that == nil {
		return po == nil
	}

	that1, ok := that.(*ProofOps)
	if !ok {
		that2, ok := that.(ProofOps)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return po == nil
	} else if po == nil {
		return false
	}
	
	if len(po.Ops) != len(that1.Ops) {
		return false
	}
	
	for i, op := range po.Ops {
		if op.Type != that1.Ops[i].Type {
			return false
		}
		if !bytes.Equal(op.Key, that1.Ops[i].Key) {
			return false
		}
		if !bytes.Equal(op.Data, that1.Ops[i].Data) {
			return false
		}
	}
	
	return true
} 