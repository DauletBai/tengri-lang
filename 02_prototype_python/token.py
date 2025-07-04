# 02_prototype_python/token.py
class Token:
    """Описывает один 'распознанный' элемент кода."""
    def __init__(self, type, value, line=1, column=1):
        self.type = type
        self.value = value
        self.line = line
        self.column = column

    def __repr__(self):
        """Метод для красивого вывода токена при печати."""
        return f"Token({self.type}, '{self.value}')"