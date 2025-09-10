# 02_prototype_python/interpreter.py

# Импортируем все наши узлы, чтобы Исполнитель их "знал"
from ast_nodes import *

class Environment:
    """Хранит переменные и их значения в определенной области видимости."""
    def __init__(self):
        self.variables = {}

    def set(self, name, value):
        self.variables[name] = value

    def get(self, name):
        if name not in self.variables:
            raise Exception(f"Ошибка исполнения: Переменная '{name}' не определена.")
        return self.variables[name]

class Interpreter:
    """«Іске Асырушы» — обходит Древо Мысли и выполняет его."""
    def __init__(self):
        # У каждой сессии исполнения своя "глобальная" память
        self.environment = Environment()

    def visit(self, node):
        """Главный метод-диспетчер. Вызывает нужный метод-посетитель для каждого узла."""
        method_name = f'visit_{type(node).__name__}'
        visitor = getattr(self, method_name, self.no_visit_method)
        return visitor(node)

    def no_visit_method(self, node):
        raise Exception(f"Ошибка: Нет метода visit_{type(node).__name__} для обработки узла {type(node).__name__}.")

    # --- Посетители для Выражений ---

    def visit_NumberNode(self, node):
        """Возвращает числовое значение узла."""
        return node.value

    def visit_BinOpNode(self, node):
        """Вычисляет результат бинарной операции."""
        left_val = self.visit(node.left_node)
        right_val = self.visit(node.right_node)

        op = node.op_token.type
        if op == 'Op_Plus': return left_val + right_val
        if op == 'Op_Minus': return left_val - right_val
        if op == 'Op_Multiply': return left_val * right_val
        if op == 'Op_Divide': return left_val / right_val

    def visit_VarAccessNode(self, node):
        """Извлекает значение переменной из памяти."""
        var_name = node.name_token.value
        return self.environment.get(var_name)

    # --- Посетители для Команд ---

    def visit_ConstDeclNode(self, node):
        """Сохраняет константу в память."""
        var_name = node.identifier.value
        value = self.visit(node.value_node)
        self.environment.set(var_name, value)
        return None

    def visit_VarDeclNode(self, node):
        """Сохраняет переменную в память."""
        var_name = node.identifier.value
        value = self.visit(node.value_node)
        self.environment.set(var_name, value)
        return None

    def interpret(self, ast_tree):
        """Главный метод, запускающий исполнение всей программы."""
        result = None
        for node in ast_tree:
            result = self.visit(node)
        return result