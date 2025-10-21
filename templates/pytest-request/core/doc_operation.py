"""
    @Author: TheRuffian
    @Email: bugpz2779@gmail.com
    @Blog: 'https://gavin.us.kg'
    @Github: 'https://github.com/GavinHaydy'
"""

import pandas
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
