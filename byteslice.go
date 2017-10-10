package collection

func equal(lho []byte, rho []byte) bool {
    if len(lho) != len(rho) {
        return false
    }

    for index := 0; index < len(lho); index++ {
        if lho[index] != rho[index] {
            return false
        }
    }

    return true
}

func bit(buffer []byte, index int) bool {
    byteidx := uint(index) / 8
    bitidx := 7 - (uint(index) % 8)

    return ((buffer[byteidx] & (uint8(1) << bitidx)) != 0)
}

func match(lho []byte, rho []byte, bits int) bool {
    for index := 0; index < bits; {
        if index < bits - 8 {
            if lho[index / 8] != rho[index / 8] {
                return false
            }

            index += 8
        } else {
            if bit(lho, index) != bit(rho, index) {
                return false
            }

            index++
        }
    }

    return true
}
