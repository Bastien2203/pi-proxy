pi-proxy
========


## Using docker compose

- Create `config.json` file in the root directory with your config : 

    - Key : domain name 
    - Value : 
        - host : targeted host of the service
        - port : targeted port of the service
        - middlewares : list of middlewares to apply to the service

```json
{
    "example.com" : {
        "host": "my_docker_service_or_ip_address",
        "port": 80,
        "middlewares": [
            {
                "name": "RateLimiter",
                "options": {
                    "maxRequests": 30,
                    "requestTimeout": 60
                }
            }
        ]
    }
}
```

- Create certs directory : `mkdir certs`
- Run `sudo docker compose up -d` to start the proxy


> Suggestions and feedback are welcome
