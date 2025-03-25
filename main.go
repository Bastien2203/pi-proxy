package main

import (
	"github.com/Bastien2203/pi-proxy/reverse_proxy"
)

func main() {
	proxy_config := reverse_proxy.ReadProxyConfig()
	reverse_proxy.RunReverseProxyServer(proxy_config)
	select {}
}
