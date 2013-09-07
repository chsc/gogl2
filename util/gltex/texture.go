package glutil

import (
	"image"
	"unsafe"
)

const (
	glRGB = 0
	glUNSIGNED_BYTE = 0
	glRGBA = 0
)

func ImageFromPixelData(format, type_ int32, pixels uintptr, w, h uint32) (img image.Image, err error) {
	if format == glRGB && type_ == glUNSIGNED_BYTE {
		dest := image.NewNRGBA(image.Rect(0, 0, int(w), int(h)))
		src := (*[0xfffffff - 1]byte)((unsafe.Pointer)(pixels))[0:int(w*h*3)]
		srcOffset := 0
		destOffset := len(dest.Pix) - dest.Stride
		// flip image
		for srcOffset < len(src) {
			for i := destOffset; i < destOffset+dest.Stride; i += 4 {
				dest.Pix[i+0] = src[srcOffset+0]
				dest.Pix[i+1] = src[srcOffset+1]
				dest.Pix[i+2] = src[srcOffset+2]
				dest.Pix[i+3] = 255
				srcOffset += 3
			}
			destOffset -= dest.Stride
		}
		img = dest
	} else if format == glRGBA && type_ == glUNSIGNED_BYTE {
		dest := image.NewNRGBA(image.Rect(0, 0, int(w), int(h)))
		src := (*[0xfffffff - 1]byte)(pixels)[0 : w*h*4]
		destOffset := len(dest.Pix) - dest.Stride
		// flip image
		for srcOffset := 0; srcOffset < len(src); srcOffset += dest.Stride {
			copy(dest.Pix[destOffset:destOffset+dest.Stride], src[srcOffset:srcOffset+dest.Stride])
			destOffset -= dest.Stride
		}
		img = dest
	} else { //TODO: add more
		err = errors.New("image format not supported")
	}
	return
}
func PixelDataFromImage(img image.Image) (internalFormat int32, format, type_ uint32, pixels uintptr, width, height uint32, err error) {
	var pixelSize int = 0
	var pix []uint8 = nil
	var srcStride int = 0
	switch i := img.(type) {
	case *image.Alpha:
		internalFormat = gl.ALPHA8
		format = glALPHA
		type_ = glUNSIGNED_BYTE
		pixelSize = 1
		pix = i.Pix
		srcStride = i.Stride
	case *image.Alpha16:
		internalFormat = gl.ALPHA16
		format = glALPHA
		type_ = glUNSIGNED_SHORT
		pixelSize = 2
		pix = i.Pix
		srcStride = i.Stride
	case *image.Gray:
		internalFormat = gl.R8 // luminance?
		format = glRED
		type_ = gl.UNSIGNED_BYTE
		pixelSize = 1
		pix = i.Pix
		srcStride = i.Stride
	case *image.Gray16:
		internalFormat = gl.R16
		format = glRED
		type_ = glUNSIGNED_SHORT
		pixelSize = 2
		pix = i.Pix
		srcStride = i.Stride
	case *image.NRGBA:
		internalFormat = gl.RGBA8
		format = glRGBA
		type_ = glUNSIGNED_BYTE
		pixelSize = 4
		pix = i.Pix
		srcStride = i.Stride
	case *image.NRGBA64:
		internalFormat = glRGBA16
		format = glRGBA
		type_ = glUNSIGNED_SHORT
		pixelSize = 8
		pix = i.Pix
		srcStride = i.Stride
	case *image.RGBA:
		internalFormat = glRGBA8
		format = glRGBA
		type_ = glUNSIGNED_BYTE
		pixelSize = 4
		pix = i.Pix
		srcStride = i.Stride
		fmt.Println("nanananananana Batman!")
	case *image.RGBA64:
		internalFormat = glRGBA16
		format = glRGBA
		type_ = glUNSIGNED_SHORT
		pixelSize = 8
		pix = i.Pix
		srcStride = i.Stride
	default: // TODO: add more
		err = errors.New("image type not supported")
	}
	// flip image to GL format: first pixel is lower left corner
	width, height = gl.Sizei(img.Bounds().Dx()), gl.Sizei(img.Bounds().Dy())
	data := make([]byte, int(width*height)*pixelSize)
	start := 0
	end := len(pix)
	fmt.Println(start, end, len(pix))
	destOffset := len(data) - srcStride
	for srcOffset := start; srcOffset < end; srcOffset += srcStride {
		copy(data[destOffset:destOffset+srcStride], pix[srcOffset:srcOffset+srcStride])
		destOffset -= srcStride
	}
	pixels = gl.Pointer(&data[0])
	return
}

