package main

import (
	"bytes"
	"fmt"
)

// aとbで先頭からマッチする長さ
func overlapWidth(a, b []byte) int {
	w := len(a)
	if w > len(b) {
		w = len(b)
	}

	for ; w > 0; w-- {
		if bytes.Equal(a[len(a)-w:], b[:w]) {
			return w
		}
	}
	// 全くマッチしなかった場合
	return 0
}

func main() {
	// テストケース
	a := []byte("hello fo")
	b := []byte("fo")

	// 部分一致を探す
	matchedLength := overlapWidth(a, b)
	fmt.Printf("Matched Length: %d\n", matchedLength)

}