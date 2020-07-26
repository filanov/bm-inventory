package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base32"
	"fmt"
	"io"
	"os"

	"github.com/dgrijalva/jwt-go"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/json"
)

func GenKeys(bits int) (crypto.PublicKey, crypto.PrivateKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		fmt.Printf("RSA Keys Generation error: %v\n", err)
	}
	return key.Public(), key, err
}

func GenJSJWKS(privKey crypto.PublicKey, pubKey crypto.PublicKey) ([]byte, []byte, string, error) {
	var pubJSJWKS []byte
	var privJSJWKS []byte
	var err error

	alg := "RS256"
	use := "sig"

	//Generate random kid
	b := make([]byte, 10)
	_, err = rand.Read(b)
	if err != nil {
		fmt.Printf("Kid Generation error: %v\n", err)
	}
	kid := base32.StdEncoding.EncodeToString(b)

	//  Public and private keys in JWK format
	priv := jose.JSONWebKey{Key: privKey, KeyID: kid, Algorithm: alg, Use: use}
	pub := jose.JSONWebKey{Key: pubKey, KeyID: kid, Algorithm: alg, Use: use}
	privJWKS := jose.JSONWebKeySet{Keys: []jose.JSONWebKey{priv}}
	pubJWKS := jose.JSONWebKeySet{Keys: []jose.JSONWebKey{pub}}

	privJSJWKS, err = json.Marshal(privJWKS)
	if err != nil {
		fmt.Printf("privJSJWKS Marshaling error: %v\n", err)
	}
	pubJSJWKS, err = json.Marshal(pubJWKS)
	if err != nil {
		fmt.Printf("pubJSJWKS Marshaling error: %v\n", err)
	}
	return pubJSJWKS, privJSJWKS, kid, nil
}

func main() {
	//Generate RSA Keypair
	pub, priv, _ := GenKeys(2048)

	//Generate keys in JWK format
	pubJSJWKS, privJSJWKS, kid, _ := GenJSJWKS(priv, pub)

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"account_number": "1234567",
		"is_internal":    false,
		"is_active":      true,
		"account_id":     "7654321",
		"org_id":         "1010101",
		"last_name":      "Doe",
		"type":           "User",
		"locale":         "en_US",
		"first_name":     "John",
		"email":          "jdoe123@example.com",
		"username":       "jdoe123@example.com",
		"is_org_admin":   false,
		"clientId":       "1234",
	})
	token.Header["kid"] = kid
	tokenString, err := token.SignedString(priv)

	if err != nil {
		fmt.Printf("Token Signing error: %v\n", err)
	}
	err = newFile("auth-test-pub.json", pubJSJWKS, 0444)
	if err != nil {
		fmt.Printf("Failed to write file auth-test-pub.json: %v\n", err)
	}
	err = newFile("auth-test.json", privJSJWKS, 0400)
	if err != nil {
		fmt.Printf("Failed to write file auth-test.json: %v\n", err)
	}
	err = newFile("auth-tokenString", []byte(tokenString), 0400)
	if err != nil {
		fmt.Printf("Failed to write file auth-tokenString: %v\n", err)
	}
}

func newFile(filename string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, perm)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}
