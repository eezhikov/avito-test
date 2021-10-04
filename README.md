Микросервис для работы с балансом пользователей

Список методов:

UserBalance - Получение баланса пользователя +  доп. задание (получение баланса в валюте)

    POST /balance - получение баланса в рублях

        параметры 
            id

        Request:
        {
            "id": 1
        }

        Response:
        {
            "balance": 600.5
        }

    POST /balance?currency=USD - получение баланса в валюте

        Request:
        {
            "id": 1
        }

        Response:
        {
            "balance": 8.234611697948134
        }


BalanceChange - Метод начисления и списания средств (объединен в один)

    POST /change

        параметры 
            id
            quality - сумма перевода

        Request:
        {
            "id": 1,
            "quality": 100.5
        }

        Request: 
        {
            "id": 1,
            "quality": -50
        }

        Response:
        {
            "credited": true
        }

MoneyTransaction - Метод перевода средств от пользователя к пользователю

    POST /transaction

        параметры 
            user_id_from - от кого
            user_id_to - кому
            quality - сумма

        Request:
        {
            "user_id_from": 1,
            "user_id_to": 2,
            "quality": 100.5
        }

        Response:
        {
            "transaction": true
        }

Комментарии: использовался сервис cbr-xml-daily.ru для валют. В сервисе из примера в ТЗ нельзя поменять валюту по умолчанию на рубли в бесплатном аккауте

По плану: добавить историю оперций по счету пользователя (2 доп. задание)
