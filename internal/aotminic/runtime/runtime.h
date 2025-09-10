#ifndef TENGRI_RUNTIME_H
#define TENGRI_RUNTIME_H

// Глобальные argv/argc для argi()
extern int __argc;
extern char **__argv;

// Примитивы рантайма
long print(long x);
long argi(long idx);

// Высокоточный монотонный таймер (ns)
long time_ns(void);

// Удобная печать метрики для бенча: "TIME_NS: <ns>"
long print_time_ns(long ns);

#endif // TENGRI_RUNTIME_H