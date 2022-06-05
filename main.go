package main

import (
	"crypto/rand"
	"math/big"
)

func main() {
	//start := time.Now()
	//message, _ := rand.Int(rand.Reader, new(big.Int).Exp(big.NewInt(2), big.NewInt(1000), nil))
	//println("Message:   ", message.Text(10))
	//publicKey, privateKey := generateKeys(message.Bytes())
	//encrypted := encDec(publicKey, message.Bytes())
	//println("Encrypted: ", new(big.Int).SetBytes(encrypted).Text(10))
	//decrypted := encDec(privateKey, encrypted)
	//println("Decrypted: ", new(big.Int).SetBytes(decrypted).Text(10))
	//
	//elapsed := time.Since(start)
	//log.Printf("Execution time %s", elapsed)

	key := make([]byte, 125)
	rand.Read(key)
	fk := newFrog(key)
	message := make([]byte, 32)
	rand.Read(message)
	println("Message: ", new(big.Int).SetBytes(message).Text(10))
	message = addPadding(message, 16)
	enc := Encrypt(fk, message)
	println("Encrypted: ", new(big.Int).SetBytes(enc).Text(10))
	dec := Decrypt(fk, enc)
	dec = removePadding(dec)
	println("Decrypted: ", new(big.Int).SetBytes(dec).Text(10))
}
