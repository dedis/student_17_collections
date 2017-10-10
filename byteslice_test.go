package collection

import "testing"
import "encoding/hex"

func TestBytesliceEqual(test *testing.T) {
    lho, _ := hex.DecodeString("85f46bd1ba19d1014b1179edd451ece95296e4a8c765ba8bba86c16893906398")
    rho, _ := hex.DecodeString("85f46bd1ba19d1014b1179edd451ece95296e4a8c765ba8bba86c16893906398")
    cutrho, _ := hex.DecodeString("85f46bd1ba19d1014b1179edd451ece95296e4a8c765ba8bba86c168939063")
    alterrho, _ := hex.DecodeString("85f46bd1ba19d1014b1179edd452ece95296e4a8c765ba8bba86c16893906398")

    if !(equal(lho, rho)) {
        test.Error("[byteslice.go]", "[equal]", "equal() returns false on two equal buffers.")
    }

    if equal(lho, cutrho) {
        test.Error("[byteslice.go]", "[equal]", "equal() returns true on two buffers of different length.")
    }

    if equal(lho, alterrho) {
        test.Error("[byteslice.go]", "[equal]", "equal() returns true on two different buffers.")
    }
}

func TestBytesliceBit(test *testing.T) {
    buffer, _ := hex.DecodeString("85f46bd1ba19d1014b1179edd451ece95296e4a8c765ba8bba86c16893906398")
    reference := "1000010111110100011010111101000110111010000110011101000100000001010010110001000101111001111011011101010001010001111011001110100101010010100101101110010010101000110001110110010110111010100010111011101010000110110000010110100010010011100100000110001110011000"

    for index := 0; index < 8 * len(buffer); index++ {
        bit := bit(buffer, index)

        if (bit && reference[index : index + 1] == "0") || (!bit && reference[index : index + 1] == "1") {
            test.Error("[byteslice.go]", "[bit]", "Wrong bit detected on test buffer.")
            break
        }
    }
}

func TestBytesliceMatch(test *testing.T) {
    min := func(lho, rho int) int {
        if lho < rho {
            return lho
        } else {
            return rho
        }
    }

    type round struct {
        lho string
        rho string
        bits int
    }

    rounds := []round{
        {"85f46bd1ba19d1014b1179edd451ece95296e4a8c765ba8bba86c16893906398", "85f46bd1ba19d1014b1179edd451ece95296e4a8c765ba8bba86c16893906398", 256},
        {"fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0", "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 252},
        {"85f46bd1ba19d1014b1179edd451ece95296e4a8c765ba8bba86c16893906390", "85f46bd1ba19d1014b1179edd451ece95296e4a8c765ba8bba86c16893906398", 252},
        {"85f46bd1ba1ad1014b1179edd451ece95296e4a8c765ba8bba86c16893906398", "85f46bd1ba19d1014b1179edd451ece95296e4a8c765ba8bba86c16893906398", 46},
        {"85f46bd1ba18d1014b1179edd451ece95296e4a8c765ba8bba86c16893906398", "85f46bd1ba19d1014b1179edd451ece95296e4a8c765ba8bba86c16893906398", 47},
    }

    for _, round := range(rounds) {
        lho, _ := hex.DecodeString(round.lho)
        rho, _ := hex.DecodeString(round.rho)

        for index := 0; index <= 8 * min(len(lho), len(rho)); index++ {
            if (match(lho, rho, index) && index > round.bits) || (!match(lho, rho, index) && index <= round.bits) {
                test.Error("[byteslice.go]", "[match]", "Wrong matching on test buffers.")
            }
        }
    }
}
