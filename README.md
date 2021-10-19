1. 创建镜像
docker build -t httpdemo -f .\Dockerfile .

2. 运行容器
docker run -d -p 8888:8888 --name httpdemo httpdemo

3. 获取容器PID
docker insepct -f {{ State.Pid }} $ContainerId

4. 进入容器
nsenter --target $PID --mount --uts --ipc --net --pid

5. 推送至官方镜像库
docker login
docker tag httpdemo $user/httpdemo:v1
docker push $user/httpdemo
