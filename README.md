# golang-concurrency
Задание:
Написать утилиту, которая будет вычислять "производительность". Утилита параллельно рассылает http-запросы на какой-либо ресурс или несколько ресурсов (например, https://google.com).
Программа принимает следующие параметры на вход:
- адрес ресурса(ов)
- количество запросов, которые необходимо выполнить
- таймаут (время ожидания, после которого мы ответа уже не ждём)

Программа собирает и выводит на экран следующие данные:
- время, за которое отработали все запросы
- среднее время на запрос
- максимальное/минимальное время возвращение ответа
- количество ответов, которых не дождались

Пример работы программы:
![alt text](https://github.com/bukhavtsov/golang-concurrency/blob/task-2/img/screenshots/screenshot_1.png)
![alt text](https://github.com/bukhavtsov/golang-concurrency/blob/task-2/img/screenshots/screenshot_2.png)
