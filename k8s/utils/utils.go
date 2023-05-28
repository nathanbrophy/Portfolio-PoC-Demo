package utils

// Int32PointerGenerator is a wrapper and will return a memory address pointer to the passed in date
func Int32PointerGenerator(x int32) *int32 {
	return &x
}

// Int64PointerGenerator is a wrapper and will return a memory address pointer to the passed in date
func Int64PointerGenerator(x int64) *int64 {
	return &x
}

// StringPointerGenerator is a wrapper and will return a memory address pointer to the passed in date
func StringPointerGenerator(x string) *string {
	return &x
}
