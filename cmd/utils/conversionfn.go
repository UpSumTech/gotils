package utils

func Int32Ptr(i int) *int32 {
	x := int32(i)
	return &x
}

func Int64Ptr(i int) *int64 {
	x := int64(i)
	return &x
}

func BoolPtr(b bool) *bool {
	return &b
}

func StringPtr(s string) *string {
	return &s
}
