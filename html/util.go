package html

import (
	"fmt"
	"github.com/nine-lives-later/octopus-deploy-decryptor/pkg/projectExport"
)

func decryptedValue(v projectExport.Decryptable, key []byte) string {
	vv, err := v.DecryptedValue(key)
	if err != nil {
		return fmt.Sprintf("DECRYPT ERROR: %s", err.Error())
	}

	return vv
}
