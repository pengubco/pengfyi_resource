"""Runtime bootstrap: runs before the workload's own code.

A real one would fetch config, wire up logging and metrics, then exec the
workload. This one prints enough to prove it came from the image cache.
"""
import platform
import sys

import requests


def main() -> int:
    print(f"python-runtime: python {platform.python_version()} on {platform.machine()}")
    print(f"python-runtime: requests {requests.__version__} preinstalled")
    if len(sys.argv) > 1:
        print(f"python-runtime: would exec workload: {' '.join(sys.argv[1:])}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
