package templates

import (
	"Gecko/internal/pkg/biz/tools"
	"Gecko/internal/pkg/dal/rao"
)

func initFunc() rao.Template {
	content := `
`
	return rao.Template{
		Content:  content,
		FileName: "core/__init__.py",
	}
}

func confFunc() rao.Template {
	content := `"""
    @Author: Gavin
    @Email: bugpz2779@gmail.com
    @Blog: 'https://gavin.us.kg'
    @Github: 'https://github.com/GavinHaydy'
"""
import redis
import yaml
import configparser
import os


def rds(host: str = 'localhost', port: int = 6379, db: int = 0, password: str = None):
    """

    Args:
        host: redis server ip
        port: redis server port
        db: redis db number
        password: redis password -> str|None

    Returns: redis client

    """
    try:
        if password:
            return redis.StrictRedis(host, port, db, password)
        else:
            return redis.StrictRedis(host=host, port=port, db=db)
    except Exception as e:
        return e


def read_yaml(filepath):
    """

    Args:
        filepath: yaml file path

    Returns: dict

    """
    return yaml.load(open(filepath, 'r').read(), yaml.FullLoader)


def read_ini(filename, encoding='utf-8'):
    """
    读取 config 目录下的 ini 文件，并返回 configparser.ConfigParser 对象。

    :param filename: ini 文件名（支持相对路径）
    :param encoding: 文件编码
    :return: configparser.ConfigParser 对象
    :raises FileNotFoundError: 文件不存在时抛出
    """
    base_dir = os.path.dirname(os.path.abspath(__file__))
    config_path = os.path.join(base_dir, '../config', filename)

    if not os.path.exists(config_path):
        raise FileNotFoundError(f"配置文件不存在: {config_path}")

    config = configparser.ConfigParser()
    config.read(config_path, encoding=encoding)

    return config
`
	return rao.Template{
		Content:  content,
		FileName: "core/conf.py",
	}
}

func docOperationFunc() rao.Template {
	content := `import pandas
import json
import numpy


class NpEncoder(json.JSONEncoder):
    def default(self, obj):
        if isinstance(obj, numpy.integer):  # numpy.integer 是int64
            return int(obj)
        elif isinstance(obj, numpy.floating):
            return float(obj)
        elif isinstance(obj, numpy.ndarray):
            return obj.tolist()
        else:
            return super(NpEncoder, self).default(obj)


class OperationsOfData:
    def __init__(self, filepath, **kwargs):
        """

        Args:
            file_name: data file path
            **kwargs:
        """
        # file_type = os.path.splitext(self.file)[-1][1:]
        self.file = filepath
        if self.file.endswith('csv'):
            self.data = pandas.read_csv(self.file, **kwargs)
        elif self.file.endswith('xls'):
            self.data = pandas.read_excel(self.file, **kwargs)
        elif self.file.endswith('xlsx'):
            self.data = pandas.read_excel(self.file, **kwargs)
        else:
            raise TypeError("An unsupported file type")

    def to_list(self):
        if self.data is not None:
            result = json.loads(json.dumps(self.data.values.tolist(), cls=NpEncoder))
            if len(result) == 0:
                raise ValueError("Empty File")
            else:
                return result

    def to_dict(self):
        if self.data is not None:
            result = json.dumps([self.data.loc[i].to_dict() for i in self.data.index.values],
                                cls=NpEncoder, ensure_ascii=False)
            if len(result) <= 2:
                raise ValueError("Empty File")
            else:
                return result
`
	return rao.Template{
		Content:  content,
		FileName: "core/doc_operation.py",
	}
}

func getKeyWord() rao.Template {
	content := `import jsonpath


class GetKeyword:
    @staticmethod
    def get_keyword(source_data, keyword):
        """
        通过关键字获取对应的值,如果有多个值,默认获取第一个,如果没有返回False
        :param source_data: 源数据
        :param keyword: 关键字
        :return: 关键字对应的第一个值/错误信息
        """

        try:
            return jsonpath.jsonpath(source_data, f'$..{keyword}')[0]
        except TypeError:
            return False

    @staticmethod
    def get_keywords(source_data, keyword):
        """
        通过关键字获取对应的所有值,如果不存在,返回False
        :param source_data: 源数据
        :param keyword: 关键字
        :return: list/错误信息
        """
        return jsonpath.jsonpath(source_data, f'$..{keyword}')

`
	return rao.Template{
		Content:  content,
		FileName: "core/get_keyword.py",
	}
}

func logger() rao.Template {
	content := `import logging
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

`
	return rao.Template{
		Content:  content,
		FileName: "core/logger.py",
	}
}

func ApiPytest() []rao.Template {
	return tools.MakeTemplates(
		initFunc, confFunc, docOperationFunc, getKeyWord, logger)
}

//func ApiPytest() []rao.Template {
//	var result []rao.Template
//	result = append(result, confFunc())
//	result = append(result, docOperationFunc())
//	result = append(result, getKeyWord())
//	return result
//}
