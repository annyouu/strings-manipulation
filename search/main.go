package main

import (
	"bytes"
	"fmt"
	"golang.org/x/text/transform"
)

// Replacer型の定義
type Replacer struct {
	transform.NopResetter
	old, new []byte  //置換前と置換後
}

// Transformメソッド: oldをnewに置換する
func (r *Replacer) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	for {
		// srcのnSrc番目からr.oldを探す
		i := bytes.Index(src[nSrc:], r.old)
		// なかった時は文字をそのままコピーして終わる
		if i == -1 {
			n := copy(dst[nDst:], src[nSrc:])
			nDst += n
			nSrc += n
			return nDst, nSrc, nil
		}

		// 見つけた位置までコピーして書き込む
		n := copy(dst[nDst:], src[nSrc:nSrc+i])
		nDst += n
		nSrc += n

		if n < i {
			err = transform.ErrShortDst
			return
		}

		// 置換部分を新しい文字列に置き換える
		n = copy(dst[nDst:], r.new)
		nDst += n
		nSrc += len(r.old)
	}
}

func NewReplacer(old, new []byte) *Replacer {
	return &Replacer{old: old, new: new}
}

func main() {
	old := []byte("郷")
	new := []byte("Go")

	r := NewReplacer(old, new)

	// 入力したやつ
	src := []byte("郷に入っては郷に従え")

	// 出力用バッファ
	dst := make([]byte, len(src))

	// transform.Readerを使って変換
	nDst, _, err := r.Transform(dst, src, true)

	if err != nil {
		fmt.Println("エラー:", err)
	} else {
		fmt.Printf("変換後: %s\n", dst[:nDst])
	}
}