package random

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewRandomString(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{name: "size = 1", size: 1},
		{name: "size = 5", size: 5},
		{name: "size = 10", size: 10},
		{name: "size = 20", size: 20},
		{name: "size = 30", size: 30},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str1 := NewRandomString(tt.size)
			time.Sleep(1 * time.Millisecond)
			str2 := NewRandomString(tt.size)

			// Выводим строки, чтобы увидеть, что происходит
			fmt.Printf("Generated strings: str1 = %s, str2 = %s\n", str1, str2)

			// Проверяем длину строк
			assert.Len(t, str1, tt.size)
			assert.Len(t, str2, tt.size)

			// Проверяем, что строки разные
			if str1 == str2 {
				t.Errorf("Expected different strings but got: str1 = %s, str2 = %s", str1, str2)
			}
		})
	}
}
