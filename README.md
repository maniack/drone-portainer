# drone-portainer

Deploy docker stack to Portainer

```
pipeline:
  deploy:
    image: maniack/drone-portainer
    portainer: http://portainer:5000
    insecure: true
    endpoint: local
    stack: nginx
    file: docker-stack.yml
    environment:
      - DEBUG=true
    debug: true
    secrets: [ portainer_username, portainer_password ]
```

```
pipeline:
  deploy:
    image: maniack/drone-portainer
    portainer: http://portainer:5000
    insecure: true
    endpoint: local
    stack: nginx
    environment:
      - DEBUG=true
    config: |
            version: '3.5'
            services:
              test:
                image: alpine:latest
    debug: true
    secrets: [ portainer_username, portainer_password ]
```