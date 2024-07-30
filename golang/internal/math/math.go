package math

func NextPowerOf2(n uint64) uint64 {
	if n == 0 {
		return 1
	}
	// subtract 1 from n (in case we already have a power of 2)
	n--
	// spread all the higher 1 bits to the right
	n |= (n >> 1) // 1000_0000 -> 1100_0000
	n |= (n >> 2) // 1100_0000 -> 1111_0000
	n |= (n >> 4) // 1111_0000 -> 1111_1111
	n |= (n >> 8)
	n |= (n >> 16)
	n |= (n >> 32)
	// add back 1 to get the power of 2
	return n + 1
}
