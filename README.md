# drone-portainer

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
    secrets: [ portainer_usersname, portainer_password ]
```
