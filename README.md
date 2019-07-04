# deploy-console
- 镜像启动命令
    ```bash
     docker run --rm -it  -v /var/run/docker.sock:/var/run/docker.sock -p90:8080 -e ENV=local  192.168.31.188/test/console:v1.0.0
    ```