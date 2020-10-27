# aviasales

```
Нужно разработать HTTP сервис для быстрого поиска анаграмм в словаре.
Два слова считаются анаграммами, если одно можно получить из другого перестановкой букв (без учета регистра).
Примеры анаграмм:
["foobar", "barfoo", "boofar"]
["живу", "вижу"]
["Abba", "BaBa"]
Примеры строк, не являющихся анаграммами:
["abba", "bba"] - во второй строке только одна буква "а"
Сервис должен предоставлять эндпойнт для загрузки списка слов в формате json. Пример использования:
curl localhost:8080/load -d '["foobar", "aabb", "baba", "boofar", "test"]'
И эндпойнт для поиска анаграмм по слову в загруженном словаре. Примеры использования:
curl 'localhost:8080/get?word=foobar' => ["foobar","boofar"]
curl 'localhost:8080/get?word=raboof' => ["foobar","boofar"]
curl 'localhost:8080/get?word=abba' => ["aabb","baba"]
curl 'localhost:8080/get?word=test' => ["test"]
curl 'localhost:8080/get?word=qwerty' => null
