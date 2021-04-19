package pnm

import (
	"errors"
	"image"
	"image/color"
	"strconv"
)

var (
	errBadPGMSample = errors.New("pnm: PGM画像のサンプル値が不正です")
)

// Reader用
func (d *pnmDecoder) pgmReadRaster() (image.Image, error) {
	var (
		i, j        int
		b           byte
		pixel       int
		readBytes   []byte
		err         error
		overFF      bool
		enSampleEnd bool
	)
	overFF = (d.h.maxValue > 255)
	img := image.NewGray16(image.Rect(0, 0, d.h.width, d.h.height))

	enSampleEnd = false
	for i = 0; i < d.h.height; i++ {
		for j = 0; j < d.h.width; {
			b, err = d.reader.ReadByte()
			if err != nil {
				return nil, errBadPGMSample
			}
			switch d.h.magicNumber {
			case "P2":
				if enSampleEnd {
					if isWhiteSpece(b) {
						pixel, err = strconv.Atoi(string(readBytes))
						if err != nil {
							return nil, errBadPGMSample
						}
						img.SetGray16(j, i,
							color.Gray16{uint16(pixel * 65536.0 / d.h.maxValue)},
						)
						readBytes = []byte{}
						enSampleEnd = false
						j += 1
					} else {
						readBytes = append(readBytes, b)
					}
				} else if !isWhiteSpece(b) {
					readBytes = append(readBytes, b)
					enSampleEnd = true
				}
			case "P5":
				if overFF {
					if enSampleEnd {
						pixel = (pixel << 8) | int(b)
						img.SetGray16(j, i,
							color.Gray16{uint16(pixel * 65536.0 / d.h.maxValue)},
						)
						enSampleEnd = false
						j += 1
					} else {
						pixel = int(b)
						enSampleEnd = true
					}
				} else {
					pixel = int(b)
					img.SetGray16(j, i,
						color.Gray16{uint16(pixel * 65536.0 / d.h.maxValue)},
					)
					j += 1
				}
			}
		}
	}
	return img, nil
}

// Writer用
func (e *pnmEncoder) pgmWriteRasterPlain(img image.Image) error {
	var (
		i, j int
		y    uint32
	)

	rect := img.Bounds()
	for i = rect.Min.Y; i < rect.Max.Y; i++ {
		for j = rect.Min.X; j < rect.Max.X; j++ {
			y, _, _, _ = rect.At(j, i).RGBA()
			e.writer.WriteString(strconv.Itoa(int(y)))
			if j == rect.Max.X-1 {
				e.writer.WriteRune('\n')
			} else {
				e.writer.WriteRune(' ')
			}
		}
	}
	return nil
}
func (e *pnmEncoder) pgmWriteRasterBinary(img image.Image) error {
	var (
		i, j   int
		y      uint32
		overFF bool
	)

	overFF = (e.h.maxValue > 255)
	rect := img.Bounds()
	for i = rect.Min.Y; i < rect.Max.Y; i++ {
		for j = rect.Min.X; j < rect.Max.X; j++ {
			y, _, _, _ = rect.At(j, i).RGBA()
			if overFF {
				e.writer.Write([]byte{byte(y >> 8), byte(y & 0xFF)})
			} else {
				e.writer.WriteByte(byte(y & 0xFF))
			}
		}
	}
	return nil
}
