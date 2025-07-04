# 02_prototype_python/ast_nodes.py

class ASTNode:
    """Базовый класс для всех узлов дерева."""
    pass

# ... (NumberNode, BinOpNode, VarDeclNode, ConstDeclNode остаются без изменений)
class NumberNode(ASTNode):
    def __init__(self, token):
        self.token = token
        self.value = token.value
    def __repr__(self):
        return f"NumberNode({self.value})"

class BinOpNode(ASTNode):
    def __init__(self, left_node, op_token, right_node):
        self.left_node = left_node
        self.op_token = op_token
        self.right_node = right_node
    def __repr__(self):
        return f"BinOpNode(left={self.left_node}, op='{self.op_token.value}', right={self.right_node})"

class VarDeclNode(ASTNode):
    def __init__(self, type_rune_token, identifier_token, value_node):
        self.type_rune = type_rune_token
        self.identifier = identifier_token
        self.value_node = value_node
    def __repr__(self):
        return f"VarDecl(name='{self.identifier.value}', type='{self.type_rune.value}', value={self.value_node})"

class ConstDeclNode(ASTNode):
    def __init__(self, type_rune_token, identifier_token, value_node):
        self.type_rune = type_rune_token
        self.identifier = identifier_token
        self.value_node = value_node
    def __repr__(self):
        return f"ConstDecl(name='{self.identifier.value}', type='{self.type_rune.value}', value={self.value_node})"

# --- НОВЫЕ УЗЛЫ ДЛЯ ФУНКЦИЙ ---

class ParamNode(ASTNode):
    """Узел для описания одного параметра функции (например, □ san)."""
    def __init__(self, type_rune_token, identifier_token):
        self.type_rune = type_rune_token
        self.identifier = identifier_token
    def __repr__(self):
        return f"Param(name='{self.identifier.value}', type='{self.type_rune.value}')"

class FuncDefNode(ASTNode):
    """Узел для определения функции (Π)."""
    def __init__(self, name_token, params_list, return_type_rune, body_node):
        self.name_token = name_token
        self.params = params_list
        self.return_type = return_type_rune
        self.body = body_node
    def __repr__(self):
        return f"FuncDef(name='{self.name_token.value}', params={self.params}, returns='{self.return_type.value if self.return_type else 'None'}', body={self.body})"

class FuncCallNode(ASTNode):
    """Узел для вызова функции."""
    def __init__(self, name_token, args_list):
        self.name_token = name_token
        self.args = args_list
    def __repr__(self):
        return f"FuncCall(name='{self.name_token.value}', args={self.args})"

class VarAccessNode(ASTNode):
    """Узел для доступа к значению переменной по имени."""
    def __init__(self, name_token):
        self.name_token = name_token
    def __repr__(self):
        return f"VarAccess(name='{self.name_token.value}')"

class ReturnNode(ASTNode):
    """Узел для команды возврата из функции (→)."""
    def __init__(self, value_node):
        self.value_to_return = value_node
    def __repr__(self):
        return f"Return({self.value_to_return})"