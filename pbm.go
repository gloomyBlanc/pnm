package pnm

import (
	"errors"
	"image"
	"image/color"
)

var (
	errBadPBMSample = errors.New("pnm: PBM画像のサンプル値が不正です")
)

// Reader用
func (d *pnmDecoder) pbmReadRaster() (image.Image, error) {
	var (
		i, j, k int
		b       byte
		err     error
	)

	img := image.NewGray(image.Rect(0, 0, d.h.width, d.h.height))
	for i = 0; i < d.h.height; i++ {
		for j = 0; j < d.h.width; {
			b, err = d.reader.ReadByte()
			if err != nil {
				return nil, errBadPBMSample
			}
			switch d.h.magicNumber {
			case "P1":
				if !isWhiteSpece(b) {
					img.SetGray(j, i, color.Gray{255 * (b - '0')})
					j += 1
				}
			case "P4":
				for k = 0; k < 8; k++ {
					img.SetGray(j+k, i, color.Gray{255 * ((b >> (7 - k)) & 1)})
				}
				j += 8
			}
		}
	}
	return img, nil
}

// Writer用
func (e *pnmEncoder) pbmWriteRasterPlain(img image.Image) error {
	var (
		i, j int
		y    uint32
	)

	rect := img.Bounds()
	for i = rect.Min.Y; i < rect.Max.Y; i++ {
		for j = rect.Min.X; j < rect.Max.X; j++ {
			y, _, _, _ = rect.At(j, i).RGBA()
			if y == 0 {
				e.writer.WriteRune('0')
			} else {
				e.writer.WriteRune('1')
			}
			if j == rect.Max.X-1 {
				e.writer.WriteRune('\n')
			} else {
				e.writer.WriteRune(' ')
			}
		}
	}
	return nil
}
func (e *pnmEncoder) pbmWriteRasterBinary(img image.Image) error {
	var (
		i, j, k int
		y       uint32
		b       byte
	)

	rect := img.Bounds()
	for i = rect.Min.Y; i < rect.Max.Y; i++ {
		for j = rect.Min.X; j < rect.Max.X; j += 8 {
			b = 0x00
			for k = 0; k < 8; k++ {
				if (j + k) < rect.Max.X {
					y, _, _, _ = rect.At(j+k, i).RGBA()
				} else {
					y = 0
				}
				b = b << 1
				if y != 0 {
					b |= 0x01
				}
			}
			e.writer.WriteByte(b)
		}
	}
	return nil
}
