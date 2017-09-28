package collection

import "crypto/sha256"
import "github.com/dedis/protobuf"

func hash(pointer interface{}) ([sha256.Size]byte, error) {
    buffer, err := protobuf.Encode(pointer)

    if(err != nil) {
        return [sha256.Size]byte{}, err
    } else {
        return sha256.Sum256(buffer), nil
    }
}
