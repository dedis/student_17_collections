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
        mask := mask{maskvalue, round.bits}

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

func TestScopeMethods(test *testing.T) {
    scope := scope{}

    bufferslice, _ := hex.DecodeString("1234567890")
    scope.Add(bufferslice, 3)

    if len(scope.masks) != 1 {
        test.Error("[add]", "Add does not add to masks.")
    }

    if !match(scope.masks[0].value, bufferslice, 24) || scope.masks[0].bits != 3 {
        test.Error("[add]", "Add adds wrong mask.")
    }

    bufferslice, _ = hex.DecodeString("0987654321")
    scope.Add(bufferslice, 40)

    if len(scope.masks) != 2 {
        test.Error("[add]", "Add does not add to masks.")
    }

    if !match(scope.masks[1].value, bufferslice, 40) || scope.masks[1].bits != 40 {
        test.Error("[add]", "Add adds wrong mask.")
    }

    scope.All()

    if len(scope.masks) != 0 {
        test.Error("[all]", "All does not wipe masks.")
    }

    if !(scope.all) {
        test.Error("[all]", "All does not properly set the all flag.")
    }

    scope.Add(bufferslice, 40)
    scope.None()

    if len(scope.masks) != 0 {
        test.Error("[none]", "None does not wipe masks.")
    }

    if scope.all {
        test.Error("[none]", "None does not properly set the all flag.")
    }
}

func TestScopeMatch(test *testing.T) {
    scope := scope{}

    bufferslice, _ := hex.DecodeString("85f46bd1ba18d1014b1179edd451ece95296e4a8c765ba8bba86c16893906398")
    var buffer [csha256.Size]byte

    for index := 0; index < csha256.Size; index++ {
        buffer[index] = bufferslice[index]
    }

    scope.None()
    if scope.match(buffer) {
        test.Error("[all]", "No-mask match succeeds after None().")
    }

    scope.All()
    if !(scope.match(buffer)) {
        test.Error("[all]", "No-mask match fails after All().")
    }

    nomatch, _ := hex.DecodeString("fa91")
    scope.Add(nomatch, 16)

    if scope.match(buffer) {
        test.Error("[match]", "Scope match succeeds on a non-matching mask.")
    }

    maybematch, _ := hex.DecodeString("86")
    scope.Add(maybematch, 8)

    if scope.match(buffer) {
        test.Error("[match]", "Scope match succeeds on a non-matching mask.")
    }

    scope.Add(maybematch, 6)

    if !(scope.match(buffer)) {
        test.Error("[match]", "Scope match fails on matching mask.")
    }
}
