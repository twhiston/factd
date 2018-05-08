FROM frolvlad/alpine-glibc:latest

LABEL maintainer="Tom Whiston <tom.whiston@gmail.com>"
LABEL RUN="docker run -it --privileged                       \
                --ipc=host --net=host --pid=host -e HOST=/host                \
                -v /run:/run -v /var/log:/var/log                             \
                -v \${HOME}/:/root \
                -v /etc/localtime:/etc/localtime -v /:/host tomwhiston/factd:latest run"
ENV HOST_VAR=/host/var \
    HOST_SYS=/host/sys \
    HOST_PROC=/host/proc \
    HOST_ETC=/host/etc

COPY ./dist/factd /factd

ENTRYPOINT ["/factd"]