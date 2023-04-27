package internal_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BoolPointer(v bool) *bool {
	newValue := v
	return &newValue
}

func GetGoldenFile(t *testing.T, name string) string {
	content, err := os.ReadFile(fmt.Sprintf("goldenfiles/%s.json", name))
	assert.Nil(t, err)
	return string(content)
}