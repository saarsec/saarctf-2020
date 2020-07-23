package network

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"
)

var maxChallenges = 1
var allowedTimeDifference = 1 * time.Minute

type ClientHandshake struct {
	PubKey rsa.PublicKey `json:"pub,string"`
	Nonce  string        `json:"nonce,string"`
	Expiry time.Time     `json:"expiry,string"`
	Target string        `json:"target"`
	Hmac   []byte        `json:"hmac,string"`
}

type ServerHandshake struct {
	Type int `json:"type,string"`
}

type ServerChallenge struct {
	Nonce string `json:"nonce,string"`
}

func (clientHandshake *ClientHandshake) MarshalJSON() ([]byte, error) {
	return json.Marshal(clientHandshake.buildMap())
}

func (clientHandshake *ClientHandshake) buildMap() map[string]string {
	return map[string]string{
		"pub":    base64.URLEncoding.EncodeToString(x509.MarshalPKCS1PublicKey(&clientHandshake.PubKey)),
		"nonce":  clientHandshake.Nonce,
		"expiry": clientHandshake.Expiry.Format(time.RFC3339),
		"target": clientHandshake.Target,
		"hmac":   string(clientHandshake.Hmac),
	}
}

func (clientHandshake *ClientHandshake) fromMap(cmap map[string]string) error {
	tmp, err := base64.URLEncoding.DecodeString(cmap["pub"])
	if err != nil {
		return err
	}
	pubkey, err := x509.ParsePKCS1PublicKey(tmp)
	if err != nil {
		tmp, err := x509.ParsePKIXPublicKey(tmp)
		if err != nil {
			return err
		}
		pubkey = tmp.(*rsa.PublicKey)
	}
	clientHandshake.PubKey = *pubkey
	clientHandshake.Nonce = cmap["nonce"]
	t, err := time.Parse(time.RFC3339, cmap["expiry"])
	if err != nil {
		return err
	}
	clientHandshake.Expiry = t
	clientHandshake.Target = cmap["target"]
	clientHandshake.Hmac = []byte(cmap["hmac"])
	return nil
}

func (clientHandshake *ClientHandshake) sign(key *rsa.PrivateKey) error {
	tmpN, err := base64.URLEncoding.DecodeString(clientHandshake.Nonce)
	if err != nil {
		return err
	}
	// check if we have gameserver <-> team our minimum expected nonce length for hardcoded length later on
	if len(tmpN) < 8 {
		return errors.New("nonce to short")
	}
	clientHandshake.PubKey = key.PublicKey
	//t, err := time.Parse(time.RFC3339, clientHandshake.Expiry.Format(time.RFC3339))
	tStr := clientHandshake.Expiry.Format("2006-01-02T15:04:05-0700")
	for len(tStr) < 24 {
		tStr += "\x00"
	}
	// len(t.String) = 24; len(Target) can be 7 to 15; len(tmpN) = 8 (should be)
	message := []byte(string(tmpN) + clientHandshake.Target + tStr)
	log.Print("Client: message to calculate hash of ", string(message))
	// TODO vuln3
	hashed := hash(message[0:(8 + len(clientHandshake.Target) + 24)])

	// log.Print("Client: calculated hash of message ", string(hashed))
	tmp, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA512, hashed)
	if err != nil {
		return err
	}
	clientHandshake.Hmac = []byte(base64.URLEncoding.EncodeToString(tmp))
	return nil
}

func (clientHandshake *ClientHandshake) verify(server *Server, nonce string) bool {
	cmap := clientHandshake.buildMap()
	if clientHandshake.Target != server.ip {
		log.Printf("expected: %s - got: %s", server.ip, clientHandshake.Target)
		return false
	}
	t := time.Now()
	//fmt.Printf("our timestamp: %s\n", t)
	//fmt.Printf("client timestamp: %s\n", clientHandshake.Expiry)
	if clientHandshake.Expiry.Before(t) {
		log.Printf("invalid timestamp, now %s got %s", t, clientHandshake.Expiry)
		return false
	}
	tmpN, err := base64.URLEncoding.DecodeString(nonce)
	if err != nil {
		log.Printf("couldnt decode nonce")
		return false
	}
	// check if we have our minimum expected nonce length for hardcoded length later on
	if len(tmpN) < 8 {
		return false
	}
	if !strings.Contains(clientHandshake.Nonce, nonce) {
		log.Printf("expected different nonce, wanted " + nonce + " got " + clientHandshake.Nonce)
		return false
	}
	hmac, err := base64.URLEncoding.DecodeString(cmap["hmac"])
	if err != nil {
		log.Printf("couldnt decode hmac")
		return false
	}
	//tStr := clientHandshake.Expiry.String()
	tStr := clientHandshake.Expiry.Format("2006-01-02T15:04:05-0700")
	for len(tStr) < 24 {
		tStr += "\x00"
	}
	message := []byte(string(tmpN) + clientHandshake.Target + tStr)
	// log.Print("Server: message to calculate hash of ", hex.Dump(message[0:(8+len(clientHandshake.Target)+37)]))
	// TODO vuln3
	hashed := hash(message[0:(8 + len(clientHandshake.Target) + 24)])
	log.Print(string(message[0:(8 + len(clientHandshake.Target) + 24)]))
	log.Print("Server: calculated hash of message, got ", base64.URLEncoding.EncodeToString(hashed))
	err = rsa.VerifyPKCS1v15(&clientHandshake.PubKey, crypto.SHA512, hashed, []byte(hmac))
	if err != nil {
		log.Printf("hmac did not match")
		return false
	}
	return true
}

func newServerChallenge() (*ServerChallenge, error) {
	nonce, err := newServerNonce()
	if err != nil {
		return nil, err
	}
	return &ServerChallenge{
		Nonce: base64.URLEncoding.EncodeToString(nonce),
	}, nil
}

func newServerNonce() ([]byte, error) {
	nonce := make([]byte, 8)
	_, err := rand.Read(nonce)
	if err != nil {
		return nil, err
	}
	return nonce, nil
}

func newClientHandshake(target string, nonce string) (*ClientHandshake, error) {
	return &ClientHandshake{
		PubKey: rsa.PublicKey{},
		Nonce:  nonce,
		Expiry: time.Now().Add(allowedTimeDifference),
		Target: target,
		Hmac:   nil,
	}, nil
}

func hash(data []byte) []byte {
	hash := sha512.New()
	hash.Write(data)
	return hash.Sum(nil)
}
