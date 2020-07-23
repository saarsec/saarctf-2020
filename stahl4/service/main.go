package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"io/ioutil"
	"math/big"
	"os"
)

const logistic = 1
const production = 2
const port = "21485"
const privPath = "./data/priv.key"

var privK *rsa.PrivateKey

func main() {
	// generate key files if they do not exist
	_, privErr := os.Stat(privPath)
	if privErr != nil {
		var err error
		privK, err = rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			panic(err)
		}
		privKey := x509.MarshalPKCS1PrivateKey(privK)
		err = ioutil.WriteFile(privPath, privKey, 0644)
		if err != nil {
			panic(err)
		}
		println("generated new keyfile")
	} else {
		tmp, err := ioutil.ReadFile(privPath)
		if err != nil {
			panic(err)
		}
		privK, err = x509.ParsePKCS1PrivateKey(tmp)
		if err != nil {
			panic(err)
		}
		println("loaded existing keyfile")
	}

	//cmd := exec.Command("ls")
	//cmd.Stdin = strings.NewReader("-lah")
	//var out bytes.Buffer
	//cmd.Stdout = &out
	//err := cmd.Run()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("result:\n %s", out.String())

	// create random integer between 0 and 1
	rnd, err := rand.Int(rand.Reader, big.NewInt(2))
	if err != nil {
		panic(err)
	} else {
		typ := uint8(rnd.Int64() + 1)
		println("Starting machine with typ: ", typ)
		startMachine(typ)
	}
}
