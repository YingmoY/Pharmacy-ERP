"""Example local service manager configuration.

Copy this file to service_manager.py and adjust paths for your machine.
The real service_manager.py is ignored because it often contains local
absolute paths and private RabbitMQ credentials.
"""

import os
import subprocess

BASE = os.path.dirname(os.path.abspath(__file__))
PNPM = r"C:\Path\To\pnpm.cmd"
AIR = r"C:\Path\To\air.exe"

SERVICES = [
    {
        "name": "PharmacyERP",
        "cwd": os.path.join(BASE, "PharmacyERP"),
        "cmd": [AIR],
    },
    {
        "name": "PharmacyERP_front",
        "cwd": os.path.join(BASE, "PharmacyERP_front"),
        "cmd": [PNPM, "run", "dev"],
    },
    {
        "name": "medicare-gateway",
        "cwd": os.path.join(BASE, "medicare-gateway"),
        "cmd": ["go", "run", ".\\cmd\\gateway"],
        "env": {
            "MEDICARE_RABBIT_URL": "amqp://guest:guest@localhost:5672/",
            "MEDICARE_ENABLE_RABBITMQ": "true",
        },
    },
    {
        "name": "ai-services",
        "cwd": os.path.join(BASE, "ai-services"),
        "cmd": [os.path.join(BASE, "ai-services", ".venv", "Scripts", "python.exe"), "main.py"],
    },
]


def main() -> None:
    processes = []
    for service in SERVICES:
        env = os.environ.copy()
        env.update(service.get("env", {}))
        processes.append(
            subprocess.Popen(
                service["cmd"],
                cwd=service["cwd"],
                env=env,
            )
        )

    try:
        for process in processes:
            process.wait()
    except KeyboardInterrupt:
        for process in processes:
            process.terminate()


if __name__ == "__main__":
    main()
