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

	"github.com/Bastien2203/pi-proxy/middlewares"
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
		HostPolicy: autocert.HostWhitelist(getAllHosts(config)...),
	}

	server := &http.Server{
		Addr: ":443",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if target, ok := (*config)[r.Host]; ok {
				middlewares.ApplyMiddlewares(reverseProxy(target), target.Middlewares).ServeHTTP(w, r)
			} else {
				http.Error(w, "Not Found", http.StatusNotFound)
			}
		}),
		TLSConfig: &tls.Config{GetCertificate: manager.GetCertificate},
	}

	fmt.Println("Server is running on port 443 with HTTPS")
	if err := server.ListenAndServeTLS("", ""); err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
}

func getAllHosts(config *ProxyConfig) []string {
	var hosts []string
	for host := range *config {
		hosts = append(hosts, host)
	}
	return hosts
}
