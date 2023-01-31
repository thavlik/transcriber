package s3

func EncodeHash(hash string) string {
	if len(hash) != 32 {
		panic("hash must be 32 character md5")
	}
	// re-encode the hash as a directory structure
	// 12345678901234567890123456789012 -> 12/34/56789012345678901234567890
	// this way each leaf directory will have at most 256 entries
	return hash[:2] + "/" + hash[2:4] + "/" + hash[4:]
}
