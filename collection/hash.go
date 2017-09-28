package collection

import "crypto/sha256"
import "github.com/dedis/protobuf"

/* hash computes the SHA256 hash of a struct. Provided with a pointer to the
   object to hash, it encodes it using protobuf, and computes the SHA256 hash
   of the result.

   Example usage:
        type twonumbers struct {
            i uint32
            j uint32
        }

        buffer, err := hash(&twonumbers{44, 59})
        if(err != nil) {
            fmt.Println("Error:", err)
        } else {
            fmt.Printf("%x\n", buffer) // Prints e3b0c44298fc1c14...
        }
*/
func hash(pointer interface{}) ([sha256.Size]byte, error) {
    buffer, err := protobuf.Encode(pointer)

    if(err != nil) {
        return [sha256.Size]byte{}, err
    } else {
        return sha256.Sum256(buffer), nil
    }
}
