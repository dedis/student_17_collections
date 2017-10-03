package collection

import "testing"
import csha256 "crypto/sha256"
import "encoding/hex"

func TestMaskMatch(test *testing.T) {
    type round struct {
        buffer string
        mask string
        bits int
        expected bool
    }

    rounds := []round{
        {"85f46bd1ba19d1014b1179edd451ece95296e4a8c765ba8bba86c16893906398", "85f46bd1ba19d1014b1179edd451ece95296e4a8c765ba8bba86c16893906398", 256, true},
        {"85f46bd1ba19d1014b1179edd451ece95296e4a8c765ba8bba86c16893906398", "85f46bd1ba19d1014b1179edd451ece95296e4a8c765ba8bba86c16893906398", 25, true},
        {"fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0", "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 252, true},
        {"fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0", "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 253, false},
        {"fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0", "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 0, true},
        {"85f46bd1ba19d1014b1179edd451ece95296e4a8c765ba8bba86c16893906390", "85f46bd1ba19d1014b1179edd451ece95296e4a8c765ba8bba86c16893906398", 252, true},
        {"85f46bd1ba19d1014b1179edd451ece95296e4a8c765ba8bba86c16893906390", "85f46bd1ba19d1014b1179edd451ece95296e4a8c765ba8bba86c16893906398", 254, false},
        {"85f46bd1ba19d1014b1179edd451ece95296e4a8c765ba8bba86c16893906390", "85f46bd1ba19d1014b1179edd451ece95296e4a8c765ba8bba86c16893906398", 2, true},
        {"85f46bd1ba1ad1014b1179edd451ece95296e4a8c765ba8bba86c16893906398", "85f46bd1ba19d1014b1179edd451ece95296e4a8c765ba8bba86c16893906398", 46, true},
        {"85f46bd1ba1ad1014b1179edd451ece95296e4a8c765ba8bba86c16893906398", "85f46bd1ba19d1014b1179edd451ece95296e4a8c765ba8bba86c16893906398", 47, false},
        {"85f46bd1ba1ad1014b1179edd451ece95296e4a8c765ba8bba86c16893906398", "85f46bd1ba19d1014b1179edd451ece95296e4a8c765ba8bba86c16893906398", 45, true},
        {"85f46bd1ba18d1014b1179edd451ece95296e4a8c765ba8bba86c16893906398", "85f46bd1ba19d1014b1179edd451ece95296e4a8c765ba8bba86c16893906398", 47, true},
        {"85f46bd1ba18d1014b1179edd451ece95296e4a8c765ba8bba86c16893906398", "85f46bd1ba19d1014b1179edd451ece95296e4a8c765ba8bba86c16893906398", 48, false},
        {"85f46bd1ba18d1014b1179edd451ece95296e4a8c765ba8bba86c16893906398", "85f46bd1ba19d1014b1179edd451ece95296e4a8c765ba8bba86c16893906398", 1, true},
    }

    for _, round := range(rounds) {
        maskvalue, _ := hex.DecodeString(round.mask)
        mask := Mask{maskvalue, round.bits}

        bufferslice, _ := hex.DecodeString(round.buffer)
        var buffer [csha256.Size]byte

        for index := 0; index < csha256.Size; index++ {
            buffer[index] = bufferslice[index]
        }

        if mask.match(buffer) != round.expected {
            test.Error("[match]", "Wrong match on reference round.", round.buffer, round.mask, round.bits, round.expected)
        }
    }
}

func TestScopeMatch(test *testing.T) {
    scope := Scope{}

    bufferslice, _ := hex.DecodeString("85f46bd1ba18d1014b1179edd451ece95296e4a8c765ba8bba86c16893906398")
    var buffer [csha256.Size]byte

    for index := 0; index < csha256.Size; index++ {
        buffer[index] = bufferslice[index]
    }

    scope.Default = false
    if scope.match(buffer) {
        test.Error("[default]", "No-mask match succeeds with false default.")
    }

    scope.Default = true
    if !(scope.match(buffer)) {
        test.Error("[default]", "No-mask match fails with true default.")
    }

    nomatch, _ := hex.DecodeString("fa91")
    scope.Masks = append(scope.Masks, Mask{nomatch, 16})

    if scope.match(buffer) {
        test.Error("[match]", "Scope match succeeds on a non-matching mask.")
    }

    maybematch, _ := hex.DecodeString("86")
    scope.Masks = append(scope.Masks, Mask{maybematch, 8})

    if scope.match(buffer) {
        test.Error("[match]", "Scope match succeeds on a non-matching mask.")
    }

    scope.Masks = append(scope.Masks, Mask{maybematch, 6})

    if !(scope.match(buffer)) {
        test.Error("[match]", "Scope match fails on matching mask.")
    }
}
