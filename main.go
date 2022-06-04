package main

import (
	"crypto/rand"
	"math/big"
)

func main() {
	message, _ := rand.Int(rand.Reader, new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil))
	println("Message: ", message.Text(10))
	publicKey, privateKey := generateKeys(message.Bytes())
	encrypted := encDec(publicKey, message.Bytes())
	println("Encrypted: ", new(big.Int).SetBytes(encrypted).Text(10))
	decrypted := encDec(privateKey, encrypted)
	println("Decrypted: ", new(big.Int).SetBytes(decrypted).Text(10))
}
