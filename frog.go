package main

import (
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
	unitSize         int
}

func newFrog(key []byte) *FrogKeys {
	frogKeys := new(FrogKeys)
	frogKeys.unitSize = 16
	frogKeys.key = key
	frogKeys.encryptRoundKeys = generateKey(frogKeys.key, EncryptOrder, frogKeys.unitSize)
	frogKeys.decryptRoundKeys = generateKey(frogKeys.key, DecryptOrder, frogKeys.unitSize)
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

func Encrypt(keys *FrogKeys, data []byte) []byte {
	res := make([]byte, len(data))
	copy(res, data)

	unitCount := len(res) / keys.unitSize
	for roundNum := 0; roundNum < 8; roundNum++ {
		for j := 0; j < unitCount; j++ {
			for i := 0; i < keys.unitSize; i++ {
				res[i+j*keys.unitSize] ^= keys.encryptRoundKeys[roundNum][0][i]
				res[i+j*keys.unitSize] = keys.encryptRoundKeys[roundNum][1][res[i+j*keys.unitSize]]
				if i < keys.unitSize-1 {
					res[i+1+j*keys.unitSize] ^= res[i+j*keys.unitSize]
				}
				id := keys.encryptRoundKeys[roundNum][2][i]
				res[int(id)+j*keys.unitSize] ^= res[i+j*keys.unitSize]
			}
		}
	}
	return res
}

func Decrypt(keys *FrogKeys, data []byte) []byte {
	res := make([]byte, len(data))
	copy(res, data)

	unitCount := len(res) / keys.unitSize
	for roundNum := 7; roundNum >= 0; roundNum-- {
		for j := 0; j < unitCount; j++ {
			for i := keys.unitSize - 1; i >= 0; i-- {
				id := keys.decryptRoundKeys[roundNum][2][i]
				res[int(id)+j*keys.unitSize] ^= res[i+j*keys.unitSize]
				if i < keys.unitSize-1 {
					res[i+1+j*keys.unitSize] ^= res[i+j*keys.unitSize]
				}
				res[i+j*keys.unitSize] = keys.decryptRoundKeys[roundNum][1][res[i+j*keys.unitSize]]
				res[i+j*keys.unitSize] ^= keys.decryptRoundKeys[roundNum][0][i]
			}
		}
	}
	return res
}

func generateKey(key []byte, order int, unitSize int) [][][]byte {
	expandedKey := expandKey(key, 2304)
	expandedMasterKey := expandKey(MasterKey, 2304)
	expandedKey = new(big.Int).Xor(new(big.Int).SetBytes(expandedKey), new(big.Int).SetBytes(expandedMasterKey)).Bytes()
	preliminaryExpandedKey := FormatExpandedKey(expandedKey, EncryptOrder)
	IV := make([]byte, unitSize)
	copy(IV, expandedKey[:unitSize])
	IV[0] ^= byte(len(key))
	res := FormatEmptyArray(preliminaryExpandedKey, IV, unitSize)
	return FormatExpandedKey(res, order)
}

func expandKey(key []byte, newSize int) []byte {
	res := make([]byte, newSize)
	for i := range res {
		res[i] = key[i%len(key)]
	}
	return res
}

func FormatKey(key []byte) {
	U := make([]byte, len(key))
	for i := range U {
		U[i] = byte(i)
	}

	prevId := 0
	var currId int
	for i := range key {
		currId = (prevId + int(key[i])) % len(U)
		prevId = currId
		key[i] = U[currId]
		U = append(U[:currId], U[currId+1:]...)
	}
}

func ReverseKey(key []byte) []byte {
	res := make([]byte, len(key))
	for i := range key {
		res[key[i]] = byte(i)
	}
	return res
}

func connectElements(permutation []byte) {
	isConnected := make([]bool, len(permutation))
	for i := range isConnected {
		isConnected[i] = false
	}

	id := 0
	for {
		isConnected[id] = true
		if isConnected[permutation[id]] {
			nextNotConnected := -1
			for i := range isConnected {
				if isConnected[i] == false {
					nextNotConnected = i
					break
				}
			}
			if nextNotConnected == -1 {
				permutation[id] = 0
				break
			} else {
				permutation[id] = byte(nextNotConnected)
			}
		}
		id = int(permutation[id])
	}
}

func FormatExpandedKey(key []byte, order int) [][][]byte {
	res := make([][][]byte, 8)
	for i := 0; i < 8; i++ {
		keyComponent1 := make([]byte, 16)
		keyComponent2 := make([]byte, 256)
		keyComponent3 := make([]byte, 16)
		currentId := i * 288
		copy(keyComponent1, key[currentId:currentId+16])
		copy(keyComponent2, key[currentId+16:currentId+272])
		copy(keyComponent3, key[currentId+272:currentId+288])

		FormatKey(keyComponent2)

		if order == DecryptOrder {
			keyComponent2 = ReverseKey(keyComponent2)
		}
		FormatKey(keyComponent3)
		connectElements(keyComponent3)
		for j := 0; j < 16; j++ {
			if int(keyComponent3[j]) == j+1 {
				keyComponent3[j] = byte((j + 2) % 16)
			}
		}
		res[i] = [][]byte{keyComponent1, keyComponent2, keyComponent3}
	}
	return res
}

func FormatEmptyArray(key [][][]byte, IV []byte, unitSize int) []byte {
	unitCount := 2304 / unitSize
	buf := make([]byte, unitSize)
	res := make([]byte, 2304)
	for i := 0; i < unitCount; i++ {
		EncryptUnit(buf, IV, key, 0, res, i*unitSize, unitSize)
	}
	return res
}

func EncryptUnit(buf []byte, IV []byte, roundKeys [][][]byte, iShift int, res []byte, oShift int, unitSize int) {
	copy(buf[iShift:iShift+unitSize], res[oShift:oShift+unitSize])

	for i := 0; i < unitSize; i++ {
		res[i] ^= IV[i]
	}

	for roundNum := 0; roundNum < 8; roundNum++ {
		for i := 0; i < unitSize; i++ {
			res[oShift+i] ^= roundKeys[roundNum][0][i]
			res[oShift+i] = roundKeys[roundNum][1][res[oShift+i]]
			if i < unitSize-1 {
				res[oShift+i+1] ^= res[oShift+i]
			}
			id := roundKeys[roundNum][2][i]
			res[oShift+int(id)] ^= res[oShift+i]
		}
	}
}

func addPadding(data []byte, size int) []byte {
	if len(data)%size == 0 {
		return data
	}
	count := size - len(data)%size
	padding := make([]byte, count)
	for i := range padding {
		padding[i] = byte(count)
	}
	return append(data, padding...)
}

func removePadding(data []byte) []byte {
	for i := len(data) - 1; i >= 0; i-- {
		if i == len(data)-int(data[i]) {
			for j := i; j < len(data); j++ {
				if data[j] != data[i] {
					return data
				}
			}
			data = data[:i]
			break
		}
	}
	return data
}
