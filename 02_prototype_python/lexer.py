# 02_prototype_python/lexer.py # В реальном проекте это будет в отдельном файле
from token import Token

class Lexer:
    """«Танушы» — распознающий. Финальная версия прототипа."""
    def __init__(self, source_code):
        self.code = source_code
        self.position = 0
        # Полный словарь для мгновенного распознавания одиночных символов
        self.token_map = {
            'Π': 'Runa_Func_Def', '↑': 'Runa_EntryPoint', 'Y': 'Runa_If',
            'Q': 'Runa_True', 'I': 'Runa_False', '↻': 'Runa_Loop',
            '→': 'Runa_Return', '⁞': 'Runa_Log', '—': 'Runa_Var',
            'Λ': 'Runa_Const', '□': 'Runa_Type_Int', '⊡': 'Runa_Type_Float',
            '∞': 'Runa_Type_Str', '◇': 'Runa_Type_Char', '≡': 'Runa_Type_Collection',
            '+': 'Op_Plus', '-': 'Op_Minus', '*': 'Op_Multiply', '/': 'Op_Divide',
            ':': 'Op_Assign', '@': 'Op_Access_At', '#': 'Op_Get_Size',
            '∈': 'Op_In', '(': 'Sep_LParen', ')': 'Sep_RParen',
            '[': 'Sep_LBracket', ']': 'Sep_RBracket', ',': 'Sep_Comma'
        }

    def _skip_noise(self):
        """Пропускает пробелы и комментарии."""
        while self.position < len(self.code):
            if self.code[self.position].isspace():
                self.position += 1
                continue
            if self.code[self.position:self.position+2] == '//':
                while self.position < len(self.code) and self.code[self.position] != '\n':
                    self.position += 1
                continue
            break

    def _read_identifier(self):
        """Читает идентификатор (имя)."""
        start_pos = self.position
        while self.position < len(self.code) and self.code[self.position].isalnum():
            self.position += 1
        return self.code[start_pos:self.position]

    def _read_number(self):
        """Читает целое или вещественное число."""
        start_pos = self.position
        is_float = False
        while self.position < len(self.code) and self.code[self.position].isdigit():
            self.position += 1
        if self.position < len(self.code) and self.code[self.position] == '.':
            is_float = True
            self.position += 1
            while self.position < len(self.code) and self.code[self.position].isdigit():
                self.position += 1
        
        num_str = self.code[start_pos:self.position]
        return Token('FloatLiteral', float(num_str)) if is_float else Token('IntegerLiteral', int(num_str))

    def get_next_token(self):
        """Главный метод, возвращающий следующий токен из кода."""
        self._skip_noise()

        if self.position >= len(self.code):
            return Token('EOF', None)

        current_char = self.code[self.position]

        # Проверка на составные операторы
        if current_char == '=' and self.position + 1 < len(self.code) and self.code[self.position+1] == '=':
            self.position += 2
            return Token('Op_Equal', '==')
        if current_char == '<' and self.position + 1 < len(self.code) and self.code[self.position+1] == '-':
            self.position += 2
            return Token('Op_Push', '←')

        # Проверка по словарю одиночных символов
        if current_char in self.token_map:
            token_type = self.token_map[current_char]
            token = Token(token_type, current_char)
            self.position += 1
            return token

        # Проверка на число
        if current_char.isdigit():
            return self._read_number()

        # Проверка на идентификатор
        if current_char.isalpha():
            identifier = self._read_identifier()
            return Token('Identifier', identifier)

        raise Exception(f"Ошибка: Нераспознанный символ '{current_char}'")