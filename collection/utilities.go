package collection

func bit(buffer []byte, index int) bool {
    byteidx := uint(index) / 8
    bitidx := 7 - (uint(index) % 8)

    return ((buffer[byteidx] & (uint8(1) << bitidx)) != 0)
}

func match(buffer []byte, reference []byte, bits int) bool {
    for index := 0; index < bits; {
        if index < bits - 8 {
            if buffer[index / 8] != reference[index / 8] {
                return false
            }

            index += 8
        } else {
            if bit(buffer, index) != bit(reference, index) {
                return false
            }

            index++
        }
    }

    return true
}
