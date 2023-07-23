import socket
import httpx
import re
from multiprocessing import Process
from time import sleep
import subprocess


def get_open_ports():
    proc = subprocess.run(["ss", "-tunl"], capture_output=True, text=True)
    return list(map(int, set(re.findall("(?<=[^:]:)\d+", proc.stdout))))


def report_open_ports(open_ports):
    try:
        transport = httpx.HTTPTransport(uds="/tmp/server.sock")
        client = httpx.Client(transport=transport)
        payload = {"container_id": socket.getfqdn(), "open_ports": open_ports}
        response = client.post(f"http://api/portreport", json=payload)
        if response.status_code == 200:
            print("Successfully sent open ports to the server.")
        else:
            print(f"Failed to send open ports. Server returned status code: {response.status_code}")
    except Exception as e:
        print(f"Error sending open ports to the server: {e}")


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
    Process(target=port_scan).start()
    subprocess.run(["ttyd", "-W", "bash"])
