package main

import (
	"bytes"
	"fmt"
	"golang.org/x/text/transform"
	"io"
)

// Replacer型の定義
type Replacer struct {
	old, new []byte // 置換前と置換後
	preDst []byte
}

// Resetメソッド: (状態をリセットする)
func (r *Replacer) Reset() {
	r.preDst = nil
}

// Transformメソッド: oldをnewに置換する
func (r *Replacer) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	// 前回書ききれなかったデータがあれば、それを最初に書き込む
	if len(r.preDst) > 0 {
		n := copy(dst[nDst:], r.preDst)
		nDst += n
		r.preDst = r.preDst[n:]
	}

	// メインの置換処理
	for {
		// srcのnSrc番目からr.oldを探す
		i := bytes.Index(src[nSrc:], r.old)
		// 見つからなかった場合
		if i == -1 {
			// 残りの部分をそのままコピーする
			n := copy(dst[nDst:], src[nSrc:])
			nDst += n
			nSrc += n
			return nDst, nSrc, nil
		}

		// 見つけた位置までコピーして書き込む
		n := copy(dst[nDst:], src[nSrc:nSrc+i])
		nDst += n
		nSrc += n

		// コピーが足りなかった場合は、残りをpreDstに退避させる
		if n < i {
			// 書き込めなかった分を退避
			r.preDst = append(r.preDst, src[nSrc+n:]...)
			err = transform.ErrShortDst
			return nDst, nSrc, err
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
	r := NewReplacer([]byte("郷"), []byte("Go"))

	// Transformに基づいた出力を作成するWriter
	w := transform.NewWriter(io.Discard, r)

	// 入力文字を変換し、結果を標準出力に出力
	input := "郷に入って郷に従え"
	n, err := io.Copy(w, bytes.NewReader([]byte(input)))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	//　書き込んだバイト数を表示
	fmt.Println("Written bytes:", n)
}