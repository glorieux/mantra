package slice

import (
	"testing"
)

var (
	a = []int{0, 1, 2, 3}
	b = []int{0, 1, 2, 3}
	c = []int{42, 42, 42, 42}
	d = []int{0, 1, 2}
)

func TestIntDeepEqual(t *testing.T) {
	if IntDeepEqual(a, d) {
		t.Error("Slices of different sizes can't be equal")
	}
	if IntDeepEqual(nil, []int{}) {
		t.Error("Slices can't be nil")
	}
	if !IntDeepEqual(a, b) {
		t.Error("Slices should be equals")
	}
	if IntDeepEqual(a, c) {
		t.Error("Slices should be different")
	}
}

func BenchmarkIntDeepEqual(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IntDeepEqual(a, c)
	}
}

func TestDelete(t *testing.T) {
	arr := []string{"a", "b", "c"}
	cleanArr := Remove(arr, 0)
	if len(cleanArr) > 2 {
		t.Error("Should only contain two elements")
	}
}

func BenchmarkDelete(b *testing.B) {
	arr := []string{"a", "b", "c"}
	for i := 0; i < b.N; i++ {
		Remove(arr, 1)
	}
}

func TestContains(t *testing.T) {
	if !Contains("test", []string{"this", "contains", "test"}) {
		t.Error("Slice contains test")
	}
	if Contains("unknown", []string{"this", "contains", "test"}) {
		t.Error("Slice doesn't contain test")
	}
}
