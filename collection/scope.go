package collection

import csha256 "crypto/sha256"

// Mask

type Mask struct {
    Value []byte
    Bits int
}

// Methods

func (mask Mask) match(buffer [csha256.Size]byte) bool {
    return match(buffer[:], mask.Value[:], mask.Bits)
}

// Scope

type Scope struct {
    Masks []Mask
    Default bool
}

// Methods

func (scope Scope) match(buffer [csha256.Size]byte) bool {
    if len(scope.Masks) == 0 {
        return scope.Default
    }

    for index := 0; index < len(scope.Masks); index++ {
        if scope.Masks[index].match(buffer) {
            return true
        }
    }

    return false
}
