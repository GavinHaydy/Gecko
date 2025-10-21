"""
    @Author: Gavin
    @Email: bugpz2779@gmail.com
    @Blog: 'https://gavin.us.kg'
    @Github: 'https://github.com/GavinHaydy'
"""

import logging
import os
import sys
from datetime import datetime


class ColorFormatter(logging.Formatter):
    COLORS = {
        "DEBUG": "\033[37m",
        "INFO": "\033[36m",
        "WARNING": "\033[33m",
        "ERROR": "\033[31m",
        "CRITICAL": "\033[41m",
    }
    RESET = "\033[0m"

    def format(self, record):
        color = self.COLORS.get(record.levelname, "")
        message = super().format(record)
        return f"{color}{message}{self.RESET}"


LOG_DIR = os.path.join(os.path.dirname(__file__), "..", "logs")
os.makedirs(LOG_DIR, exist_ok=True)
log_file = os.path.join(LOG_DIR, f"{datetime.now().strftime('%Y-%m-%d')}.log")

fmt = "%(asctime)s | %(levelname)-8s | %(filename)s:%(lineno)d | %(message)s"
formatter = logging.Formatter(fmt, "%Y-%m-%d %H:%M:%S")

file_handler = logging.FileHandler(log_file, encoding="utf-8")
file_handler.setFormatter(formatter)

console_handler = logging.StreamHandler(sys.stdout)
console_handler.setFormatter(ColorFormatter(fmt))

logger = logging.getLogger("gecko")
logger.setLevel(logging.DEBUG)
logger.addHandler(console_handler)
logger.addHandler(file_handler)
