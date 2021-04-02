# Описание
Набор плагинов, шаблонов и других настроек для мониторинга Transporter (Microporter) с помощью zabbix.

# Плагин агента
### Сборка
zabbix_agent2 с плагинами собирается статически. То есть на выходе - один бинарник и нет возможности поставлять плагины в виде отдельных dll.
1. Скачать репозиторий zabbix нужной версии: `git clone https://github.com/zabbix/zabbix.git --depth 1 --branch release/5.2`
2. Выполнить конфигурацию как в документации с флагом `--enable-agent2`
3. Переместить содержимое (или сделать ссылку) на agent2_plugin в `zabbix/src/go/plugins/microporter`
4. В `zabbix/src/go/plugins/plugins_linux.go` добавить импорт `"zabbix.com/plugins/microporter"`
5. В папке `zabbix/src/go` выполнить `go get github.com/powerman/rpc-codec/jsonrpc2`
6. В папке `zabbix/src/go` выполнить `make`
7. `zabbix/src/go/bin/zabbix_agent2` - агент, `zabbix/src/go/conf/zabbix_agent2.conf` - конфиг

### Метрики
- `porter.stat[<port>]` - статистика сервиса
    - `bitrate`
    - `buffers-mem-usage`
	- `cpu-usage`
	- `mem-usage`
	- `pid`
	- `streaming-count`
	- `streams-count`
	- `threads-count`
	- `uptime`
- `porter.streams.discovery[<port>]` - низкоуровневое обнаружение потоков. Поля:
    - `{#ID}` - идентификатор потока
    - `{#NAME}` - имя потока
- `porter.streams[<port>]` - список всех потоков. Массив объектов со следующими полями:
    - `id` - идентификатор потока
    - `name` - имя потока
    - `enabled` - включен ли поток
    - `state` - статус потока
    - `uptime` - время работы потока
    - `bitrate` - битрейт потока

# Шаблон для Zabbix
Шаблон содержит набор метрик (items), триггеров и правило низкоуровневого обнаружения видеопотоков. В свою очередь, правило обнаружения создает метрики и треггеры.
Макросы (параметры) шаблона:
- `PORTER.PORT` - порт сервиса, по умолчанию 8066.


# Интеграция Telegram в Zabbix 
1. Создание бота:
Отправте "/newbot" боту @BotFather и следуйте инструкции. Токен, выданный @BotFather, будет использован для конфигурации webhook'a в Zabbix
2. Для того чтобы получать уведомления в группу, вы должны получить ID данной группы. Для этого создайте груповой чат и добавте в чат @myidbot и вашего созданного бота. Отправте "/getgroupid@myidbot". Затем отправте "/start@<имя_вашего_бота>".
3. Настройка Zabbix 
Перейдите в раздел "Administration > Media types". Выберете в списке Teleram. В настройка webhook'a необходимо ввести telegramToken, id группового чата и telegramParseMode (указывайте Markdown).
4. Для проверки работоспособность бота необходимо перейте на вкладку "Administration > Media types", найти в списке Telegram и нажать на кнопку тест. При правильной настройки вам должно прийти тестовое сообщение в указанный вами чат.
