# 02_prototype_python/parser.py

from token import Token
from ast_nodes import *

class Parser:
    """«Құрастырушы» — архитектор. Финальная версия прототипа."""
    def __init__(self, tokens):
        self.tokens = tokens
        self.position = 0

    def get_current_token(self):
        return self.tokens[self.position] if self.position < len(self.tokens) else self.tokens[-1]

    def advance(self):
        if self.position < len(self.tokens):
            self.position += 1

    # --- РАЗБОР ВЫРАЖЕНИЙ (Expressions) ---

    def parse_factor(self):
        token = self.get_current_token()
        if token.type == 'IntegerLiteral':
            self.advance()
            return NumberNode(token)
        elif token.type == 'Identifier':
            if self.position + 1 < len(self.tokens) and self.tokens[self.position + 1].type == 'Sep_LParen':
                return self.parse_function_call()
            else:
                self.advance()
                return VarAccessNode(token)
        elif token.type == 'Sep_LParen':
            self.advance()
            node = self.parse_expression()
            if self.get_current_token().type != 'Sep_RParen':
                raise Exception("Ошибка: Ожидалась ')'")
            self.advance()
            return node
        raise Exception(f"Ошибка: Неожиданный токен в 'factor': {token}")

    def parse_term(self):
        node = self.parse_factor()
        while self.get_current_token().type in ('Op_Multiply', 'Op_Divide'):
            op = self.get_current_token(); self.advance()
            node = BinOpNode(left_node=node, op_token=op, right_node=self.parse_factor())
        return node

    def parse_expression(self):
        node = self.parse_term()
        while self.get_current_token().type in ('Op_Plus', 'Op_Minus'):
            op = self.get_current_token(); self.advance()
            node = BinOpNode(left_node=node, op_token=op, right_node=self.parse_term())
        return node

    def parse_function_call(self):
        name_token = self.get_current_token()
        self.advance()
        self.advance() # (
        args = []
        while self.get_current_token().type != 'Sep_RParen':
            args.append(self.parse_expression())
            if self.get_current_token().type == 'Sep_Comma':
                self.advance()
        self.advance() # )
        return FuncCallNode(name_token, args)

    # --- РАЗБОР КОМАНД (Statements) ---

    def parse_declaration(self):
        decl_rune = self.get_current_token(); self.advance()
        type_rune = self.get_current_token(); self.advance()
        identifier = self.get_current_token(); self.advance()
        if self.get_current_token().type != 'Op_Assign': raise Exception("Ошибка: Ожидался ':'")
        self.advance()
        value_node = self.parse_expression()
        if decl_rune.type == 'Runa_Const':
            return ConstDeclNode(type_rune, identifier, value_node)
        else:
            return VarDeclNode(type_rune, identifier, value_node)

    def parse_return_statement(self):
        self.advance() # →
        value_node = self.parse_expression()
        return ReturnNode(value_node)

    def parse_statement_list(self):
        statements = []
        while self.get_current_token().type != 'Sep_RParen':
            statements.append(self.parse_statement())
        return statements

    def parse_function_definition(self):
        self.advance() # Π
        name_token = self.get_current_token(); self.advance()
        self.advance() # (
        params = []
        while self.get_current_token().type != 'Sep_RParen':
            param_type = self.get_current_token(); self.advance()
            param_name = self.get_current_token(); self.advance()
            params.append(ParamNode(param_type, param_name))
            if self.get_current_token().type == 'Sep_Comma': self.advance()
        self.advance() # )
        return_type = None
        if self.get_current_token().type == 'Runa_Return':
            self.advance(); return_type = self.get_current_token(); self.advance()
        self.advance() # (
        body = self.parse_statement_list()
        self.advance() # )
        return FuncDefNode(name_token, params, return_type, body)

    def parse_statement(self):
        token_type = self.get_current_token().type
        if token_type in ('Runa_Const', 'Runa_Var'):
            return self.parse_declaration()
        elif token_type == 'Runa_Func_Def':
            return self.parse_function_definition()
        elif token_type == 'Runa_Return':
            return self.parse_return_statement()
        raise Exception(f"Ошибка: Неизвестная команда, начинающаяся с {self.get_current_token()}")

    def parse(self):
        """Главный метод, строит дерево для всей программы."""
        ast_tree = []
        while self.get_current_token().type != 'EOF':
            ast_tree.append(self.parse_statement())
        return ast_tree