package main

import (
	"crypto/rand"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"io/ioutil"
	"math/big"
	"strconv"
)

var window fyne.Window
var entryInput binding.String
var entryOutput binding.String
var entryPublicKey binding.String
var entryPrivateKey binding.String
var entryAdditional binding.String
var frogKeysGlobal *FrogKeys

func main() {
	a := app.New()
	window = a.NewWindow("Cryptographer")

	entryInput = binding.NewString()
	entryOutput = binding.NewString()
	entryPublicKey = binding.NewString()
	entryPrivateKey = binding.NewString()
	entryAdditional = binding.NewString()
	window.SetContent(
		container.NewGridWithRows(
			2,
			widget.NewForm(
				widget.NewFormItem("Input:", widget.NewEntryWithData(entryInput)),
				widget.NewFormItem("Output:", widget.NewEntryWithData(entryOutput)),
				widget.NewFormItem("Public key:", widget.NewEntryWithData(entryPublicKey)),
				widget.NewFormItem("Private key:", widget.NewEntryWithData(entryPrivateKey)),
				widget.NewFormItem("Additional value:", widget.NewEntryWithData(entryAdditional)),
			),
			container.NewGridWithColumns(
				3,
				container.NewGridWithRows(
					3,
					widget.NewButton("Generate LUC keys", func() {
						str, err := entryInput.Get()
						if err != nil || len(str) == 0 {
							dialog.NewInformation("Ошибка", "Некорректный ввод", window).Show()
							return
						}
						bigint, success := new(big.Int).SetString(str, 10)
						if !success {
							dialog.NewInformation("Ошибка", "Некорректный ввод", window).Show()
							return
						}
						publicKeyLUC, privateKeyLUC := generateKeys(bigint.Bytes())
						entryPublicKey.Set(new(big.Int).SetBytes(publicKeyLUC.keyValue).Text(10))
						entryPrivateKey.Set(new(big.Int).SetBytes(privateKeyLUC.keyValue).Text(10))
						entryAdditional.Set(new(big.Int).SetBytes(publicKeyLUC.mulValue).Text(10))
					}),
					widget.NewButton("Encrypt LUC", func() {
						publicKeyLUCStr, err := entryPublicKey.Get()
						if err != nil || len(publicKeyLUCStr) == 0 {
							dialog.NewInformation("Ошибка", "Некорректный ввод", window).Show()
							return
						}
						publicKeyLUC, success := new(big.Int).SetString(publicKeyLUCStr, 10)
						if !success {
							dialog.NewInformation("Ошибка", "Некорректный ввод", window).Show()
							return
						}
						MulValueStr, err := entryAdditional.Get()
						if err != nil || len(MulValueStr) == 0 {
							dialog.NewInformation("Ошибка", "Некорректный ввод", window).Show()
							return
						}
						MulValue, success := new(big.Int).SetString(MulValueStr, 10)
						if !success {
							dialog.NewInformation("Ошибка", "Некорректный ввод", window).Show()
							return
						}
						dataStr, err := entryInput.Get()
						if err != nil || len(dataStr) == 0 {
							dialog.NewInformation("Ошибка", "Некорректный ввод", window).Show()
							return
						}
						data, success := new(big.Int).SetString(dataStr, 10)
						if !success {
							dialog.NewInformation("Ошибка", "Некорректный ввод", window).Show()
							return
						}
						entryOutput.Set(new(big.Int).SetBytes(encDec(LucKey{publicKeyLUC.Bytes(), MulValue.Bytes()}, data.Bytes())).Text(10))
					}),
					widget.NewButton("Decrypt LUC", func() {
						privateKeyLUCStr, err := entryPrivateKey.Get()
						if err != nil || len(privateKeyLUCStr) == 0 {
							dialog.NewInformation("Ошибка", "Некорректный ввод", window).Show()
							return
						}
						privateKeyLUC, success := new(big.Int).SetString(privateKeyLUCStr, 10)
						if !success {
							dialog.NewInformation("Ошибка", "Некорректный ввод", window).Show()
							return
						}
						MulValueStr, err := entryAdditional.Get()
						if err != nil || len(MulValueStr) == 0 {
							dialog.NewInformation("Ошибка", "Некорректный ввод", window).Show()
							return
						}
						MulValue, success := new(big.Int).SetString(MulValueStr, 10)
						if !success {
							dialog.NewInformation("Ошибка", "Некорректный ввод", window).Show()
							return
						}
						dataStr, err := entryInput.Get()
						if err != nil || len(dataStr) == 0 {
							dialog.NewInformation("Ошибка", "Некорректный ввод", window).Show()
							return
						}
						data, success := new(big.Int).SetString(dataStr, 10)
						if !success {
							dialog.NewInformation("Ошибка", "Некорректный ввод", window).Show()
							return
						}
						entryOutput.Set(new(big.Int).SetBytes(encDec(LucKey{privateKeyLUC.Bytes(), MulValue.Bytes()}, data.Bytes())).Text(10))
					}),
				),
				container.NewGridWithRows(
					3,
					container.NewGridWithColumns(
						2,
						widget.NewButton("Generate Frog key", func() {
							keySizeStr, err := entryInput.Get()
							if err != nil || len(keySizeStr) == 0 {
								dialog.NewInformation("Ошибка", "Некорректный ввод", window).Show()
								return
							}
							keyValue, err := strconv.Atoi(keySizeStr)
							if err != nil {
								dialog.NewInformation("Ошибка", "Некорректный ввод", window).Show()
								return
							}
							if keyValue < 5 || keyValue > 125 {
								dialog.NewInformation("Ошибка", "Некорректный размер ключа", window).Show()
								return
							}
							key := make([]byte, keyValue)
							rand.Read(key)
							entryPublicKey.Set(new(big.Int).SetBytes(key).Text(10))
						}),
						widget.NewButton("Generate Frog Round Keys", func() {
							keyStr, err := entryPublicKey.Get()
							if err != nil || len(keyStr) == 0 {
								dialog.NewInformation("Ошибка", "Некорректный ввод", window).Show()
								return
							}
							key, success := new(big.Int).SetString(keyStr, 10)
							if !success {
								dialog.NewInformation("Ошибка", "Некорректный ввод", window).Show()
								return
							}
							frogKeysGlobal = newFrog(key.Bytes())
						}),
					),
					widget.NewButton("Encrypt", func() {
						if frogKeysGlobal == nil {
							dialog.NewInformation("Ошибка", "Раундовые ключи ещё не сгенерированы", window).Show()
							return
						}
						dataStr, err := entryInput.Get()
						if err != nil || len(dataStr) == 0 {
							dialog.NewInformation("Ошибка", "Некорректный ввод", window).Show()
							return
						}
						data, success := new(big.Int).SetString(dataStr, 10)
						if !success {
							dialog.NewInformation("Ошибка", "Некорректный ввод", window).Show()
							return
						}

						entryOutput.Set(new(big.Int).SetBytes(Encrypt(frogKeysGlobal, addPadding(data.Bytes(), 16))).Text(10))
					}),
					widget.NewButton("Decrypt", func() {
						if frogKeysGlobal == nil {
							dialog.NewInformation("Ошибка", "Раундовые ключи ещё не сгенерированы", window).Show()
							return
						}
						dataStr, err := entryInput.Get()
						if err != nil || len(dataStr) == 0 {
							dialog.NewInformation("Ошибка", "Некорректный ввод", window).Show()
							return
						}
						data, success := new(big.Int).SetString(dataStr, 10)
						if !success {
							dialog.NewInformation("Ошибка", "Некорректный ввод", window).Show()
							return
						}
						if len(dataStr) < 16 {
							dialog.NewInformation("Ошибка", "Некорректный ввод", window).Show()
							return
						}
						entryOutput.Set(new(big.Int).SetBytes(removePadding(Decrypt(frogKeysGlobal, data.Bytes()))).Text(10))
					}),
				),
				container.NewGridWithRows(
					2,
					widget.NewButton("Encrypt File", func() {
						if frogKeysGlobal == nil {
							dialog.NewInformation("Ошибка", "Раундовые ключи ещё не сгенерированы", window).Show()
							return
						}
						FilenameStr, err := entryInput.Get()
						if err != nil || len(FilenameStr) == 0 {
							dialog.NewInformation("Ошибка", "Некорректный ввод", window).Show()
							return
						}
						FilenameOutStr, err := entryOutput.Get()
						if err != nil || len(FilenameOutStr) == 0 {
							dialog.NewInformation("Ошибка", "Некорректный ввод", window).Show()
							return
						}
						file, err := ioutil.ReadFile(FilenameStr)
						if err != nil {
							dialog.NewInformation("Ошибка", "Некорректное имя файла", window).Show()
							return
						}
						file = addPadding(file, 16)
						enc := Encrypt(frogKeysGlobal, file)

						ioutil.WriteFile(FilenameOutStr, enc, 0777)
					}),
					widget.NewButton("Decrypt File", func() {
						if frogKeysGlobal == nil {
							dialog.NewInformation("Ошибка", "Раундовые ключи ещё не сгенерированы", window).Show()
							return
						}
						FilenameStr, err := entryInput.Get()
						if err != nil || len(FilenameStr) == 0 {
							dialog.NewInformation("Ошибка", "Некорректный ввод", window).Show()
							return
						}
						FilenameOutStr, err := entryOutput.Get()
						if err != nil || len(FilenameOutStr) == 0 {
							dialog.NewInformation("Ошибка", "Некорректный ввод", window).Show()
							return
						}
						file, err := ioutil.ReadFile(FilenameStr)
						if err != nil {
							dialog.NewInformation("Ошибка", "Некорректное имя файла", window).Show()
							return
						}
						dec := Decrypt(frogKeysGlobal, file)
						dec = removePadding(dec)

						ioutil.WriteFile(FilenameOutStr, dec, 0777)
					}),
				),
			),
		),
	)
	window.ShowAndRun()

}
