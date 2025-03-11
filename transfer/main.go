package main

import (
	"bytes"
	"fmt"
	"golang.org/x/text/transform"
	"io"
	"os"
	"strings"
)

// Replacer型の定義
type Replacer struct {
	transform.NopResetter
	old, new []byte
}

// Transformメソッド: old を new に置換する
func (r *Replacer) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	for {
		// src[nSrc:] から r.old を探す
		i := bytes.Index(src[nSrc:], r.old)
		// 見つからなかった場合、残りをそのままコピーして終了
		if i == -1 {
			n := copy(dst[nDst:], src[nSrc:])
			nDst += n
			nSrc += n
			return nDst, nSrc, nil
		}

		// 見つけた位置までをコピー
		n := copy(dst[nDst:], src[nSrc:nSrc+i])
		nDst += n
		nSrc += n

		// バッファが足りない場合
		if n < i {
			err = transform.ErrShortDst
			return
		}

		// 置換部分を新しいバイト列に書き込む
		n = copy(dst[nDst:], r.new)
		nDst += n
		nSrc += len(r.old)
	}
}

func NewReplacer(old, new []byte) *Replacer {
	return &Replacer{old: old, new: new}
}

func main() {
	t := NewReplacer([]byte("郷"), []byte("Go"))

	w := transform.NewWriter(os.Stdout, t)

	input := "郷に入っては郷に従え"

	_, err := io.Copy(w, strings.NewReader(input))
	if err != nil {
		fmt.Println("エラー:", err)
	}
}