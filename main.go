package main

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

	fk := newFrogKeys(32)
	println(fk.key)
}
