package collection

import "reflect"
import csha256 "crypto/sha256"
import "encoding/binary"

/*
    sha256 computes the SHA256 hash of one or more values in a way that takes
    into account their types, their order, and their structure.

    sha256() accepts one or more values of the following types:
     - Fixed size signed integers: `int8`, `int16`, `int32`, and `int64`.
     - Fixed size unsigned integers: `uint8`, `uint16`, `uint32`, and `uint64`.
     - `string`
     - N-dimensional slices (with N >= 1) of any of the above.

    For example, the following are valid calls to sha256:
     - `sha256(uint64(33))`
     - `sha256("Hello World!")`
     - `sha256(uint64(12), int32(-4), "What", "a", "nice", "day")`
     - `sha256([]uint64{1, 2, 3, 4, 5})`
     - `sha256([][]uint8{{0, 1, 2}, {3, 4, 5}, {6, 7, 8}})`

    But the following are *not* valid calls to sha256:
     - `sha256(33)` (33 is an `int`, whose size depends on the environment).
     - `sha256([4]int8{1, 2, 3, 4})` (arrays are not accepted, only slices).
     - `sha256(my_fancy_object)` (structs are not accepted)
     - `sha256(&my_uint8)` (pointers are not accepted)

    This hash function is endianess-independent, and its value is determined
    also by the types and the structure of the values provided. For example:
     - `sha256(int8(44)) != sha256(uint8(44))`
     - `sha256(uint8(44), uint8(55)) != sha256([]uint8{44, 55})`
     - `sha256("Hello World!") != sha256([]uint8("Hello World!"))`

    sha256 first encodes the values provided into a buffer, then returns the
    sha256 hash of the buffer. The encoding follows the following rules:
     - All arithmetic values are encoded in Big Endian format.
     - Before encoding any value, a one-byte progressive constant is encoded
       to identify the type of the value. The constant assumes values 0, 1, ...
       on `bool`, `int8`, `int16`, `int32`, `int64`, `uint8`, `uint16`,
       `uint32`, `uint64`, `[]bool`, `[]int8`, `[]int16`, `[]int32`, `[]uint8`,
       `[]uint16`, `[]uint32`, `[]uint64`, and `string` respectively.
     - Before encoding a slice value, an `uint64` is encoded to identify
       the length of the slice.
     - Slices of slices use a special one-byte constant value, and a `uint64`
       value to identify their length. Encoding is then done by recurring
       depth-first on each of their elements.
*/
func sha256(item interface{}, items... interface{}) [csha256.Size]byte {
    const(
        boolid = iota
        int8id
        int16id
        int32id
        int64id
        uint8id
        uint16id
        uint32id
        uint64id

        boolsliceid
        int8sliceid
        int16sliceid
        int32sliceid
        int64sliceid
        uint8sliceid
        uint16sliceid
        uint32sliceid
        uint64sliceid

        stringid
        arrayid
    )

    var size func(interface{}) int
    size = func(item interface{}) int {
        switch value := item.(type) {
        case bool:
            return 2
        case int8:
            return 2
        case int16:
            return 3
        case int32:
            return 5
        case int64:
            return 9
        case uint8:
            return 2
        case uint16:
            return 3
        case uint32:
            return 5
        case uint64:
            return 9
        case []bool:
            return 9 + len(value)
        case []int8:
            return 9 + len(value)
        case []int16:
            return 9 + 2 * len(value)
        case []int32:
            return 9 + 4 * len(value)
        case []int64:
            return 9 + 8 * len(value)
        case []uint8:
            return 9 + len(value)
        case []uint16:
            return 9 + 2 * len(value)
        case []uint32:
            return 9 + 4 * len(value)
        case []uint64:
            return 9 + 8 * len(value)
        case string:
            return 9 + len(value)
        }

        reflection := reflect.ValueOf(item)

        if reflection.Kind() == reflect.Slice {
            total := 9

            for index := 0; index < reflection.Len(); index++ {
                total += size(reflection.Index(index).Interface())
            }

            return total
        }

        panic("hash() only accepts: bool, int8, int16, int32, int64, uint8, uint16, uint32, uint64 and string (or N-dimensional slices of those types, with N >= 1).")
    }

    alloc := size(item)
    for _, oitem := range(items) {
        alloc += size(oitem)
    }

    buffer := make([]byte, alloc)

    var write func([]byte, interface{}) []byte
    write = func(buffer []byte, item interface{}) []byte {
        switch value := item.(type) {
        case bool:
            buffer[0] = boolid
            if value {
                buffer[1] = 1
            } else {
                buffer[1] = 0
            }

            return buffer[2:]

        case int8:
            buffer[0] = int8id
            buffer[1] = byte(value)
            return buffer[2:]

        case int16:
            buffer[0] = int16id
            binary.BigEndian.PutUint16(buffer[1:], uint16(value))
            return buffer[3:]

        case int32:
            buffer[0] = int32id
            binary.BigEndian.PutUint32(buffer[1:], uint32(value))
            return buffer[5:]

        case int64:
            buffer[0] = int64id
            binary.BigEndian.PutUint64(buffer[1:], uint64(value))
            return buffer[9:]

        case uint8:
            buffer[0] = uint8id
            buffer[1] = byte(value)
            return buffer[2:]

        case uint16:
            buffer[0] = uint16id
            binary.BigEndian.PutUint16(buffer[1:], value)
            return buffer[3:]

        case uint32:
            buffer[0] = uint32id
            binary.BigEndian.PutUint32(buffer[1:], value)
            return buffer[5:]

        case uint64:
            buffer[0] = uint64id
            binary.BigEndian.PutUint64(buffer[1:], value)
            return buffer[9:]

        case []bool:
            buffer[0] = boolsliceid
            binary.BigEndian.PutUint64(buffer[1:], uint64(len(value)))

            for index := 0; index < len(value); index++ {
                if value[index] {
                    buffer[9 + index] = 1
                } else {
                    buffer[9 + index] = 0
                }
            }

            return buffer[9 + len(value):]

        case []int8:
            buffer[0] = int8sliceid
            binary.BigEndian.PutUint64(buffer[1:], uint64(len(value)))

            for index := 0; index < len(value); index++ {
                buffer[9 + index] = byte(value[index])
            }

            return buffer[9 + len(value):]

        case []int16:
            buffer[0] = int16sliceid
            binary.BigEndian.PutUint64(buffer[1:], uint64(len(value)))

            for index := 0; index < len(value); index++ {
                binary.BigEndian.PutUint16(buffer[9 + 2 * index:], uint16(value[index]))
            }

            return buffer[9 + 2 * len(value):]

        case []int32:
            buffer[0] = int32sliceid
            binary.BigEndian.PutUint64(buffer[1:], uint64(len(value)))

            for index := 0; index < len(value); index++ {
                binary.BigEndian.PutUint32(buffer[9 + 4 * index:], uint32(value[index]))
            }

            return buffer[9 + 4 * len(value):]

        case []int64:
            buffer[0] = int64sliceid
            binary.BigEndian.PutUint64(buffer[1:], uint64(len(value)))

            for index := 0; index < len(value); index++ {
                binary.BigEndian.PutUint64(buffer[9 + 8 * index:], uint64(value[index]))
            }

            return buffer[9 + 8 * len(value):]

        case []uint8:
            buffer[0] = uint8sliceid
            binary.BigEndian.PutUint64(buffer[1:], uint64(len(value)))
            copy(buffer[9:], value)
            return buffer[9 + len(value):]

        case []uint16:
            buffer[0] = uint16sliceid
            binary.BigEndian.PutUint64(buffer[1:], uint64(len(value)))

            for index := 0; index < len(value); index++ {
                binary.BigEndian.PutUint16(buffer[9 + 2 * index:], value[index])
            }

            return buffer[9 + 2 * len(value):]

        case []uint32:
            buffer[0] = uint32sliceid
            binary.BigEndian.PutUint64(buffer[1:], uint64(len(value)))

            for index := 0; index < len(value); index++ {
                binary.BigEndian.PutUint32(buffer[9 + 4 * index:], value[index])
            }

            return buffer[9 + 4 * len(value):]

        case []uint64:
            buffer[0] = uint64sliceid
            binary.BigEndian.PutUint64(buffer[1:], uint64(len(value)))

            for index := 0; index < len(value); index++ {
                binary.BigEndian.PutUint64(buffer[9 + 8 * index:], uint64(value[index]))
            }

            return buffer[9 + 8 * len(value):]

        case string:
            buffer[0] = stringid
            binary.BigEndian.PutUint64(buffer[1:], uint64(len(value)))
            copy(buffer[9:], value)
            return buffer[9 + len(value):]
        }

        reflection := reflect.ValueOf(item)

        buffer[0] = arrayid
        binary.BigEndian.PutUint64(buffer[1:], uint64(reflection.Len()))

        cursor := buffer[9:]

        for index := 0; index < reflection.Len(); index++ {
            cursor = write(cursor, reflection.Index(index).Interface())
        }

        return cursor
    }

    cursor := write(buffer, item)

    for _, variadicitem := range(items) {
        cursor = write(cursor, variadicitem)
    }

    return csha256.Sum256(buffer)
}
