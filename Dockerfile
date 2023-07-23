FROM tsl0922/ttyd:alpine

RUN apk update \
 && apk add alpine-sdk python3 py3-pip nodejs npm vim iproute2 \
 && python -m pip install requests httpx

COPY entrypoint.py /tmp/entrypoint.py

CMD ["python", "/tmp/entrypoint.py"]