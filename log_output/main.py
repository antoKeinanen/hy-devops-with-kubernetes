
from datetime import datetime
import time
import uuid


if __name__ == "__main__":
    string = uuid.uuid4()
    while True:
        ts = datetime.now()
        print(f"{ts}: {string}")
        time.sleep(5)

