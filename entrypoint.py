import socket
import httpx
import re
from multiprocessing import Process
from time import sleep
import subprocess


def get_open_ports():
    proc = subprocess.run(["ss", "-tunl"], capture_output=True, text=True)
    return list(map(int, set(re.findall("(?<=[^:]:)\d+", proc.stdout))))


def post_to_server(path, payload):
    try:
        transport = httpx.HTTPTransport(uds="/tmp/server.sock")
        client = httpx.Client(transport=transport)
        return client.post(f"http://api{path}", json=payload)
    except Exception as e:
        print(f"Error sending open ports to the server: {e}")


def report_open_ports(open_ports):
    payload = {"container_id": socket.getfqdn(), "open_ports": open_ports}
    response = post_to_server("/portreport", payload)
    if response.status_code == 200:
        print("Successfully sent open ports to the server.")
    else:
        print(f"Failed to send open ports. Server returned status code: {response.status_code}")


def report_session_end():
    payload = {"container_id": socket.getfqdn()}
    post_to_server("/sessionend", payload)


def port_scan(interval=3):
    open_ports = None
    while True:
        print("Checking ports")
        found_ports = get_open_ports()
        if open_ports != found_ports:
            open_ports = found_ports
            report_open_ports(open_ports)
        sleep(interval)


if __name__ == "__main__":
    print("Hello world")
    Process(target=port_scan, daemon=True).start()
    subprocess.run(["ttyd", "--writable", "--once", "bash"])
    report_session_end()
