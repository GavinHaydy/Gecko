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
