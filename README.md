# Tenge Programming Language

**Tenge** — новый язык программирования, выросший из прототипа *tengri-lang*.  
Он основан на идеях агглютинативности казахского языка и ориентирован на **финтех и вычислительные задачи**.

## Mission
- Создать язык, в котором архетипы и агглютинативность казахского языка отражаются в структуре кода.
- Дать инструмент для высокопроизводительных приложений в **финансовом секторе**.
- Соревноваться с C, Rust и Go на уровне низкоуровневых вычислений.

## Benchmarks (средние результаты трёх прогонов)

| Task                | C        | Go       | Rust     | Tenge (best)     |
|---------------------|----------|----------|----------|------------------|
| Fib iter (90)       | ~41 ns   | ~5170 ns | ~197 ns  | ~50 ns           |
| Fib rec (35)        | ~43 ms   | ~51 ms   | ~46 ms   | ~44 ms           |
| Sort (100k)         | ~0.49 ms | ~0.11 ms | ~1.67 ms | ~0.82 ms (PDQ)   |
| VaR Monte Carlo(1e6)| ~185 ms  | ~177 ms  | ~85 ms   | ~31 ms (Zig+QSel)|

Tenge уже показал результаты уровня C на `fib_iter` и обогнал Rust/Go/C в Monte Carlo VaR благодаря Quickselect.

## Roadmap
- [x] Реализовать AOT-компиляцию
- [x] Добавить benchmarks (fib, sort, Monte Carlo VaR)
- [x] Улучшить sort (PDQ, radix)
- [x] Добавить Ziggurat RNG и Quickselect для VaR
- [ ] Добавить C/Rust/Go quickselect-версии (для честного сравнения)
- [ ] Доработать radix sort (16-bit passes)
- [ ] Расширить синтаксис языка (финтех-конструкции)
- [ ] Выпустить whitepaper

---

Исторический прототип: [tengri-lang](https://github.com/DauletBai/tengri-lang)

Link to the new [tenge](https://github.com/DauletBai/tenge) project, based on [tengri-lang](https://github.com/DauletBai/tengri-lang)