package utils

import (
	"testing"
	"fmt"
)

const
(
	test_PRIVATE_KEY = `-----BEGIN RSA PRIVATE KEY-----
MIICWwIBAAKBgQCVrHJFMKGM0pDHaXRHpcIaXs17Wo6QhwLNuhhA3dHFFJdEWZCu
yv8VVuDv5KSNnIy8ItKDb4uGMQzCvEm2gBUwJXpSciaaGg+kqY9W1OzVSTeb7cY/
NGDgoF0aRZZ0YiNdONzh8KN05jVITywh4o6vxgYiyJdox94WWxoYrkv5iQIDAQAB
AoGAbIIQdcjD1evxuh/hhO/OwH9qSLtmD7FRfwQjASPPKCm9YHfuREo2k6ngeQox
odiUzvAP3enIJQj6T1NhvUUuFhW+p0JcW4g6m5bTq8TBAaM0HIqy2zX9aVDPTEan
GHw/tRUWycxhp7AUkerHF1eDoY5E4FQJ8y+OCQSAaMEQyYECQQDF7yZL0+oW1vAk
JYGUxYHkndFWjYgksBa88hYbiyBc1TJGb3j+hgv5R+5koAbBbaFTBtOslX8R7l+4
lP7vJ9B9AkEAwZTuKM24ZLETmlLCBgXBdc/NnU11ta8nhZNKvO12Djod50kUGSY5
dSinHt8HQ/wnI1kPnGeaAGpB76uFllAG/QJAae+pS4RMEZVQSchZJkrfToC4/d4a
M6ibQt0+v9cipwzkL5aR54fO+Mhq6yhK9VO7uDg7Km+I5wvx51S3bUCd8QJAae/v
2aKjS29gk+7AQX163tdG5dPTHAdrsHznxLaLCcQiQ0VJ222Aui3yL0HMfxcJ8B04
HtbPf3Sm+ts58wV+nQJACStkg2q9k3/Msxl3B5Yz01Id3u/Jc98LJ3Qi0xU7Uord
UxCFvlEWgaVpD0yPsUbpirEj2byGoVZSwLUBgFICRA==
-----END RSA PRIVATE KEY-----
`
	test_PUBLIC_KEY = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCVrHJFMKGM0pDHaXRHpcIaXs17
Wo6QhwLNuhhA3dHFFJdEWZCuyv8VVuDv5KSNnIy8ItKDb4uGMQzCvEm2gBUwJXpS
ciaaGg+kqY9W1OzVSTeb7cY/NGDgoF0aRZZ0YiNdONzh8KN05jVITywh4o6vxgYi
yJdox94WWxoYrkv5iQIDAQAB
-----END PUBLIC KEY-----
`
	test_PRIVATE_PATH = "c_private.pem"
	test_PUBLIC_PATH  = "c_public.pem"
)

// Test load & parse public/private keys
func TestLoadRsaKey(t *testing.T) {
	privateKey, err := LoadPrivateKey(test_PRIVATE_PATH)
	if err != nil {
		t.Errorf("LoadPrivateKey fail!\n%#v\n", err)
	}
	fmt.Printf("PrivateKey:\n%#v\n", privateKey)

	publicKey, err := LoadPublicKey(test_PUBLIC_PATH)
	if err != nil {
		t.Errorf("LoadPublicKey fail!\n%#v\n", err)
	}

	fmt.Printf("PublicKey:\n%#v\n", publicKey)
}

// Test sign & verify
func TestSignVerify(t *testing.T) {
	privateKey, err := ParsePrivateKey([]byte(test_PRIVATE_KEY))
	if err != nil {
		t.Errorf("ParsePrivateKey failed!\n%#v\n", err)
	}
	fmt.Printf("ParsePrivateKey done!\n%#v\n", err)

	publicKey, err := ParsePublicKey([]byte(test_PUBLIC_KEY))
	if err != nil {
		t.Errorf("ParsePublicKey fail!\n%#v\n", err)
	}
	fmt.Printf("ParsePublicKey done!\n%#v\n", err)

	message := "message"
	msg := []byte(message)

	// Sign & Verify in raw bytes
	sig, err := Sign(privateKey, msg)
	if err != nil {
		t.Errorf("Sign failed!\n%#v\n", err)
	}
	fmt.Printf("Sign done：%v\n", sig)

	err = Verify(publicKey, msg, sig)
	if err != nil {
		t.Errorf("Verify failed!\n%#v\n", err)
	}
	fmt.Printf("Verify done \n")

	// Sign & Verify in base64
	sigB64, err := SignBase64(privateKey, msg)
	if err != nil {
		t.Errorf("SignBase64 fail!\n%#v\n", err)
	}
	fmt.Printf("SignBase64 done：%v\n", sigB64)

	err = VerifyBase64(publicKey, msg, sigB64)
	if err != nil {
		t.Errorf("VerifyBase64 fail!\n%#v\n", err)
	}
	fmt.Printf("SignBase64 done\n")
}
