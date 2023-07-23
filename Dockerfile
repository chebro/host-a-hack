FROM tsl0922/ttyd:alpine

RUN apk update && apk add python3 py3-pip nodejs npm

CMD ["ttyd", "-W", "bash"]