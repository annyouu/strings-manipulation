package main

import (
	"bytes"
	"fmt"
	"golang.org/x/text/transform"
	"io"
	"strings"
)

// Replacerは指定された文字列を新しい文字列に置換する
type Replacer struct {
	old []byte
	new []byte
}

// NewReplacerはReplacerを初期化するコンストラクタ
func NewReplacer(old, new string) *Replacer {
	return &Replacer{
		old: []byte(old),
		new: []byte(new),
	}
}

// Transformメソッド(Transformerインターフェースの実装)
func (r *Replacer) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	for {
		// src[nSrc:]からr.oldを探す
		i := bytes.Index(src[nSrc:], r.old)
		// なかったら、そのままコピー
		if i == -1 {
			n := len(src[nSrc:])
			var w int
			if !atEOF {
				// srcの末尾がr.oldの前方部分と一致するかチェック
				w = overlapWidth(src[nSrc:], r.old)
				if w > 0 {
					n -= w
					err = transform.ErrShortSrc
				}
			}

			m := copy(dst[nDst:], src[nSrc:nSrc+n])
			nDst += m
			nSrc += m
			if m < n {
				err = transform.ErrShortDst
				return
			}
			// 余った分があればあればここで処理する
			nSrc += w
			return nDst, nSrc, err
		}

		// r.oldが見つかった位置までdstにコピー
		m := copy(dst[nDst:], src[nSrc:nSrc+i])
		nDst += m
		nSrc += i

		// r.oldをr.newに置換してdstにコピー
		m = copy(dst[nDst:], r.new)
		nDst += m
		nSrc += len(r.old)
		// もし、r.newがdst書ききれなければエラーを返す
		if m < len(r.new) {
			err = transform.ErrShortDst
			return nDst, nSrc, err
		}
	}
}

// srcの末尾まで、r.oldの前方部分とどれだけ一致するかを計算する
func overlapWidth(src, old []byte) int {
	max := len(src)
	if len(old) < max {
		max = len(old)
	}
	for i := max; i > 0; i-- {
		if bytes.HasPrefix(old, src[len(src)-i:]) {
			return i
		}
	}
	return 0
}

func (r *Replacer) Reset() {}

func main() {
	replacer := NewReplacer("foo", "bar")
	input := "hello foo world! foo is nice."
	reader := transform.NewReader(strings.NewReader(input), replacer)
	writer := &bytes.Buffer{}

	_, err := io.Copy(writer, reader)
	if err != nil && err != io.EOF {
		fmt.Println("変換エラー:", err)
	} else {
		fmt.Println("変換結果:", writer.String())
	}
}