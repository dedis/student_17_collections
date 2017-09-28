package collection

import "testing"
import "crypto/sha256"
import "encoding/hex"

func TestHash(test *testing.T) {
    equal := func(reference string, buffer [sha256.Size]byte) bool {
        refbuffer, _ := hex.DecodeString(reference)

        for i := 0; i < sha256.Size; i++ {
            if buffer[i] != refbuffer[i] {
                return false
            }
        }

        return true
    }

    {
        oneint := struct {
            i int
        }{33}

        buffer, err := hash(&oneint)

        if err != nil {
            test.Error("[oneint]", err)
        } else if !equal("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", buffer) {
            test.Error("[oneint] SHA256 inconsistency.")
        }
    }
}
