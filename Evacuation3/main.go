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
	old, new []byte
	preDst []byte // 前回書き込めなかったデータ
	preSrc []byte // 前回余ったold分
}

// Replacerの状態をリセットする
func (r *Replacer) Reset() {
	r.preDst = nil
	r.preSrc = nil
}

// transformメソッドで実際の置換処理をする
func (r *Replacer) transform(dst, src []byte, atEOF bool) (nDst, nSrc int, preSrc []byte, err error) {
	for {
		i := bytes.Index(src[nSrc:], r.old)
		if i == -1 {
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

func (r *Replacer) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	// srcの前方にpreSrcを付加する
	_src := src
	if len(r.preSrc) > 0 {
		// preSrcをsrcの前に追加
		_src = append(r.preSrc, src...)
	}
	// 実際の置換処理を行う
	nDst, nSrc, preSrc, err := r.transform(dst, _src, atEOF)

	// 読み込んだ長さより退避させてた長さが長い場合
	if nSrc < len(r.preSrc) {
		// preSrcに残すデータをあたらに退避させる
		r.preSrc = r.preSrc[nSrc:]
		nSrc = 0
	} else {
		// 新たに余った余った部分をpreSrcに退避
		nSrc -= len(r.preSrc)
		r.preSrc = preSrc
	}
	return nDst, nSrc, err
}

func NewReplacer(old, new []byte) *Replacer {
	return &Replacer{old: old, new: new}
}

func main() {
	r := NewReplacer([]byte("apple"), []byte("orange"))

	w := transform.NewWriter(os.Stdout, r)

	// 入力文字列を変換し、標準出力に出力
	input := "i like apple pie apple"
	_, err := io.Copy(w, bytes.NewReader([]byte(input)))
	if err != nil {
		fmt.Println("Error:", err)
	}
}