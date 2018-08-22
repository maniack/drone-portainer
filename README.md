# drone-portainer

Deploy docker stack to Portainer

```
pipeline:
  deploy:
    image: maniack/drone-portainer
    portainer:
      address: http://portainer:5000
      insecure: true
      endpoint: local
    stack:
      name: nginx
      path: docker-stack.yml
      environment:
        - DEBUG=true
    secrets: [ portainer_username, portainer_password ]
```
