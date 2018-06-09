package main

import (
	"image"
	"github.com/corona10/goimagehash"
	"strconv"
	"fmt"
)

type ImageHash struct {
	Dhash string `json:"dhash"`
	Ahash string `json:"ahash"`
}

func HashImage(i image.Image) ImageHash {
	ahash, _ := goimagehash.AverageHash(i)
	dhash, _ := goimagehash.DifferenceHash(i)
	ahashString := fmt.Sprintf("%064v", strconv.FormatUint(ahash.GetHash(), 2))
	dhashString := fmt.Sprintf("%064v", strconv.FormatUint(dhash.GetHash(), 2))
	return ImageHash{Dhash: dhashString, Ahash: ahashString}
}
