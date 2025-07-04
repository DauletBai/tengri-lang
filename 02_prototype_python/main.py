# 02_prototype_python/main.py

from token import Token
from lexer import Lexer
from parser import Parser
from interpreter import Interpreter # <-- Импортируем нашего Исполнителя!

# Код на ОЙ, который использует переменные
oy_code = """
// Создаем переменную и константу
— □ a : 10
Λ □ b : 5

// Вычисляем выражение, используя созданные переменные
— □ c : a * (b + 2) 
"""

# --- ЭТАП 1: ЛЕКСЕР ---
lexer = Lexer(oy_code)
tokens = []
while True:
    token = lexer.get_next_token()
    tokens.append(token)
    if token.type == 'EOF': break

print("--- 1. Токены от Лексера ---")
print(tokens)

# --- ЭТАП 2: ПАРСЕР ---
parser = Parser(tokens)
ast = parser.parse()
print("\n--- 2. Древо Мысли от Парсера ---")
print(ast)

# --- ЭТАП 3: ИСПОЛНИТЕЛЬ ---
interpreter = Interpreter()
# Исполняем программу, обходя Древо Мысли
interpreter.interpret(ast)

print("\n--- 3. Состояние Памяти после Исполнения ---")
# Заглянем в память нашего Исполнителя, чтобы увидеть результат
print(interpreter.environment.variables)