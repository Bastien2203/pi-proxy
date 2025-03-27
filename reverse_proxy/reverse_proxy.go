package reverse_proxy

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/Bastien2203/pi-proxy/middlewares"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

type Conf struct {
	Host        string                   `json:"host"`
	Port        uint16                   `json:"port"`
	Middlewares []middlewares.Middleware `json:"middlewares"`
}

type ProxyConfig map[string]Conf

func reverseProxy(target Conf) http.Handler {
	targetURL := fmt.Sprintf("http://%s:%d", target.Host, target.Port)
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(parsedURL)

	return proxy
}

func ReadProxyConfig() *ProxyConfig {
	configFile, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	defer configFile.Close()

	byteValue, _ := io.ReadAll(configFile)
	var config ProxyConfig
	json.Unmarshal(byteValue, &config)
	return &config
}

func RunReverseProxyServer(config *ProxyConfig) {
	manager := &autocert.Manager{
		Cache:      autocert.DirCache("certs"),
		Prompt:     autocert.AcceptTOS,
		Email:      "bastiengrisvard2203@gmail.com",
		HostPolicy: autocert.HostWhitelist(getDomains(config)...),
		Client:     &acme.Client{DirectoryURL: acme.LetsEncryptURL},
	}

	go http.ListenAndServe(":80", manager.HTTPHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only redirect if it's not an ACME challenge
		if !strings.HasPrefix(r.URL.Path, "/.well-known/acme-challenge/") {
			http.Redirect(w, r, "https://"+r.Host+r.URL.RequestURI(), http.StatusMovedPermanently)
		} else {
			manager.HTTPHandler(nil).ServeHTTP(w, r)
		}
	})))

	server := &http.Server{
		Addr: ":443",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if target, ok := (*config)[r.Host]; ok {
				middlewares.ApplyMiddlewares(reverseProxy(target), target.Middlewares).ServeHTTP(w, r)
			} else {
				http.Error(w, "Not Found", http.StatusNotFound)
			}
		}),
		TLSConfig: &tls.Config{
			GetCertificate: manager.GetCertificate,
			NextProtos:     []string{"h2", "http/1.1", acme.ALPNProto},
		},
	}

	fmt.Println("Server is running on port 443 with HTTPS")
	if err := server.ListenAndServeTLS("", ""); err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
}

func getDomains(config *ProxyConfig) []string {
	domains := make([]string, 0, len(*config))
	for domain := range *config {
		domains = append(domains, domain)
	}
	return domains
}
