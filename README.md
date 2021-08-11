# go-template

Перед запуском проверить содержимое Makefile.env:

1. Не запускать команду git_clear в set_env, если клонирование было непосредственно через github template
2. Проверить имя переменной APP_NAME
3. Проверить автора репозитория в PACKAGE_NAME
4. Запустить команду make -f Makefile.env
