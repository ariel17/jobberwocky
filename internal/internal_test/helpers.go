package internal_test

func BoolPointer(v bool) *bool {
	newValue := v
	return &newValue
}