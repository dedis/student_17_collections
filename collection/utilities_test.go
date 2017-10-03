package collection

import "testing"
import "encoding/hex"

func TestBit(test *testing.T) {
    buffer, _ := hex.DecodeString("85f46bd1ba19d1014b1179edd451ece95296e4a8c765ba8bba86c16893906398")

    reference := "1000010111110100011010111101000110111010000110011101000100000001010010110001000101111001111011011101010001010001111011001110100101010010100101101110010010101000110001110110010110111010100010111011101010000110110000010110100010010011100100000110001110011000"
    for index := 0; index < 8 * len(buffer); index++ {
        bit := bit(buffer, index)

        if (bit && reference[index : index + 1] == "0") || (!bit && reference[index : index + 1] == "1") {
            test.Error("[bit]", "Wrong bit detected on test buffer.")
        }
    }
}

func TestMatch(test *testing.T) {
    min := func(a, b int) int {
        if a < b {
            return a
        } else {
            return b
        }
    }

    type round struct {
        buffer string
        reference string
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
        buffer, _ := hex.DecodeString(round.buffer)
        reference, _ := hex.DecodeString(round.reference)

        for index := 0; index <= 8 * min(len(buffer), len(reference)); index++ {
            if (match(buffer, reference, index) && index > round.bits) || (!match(buffer, reference, index) && index <= round.bits) {
                test.Error("[match]", "Wrong matching on test buffers.", index, match(buffer, reference, index))
            }
        }
    }
}
