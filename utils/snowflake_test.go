package utils

import (
	"fmt"
	"testing"
)

func TestNewSnowflake(t *testing.T) {
	idCli, err := NewSnowflake(1)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(idCli.NextID())
}

func BenchmarkSnowflake(b *testing.B) {
	idCli, err := NewSnowflake(1)
	if err != nil {
		b.Error(err)
	}
	for i := 0; i < b.N; i++ {
		if _, err = idCli.NextID(); err != nil {
			b.Error(err)
		}
	}
}
