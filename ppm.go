package pnm

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"strconv"
)

var (
	errBadPPMSample = errors.New("pnm: PPM画像のサンプル値が不正です")
)

// Reader用
func (d *pnmDecoder) ppmReadRaster() (image.Image, error) {
	var (
		i, j, k     int
		b           byte
		pixel       [3]int
		readBytes   []byte
		err         error
		overFF      bool
		enSampleEnd bool
	)
	overFF = (d.h.maxValue > 255)
	img := image.NewNRGBA64(image.Rect(0, 0, d.h.width, d.h.height))

	enSampleEnd = false
	for i = 0; i < d.h.height; i++ {
		for j = 0; j < d.h.width; j++ {
			for k = 0; k < 3; {
				b, err = d.reader.ReadByte()
				if err != nil {
					return nil, errBadPPMSample
				}
				switch d.h.magicNumber {
				case "P3":
					if enSampleEnd {
						if isWhiteSpece(b) {
							pixel[k], err = strconv.Atoi(string(readBytes))
							if err != nil {
								return nil, errBadPPMSample
							}
							readBytes = []byte{}
							enSampleEnd = false
							k += 1
						} else {
							readBytes = append(readBytes, b)
						}
					} else if !isWhiteSpece(b) {
						readBytes = append(readBytes, b)
						enSampleEnd = true
					}
				case "P6":
					if overFF {
						if enSampleEnd {
							pixel[k] = (pixel[k] << 8) | int(b)
							enSampleEnd = false
							k += 1
						} else {
							pixel[k] = int(b)
							enSampleEnd = true
						}
					} else {
						pixel[k] = int(b)
						k += 1
					}
				}
			}
			// pixel値の代入
			img.SetNRGBA64(j, i,
				color.NRGBA64{
					uint16(pixel[0] * 65536.0 / d.h.maxValue),
					uint16(pixel[1] * 65536.0 / d.h.maxValue),
					uint16(pixel[2] * 65536.0 / d.h.maxValue),
					0xFFFF,
				},
			)
			pixel = [3]int{}
		}
	}
	return img, nil
}

// Writer用
func (e *pnmEncoder) ppmWriteRasterPlain(img image.Image) error {
	var (
		i, j    int
		r, g, b uint32
	)

	rect := img.Bounds()
	for i = rect.Min.Y; i < rect.Max.Y; i++ {
		for j = rect.Min.X; j < rect.Max.X; j++ {
			r, g, b, _ = rect.At(j, i).RGBA()
			e.writer.WriteString(
				fmt.Sprintf("%d %d %d",
					r,
					g,
					b,
				),
			)
			if j == rect.Max.X-1 {
				e.writer.WriteRune('\n')
			} else {
				e.writer.WriteRune(' ')
			}
		}
	}
	return nil
}
func (e *pnmEncoder) ppmWriteRasterBinary(img image.Image) error {
	var (
		i, j    int
		r, g, b uint32
		overFF  bool
	)

	overFF = (e.h.maxValue > 255)
	rect := img.Bounds()
	for i = rect.Min.Y; i < rect.Max.Y; i++ {
		for j = rect.Min.X; j < rect.Max.X; j++ {
			r, g, b, _ = rect.At(j, i).RGBA()
			if overFF {
				e.writer.Write([]byte{
					byte(r >> 8), byte(r & 0xFF),
					byte(g >> 8), byte(g & 0xFF),
					byte(b >> 8), byte(b & 0xFF),
				})
			} else {
				e.writer.Write([]byte{
					byte(r & 0xFF),
					byte(g & 0xFF),
					byte(b & 0xFF),
				})
			}
		}
	}
	return nil
}
