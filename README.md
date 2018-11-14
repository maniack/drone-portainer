# drone-portainer

Deploy docker stack to Portainer

```
kind: pipeline
name: default

steps:
- name: deploy
  image: maniack/drone-portainer:1.0.0-rc.1
  settings:
    portainer: http://portainer:5000
    insecure: true
    username:
      from_secret: portainer_username
    password:
      from_secret: portainer_password
    endpoint: local
    stack: nginx
    file: docker-stack.yml
    environment:
      - DEBUG=true
    debug: true
```