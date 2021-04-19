package pnm

import (
	"bufio"
	"errors"
	"fmt"
	"image"
	"io"
)

var (
	errBadImageSize  = errors.New("pnm: 入力画像が不正です")
	errUnsupportType = errors.New("pnm: サポートされていないカラーモードです")
)

// グレースケールの場合はPGM形式で、カラー画像の場合はPPM形式でWriterに返す
func Encode(w io.Writer, img image.Image) error {
	var e pnmEncoder
	return e.encode(w, img, "")
}

// magicに指定の形式でWriterに返す
func EncodeWithType(w io.Writer, img image.Image, magic string) error {
	var e pnmEncoder
	return e.encode(w, img, magic)
}

type pnmHeader struct {
	magicNumber   string
	width, height int
	maxValue      int
}
type pnmEncoder struct {
	writer *bufio.Writer
	h      pnmHeader
}

///
// メソッド
///
func (e *pnmEncoder) encode(w io.Writer, img image.Image, magic string) error {
	e.writer = bufio.NewWriter(w)
	err := e.setHeader(img, magic)
	if err != nil {
		return err
	}
	switch sortPNM(e.h.magicNumber) {
	case PBM:
		fmt.Fprintf(e.writer,
			"%s\n%d %d\n",
			e.h.magicNumber,
			e.h.width,
			e.h.height,
		)
		if isPlain(e.h.magicNumber) {
			return e.pbmWriteRasterPlain(img)
		} else {
			return e.pbmWriteRasterBinary(img)
		}
	case PGM:
		fmt.Fprintf(e.writer,
			"%s\n%d %d\n%d\n",
			e.h.magicNumber,
			e.h.width,
			e.h.height,
			e.h.maxValue,
		)
		if isPlain(e.h.magicNumber) {
			return e.pgmWriteRasterPlain(img)
		} else {
			return e.pgmWriteRasterBinary(img)
		}
	case PPM:
		fmt.Fprintf(e.writer,
			"%s\n%d %d\n%d\n",
			e.h.magicNumber,
			e.h.width,
			e.h.height,
			e.h.maxValue,
		)
		if isPlain(e.h.magicNumber) {
			return e.ppmWriteRasterPlain(img)
		} else {
			return e.ppmWriteRasterBinary(img)
		}
	}
	return nil
}

func (e *pnmEncoder) setHeader(img image.Image, magic string) error {
	// 画像サイズの取得
	rect := img.Bounds()
	if rect.Dx() <= 0 || rect.Dy() <= 0 {
		return errBadImageSize
	}
	// ヘッダ情報の登録
	e.h.width = rect.Dx()
	e.h.height = rect.Dy()
	switch img.(type) {
	case *image.Gray, *image.NRGBA, *image.RGBA:
		e.h.maxValue = 255
	case *image.Gray16, *image.NRGBA64, *image.RGBA64:
		e.h.maxValue = 65535
	default:
		return errUnsupportType
	}
	if magic == "" {
		switch img.(type) {
		case *image.Gray, *image.Gray16:
			e.h.magicNumber = "P5"
		case *image.NRGBA, *image.RGBA, *image.NRGBA64, *image.RGBA64:
			e.h.magicNumber = "P6"
		default:
			return errUnsupportType
		}
	} else {
		e.h.magicNumber = magic
	}
	return nil
}

//ラスター文字列がPlain形式のmagic numberか判別
func isPlain(magic string) bool {
	var (
		i          int
		plainMagic = [3]string{"P1", "P2", "P3"}
	)
	for i = 0; i < 3; i++ {
		if magic == plainMagic[i] {
			return true
		}
	}
	return false
}
