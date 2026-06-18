import random
from datetime import datetime, timezone


def new_request_id(prefix: str = "AI") -> str:
    ts = datetime.now(timezone.utc).strftime("%Y%m%d%H%M%S%f")[:18]
    rand = str(random.randint(0, 999999)).zfill(6)
    return f"{prefix}-{ts}-{rand}"
