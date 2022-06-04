package main

import (
	"crypto/rand"
	"math/big"
)

const (
	EncryptOrder = iota
	DecryptOrder
)

type FrogKeys struct {
	key              []byte
	encryptRoundKeys [][][]byte
	decryptRoundKeys [][][]byte
}

func newFrogKeys(size int) *FrogKeys {
	frogKeys := new(FrogKeys)
	frogKeys.key = make([]byte, size)
	rand.Read(frogKeys.key)
	frogKeys.encryptRoundKeys = generateKey(frogKeys.key, EncryptOrder)
	frogKeys.decryptRoundKeys = generateKey(frogKeys.key, DecryptOrder)
	return frogKeys
}

var MasterKey = []byte{
	113, 21, 232, 18, 113, 92, 63, 157, 124, 193, 166, 197, 126, 56, 229, 229,
	156, 162, 54, 17, 230, 89, 189, 87, 169, 0, 81, 204, 8, 70, 203, 225,
	160, 59, 167, 189, 100, 157, 84, 11, 7, 130, 29, 51, 32, 45, 135, 237,
	139, 33, 17, 221, 24, 50, 89, 74, 21, 205, 191, 242, 84, 53, 3, 230,
	231, 118, 15, 15, 107, 4, 21, 34, 3, 156, 57, 66, 93, 255, 191, 3,
	85, 135, 205, 200, 185, 204, 52, 37, 35, 24, 68, 185, 201, 10, 224, 234,
	7, 120, 201, 115, 216, 103, 57, 255, 93, 110, 42, 249, 68, 14, 29, 55,
	128, 84, 37, 152, 221, 137, 39, 11, 252, 50, 144, 35, 178, 190, 43, 162,
	103, 249, 109, 8, 235, 33, 158, 111, 252, 205, 169, 54, 10, 20, 221, 201,
	178, 224, 89, 184, 182, 65, 201, 10, 60, 6, 191, 174, 79, 98, 26, 160,
	252, 51, 63, 79, 6, 102, 123, 173, 49, 3, 110, 233, 90, 158, 228, 210,
	209, 237, 30, 95, 28, 179, 204, 220, 72, 163, 77, 166, 192, 98, 165, 25,
	145, 162, 91, 212, 41, 230, 110, 6, 107, 187, 127, 38, 82, 98, 30, 67,
	225, 80, 208, 134, 60, 250, 153, 87, 148, 60, 66, 165, 72, 29, 165, 82,
	211, 207, 0, 177, 206, 13, 6, 14, 92, 248, 60, 201, 132, 95, 35, 215,
	118, 177, 121, 180, 27, 83, 131, 26, 39, 46, 12,
}

func generateKey(key []byte, order int) [][][]byte {
	expandedKey := expandKey(key, 2304)
	expandedMasterKey := expandKey(MasterKey, 2304)
	expandedKey = new(big.Int).Xor(new(big.Int).SetBytes(expandedKey), new(big.Int).SetBytes(expandedMasterKey)).Bytes()
	preliminaryExpandedKey := FormatExpandedKey(key, EncryptOrder)
	//println(expandedKey, expandedMasterKey)

	return make([][][]byte, 1)
}

func expandKey(key []byte, newSize int) []byte {
	res := make([]byte, newSize)
	for i := range res {
		res[i] = key[i%len(key)]
	}
	return res
}

//func FormatExpandedKey(key []byte, order int) []byte {
//
//}
