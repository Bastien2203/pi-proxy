package reverse_proxy

import (
	"crypto/tls"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

func GetCertificate() *tls.Config {
	manager := &autocert.Manager{
		Cache:  autocert.DirCache("certs"),
		Prompt: autocert.AcceptTOS,
		Email:  "bastiengrisvard2203@gmail.com",
		Client: &acme.Client{DirectoryURL: acme.LetsEncryptURL},
	}

	return &tls.Config{
		GetCertificate: manager.GetCertificate,
		NextProtos:     []string{"h2", "http/1.1", acme.ALPNProto},
	}
}
