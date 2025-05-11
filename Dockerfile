FROM alpine:latest

# 安装运行时可能需要的基本依赖
RUN apk add --no-cache ca-certificates alsa-lib

# 设置工作目录
WORKDIR /app

# 复制预编译的二进制文件到容器中
COPY ./bin/term-rex /app/term-rex

# 复制资源文件
COPY ./assets /app/assets/

# 设置环境变量
ENV TERM=xterm-256color

# 设置可执行权限
RUN chmod +x /app/term-rex

# 运行应用
CMD ["/app/term-rex"]
