package utils

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestInt32PointerGenerator(t *testing.T) {
	fuzz := rand.Int31()
	type args struct {
		x int32
	}
	tests := []struct {
		name string
		args args
		want int32
	}{
		{
			name: "default scenario",
			args: args{
				x: fuzz,
			},
			want: fuzz,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Int32PointerGenerator(tt.args.x); *got != tt.want {
				t.Errorf("Int32PointerGenerator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64PointerGenerator(t *testing.T) {
	fuzz := rand.Int63()
	type args struct {
		x int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "default scenario",
			args: args{
				x: fuzz,
			},
			want: fuzz,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Int64PointerGenerator(tt.args.x); *got != tt.want {
				t.Errorf("Int64PointerGenerator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringPointerGenerator(t *testing.T) {
	fuzz := fmt.Sprintf("test-%d", rand.Int31())
	type args struct {
		x string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "default scenario",
			args: args{
				x: fuzz,
			},
			want: fuzz,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringPointerGenerator(tt.args.x); *got != tt.want {
				t.Errorf("StringPointerGenerator() = %v, want %v", got, tt.want)
			}
		})
	}
}
