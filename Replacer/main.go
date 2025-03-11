package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/text/transform"
)

// Replacer型の定義
type Replacer struct {
	transform.NopResetter // Resetメソッドのデフォルト
	old, new []byte
}

// Transformメソッド：oldをnewに変換する
func (r *Replacer) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	// r.oldが空なら、そのままコピー
	if len(r.old) == 0 {
		n := copy(dst, src)
		return n, n, nil
	}

	for len(src) > 0 {
		//　変換対象の文字列を検索
		index := strings.Index(string(src), string(r.old))
		if index == -1 {
			// 変換対象がなければ、srcをそのままコピー
			n := copy(dst[nDst:], src)
			nDst += n
			nSrc += n
			return nDst, nSrc, nil
		}

		// 変換対象より前の部分をコピー
		n := copy(dst[nDst:], src[:index])
		nDst += n
		nSrc += index

		// 変換部分をnewに置き換える
		n = copy(dst[nDst:], r.new)
		nDst += n

		// 処理済み部分をスキップ
		src = src[index+len(r.old):]
		nSrc += len(r.old)
	}

	return nDst, nSrc, nil
}


func NewReplacer(old, new []byte) *Replacer {
	return &Replacer{old: old, new: new}
}

func main() {
	// 置換対象のテキスト
	input := "郷に入っては郷に従え"

	t := NewReplacer([]byte("郷"), []byte("Go"))

	// transform.NewWriterを使って、書き込んだやつを標準出力に出す
	w := transform.NewWriter(os.Stdout, t)

	fmt.Println("変換前:",input)
	fmt.Print("変換後: ")
	io.Copy(w, strings.NewReader(input))
}