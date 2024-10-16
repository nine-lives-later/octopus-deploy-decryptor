package projectExport

type Decryptable interface {
	DecryptedValue(key []byte) (string, error)
}
