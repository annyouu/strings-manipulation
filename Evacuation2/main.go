package main

import (
	"bytes"
	"fmt"
	"golang.org/x/text/transform"
	"io"
	"os"
)

// Replacerはバイト列の置換を行う
type Replacer struct {
	old, new []byte // 置換前と置換後
	preDst []byte // 前回書き込めなかったデータ
}

// ResetはReplacerの状態をリセットする
func (r *Replacer) Reset() {
	r.preDst = nil
}

// Transformはoldをnewに置換する
func (r *Replacer) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	// 前回書き込めなかった分を書き込む
	if len(r.preDst) > 0 {
		n := copy(dst, r.preDst)
		nDst += n
		r.preDst = r.preDst[n:]
		// それでもまだ足りない時
		if len(r.preDst) > 0 {
			err = transform.ErrShortDst
			return
		}
	}

	// メインの置換処理をする
	for {
		// srcのnSrc番目からr.oldを探す
		i := bytes.Index(src[nSrc:], r.old)
		if i == -1 {
			// 見つからなかった場合、残りの部分をコピーする
			n := copy(dst[nDst:], src[nSrc:])
			nDst += n
			nSrc += n
			return
		}
		
		// 見つけた位置までコピーする
		n := copy(dst[nDst:], src[nSrc:nSrc+i])
		nDst += n
		nSrc += i

		// 置換する文字をコピーして書き込む
		n = copy(dst[nDst:], r.new)
		nDst += n
		nSrc += len(r.old)

		// r.newが長くてdstに書ききれない場合、次回に持ち越し
		if n < len(r.new) {
			r.preDst = r.new[n:]
			err = transform.ErrShortDst
			return
		}
	}
}

func NewReplacer(old, new []byte) *Replacer {
	return &Replacer{old: old, new: new}
}

func main() {
	r := NewReplacer([]byte("郷"), []byte("Go"))

	// 変換を行うWriterを作成
	w := transform.NewWriter(os.Stdout, r)

	// 入力文字列を変換し、標準出力に出力
	input := "郷に入っては郷に従え"
	_, err := io.Copy(w, bytes.NewReader([]byte(input)))
	if err != nil {
		fmt.Println("Error:", err)
	}
}