package crypto

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const TestString = "encryptme"

func TestDecrypt(t *testing.T) {
	t.Run("success decryption", func(t *testing.T) {
		enryptedMess, err := Encrypt(TestString)
		assert.Nil(t, err)
		actual, err := Decrypt(enryptedMess)
		assert.Nil(t, err)
		assert.Equal(t, TestString, actual)
	})
	t.Run("failure decryption", func(t *testing.T) {
		str := fmt.Sprintf("%d%s%s", encV1, "||", TestString)
		_, err := Decrypt(str)
		assert.NotNil(t, err)
	})
}
