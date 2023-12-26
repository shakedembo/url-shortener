package utils

import (
	"hash/fnv"
	"log"

	"github.com/catinello/base62"
)

type HashProvider[T any] interface {
	Get(key T) string
}

type ConvertToBytesFunc[T any] func(T) []byte

type Fnv32aHashProvider[T any] struct {
	logger  *log.Logger
	convert ConvertToBytesFunc[T]
}

func NewFnv32aHashProvider[T any](convert ConvertToBytesFunc[T], logger *log.Logger) HashProvider[T] {
	return &Fnv32aHashProvider[T]{
		logger:  logger,
		convert: convert,
	}
}

func (c Fnv32aHashProvider[T]) Get(keyT T) string {
	key := c.convert(keyT)
	hash := fnv.New32a()

	_, err := hash.Write(key)
	if err != nil {
		c.logger.Printf("Error occurred trying to hash the key `%v`", keyT)
		return "0"
	}
	num := hash.Sum32()
	return base62.Encode(int(num))
}
