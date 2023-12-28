import threading
import main

if __name__ == "__main__":
    main.InitServer()
    t2 = threading.Thread(target=main.ShellOsHttp)
    t1 = threading.Thread(target=main.HeartbeatHttp)
    t2.start()
    t1.start()
