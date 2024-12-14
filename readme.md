### 1. Параметры файла devices.json

```json
{
    "0190": {
        "Vendor": "Beward DKS (DTMF-открытие при вызове)",
        "Adress": "Сысоева 18 под 3",
        "URL": "http://log:pass@192.168.1.1:8080",
    }
}
```

* Vendor - Указание вендора производителя
* Адрес - Адрес для отображения в TG
* URL - Адрес оборудования, если на панели требуется авторизация укажите ее в url

### 2. Структура Alert и Devices

```go
type Alert struct {
	Addr		string
	Time 		time.Time
	Status		int
	LogAlert 	bool
}

type Devices struct {
	Vendor 	string `json:"Vendor"`
	Adress 	string `json:"Adress"`
	Url 	string `json:"URL"`
}

```

```go
/*

	ПРИМЕРЫ АЛЕРТОВ telegram

*/

// При статусе StatusTimeout
/*
❗ ERROR (StatusTimeout)
	Addr: Груднова д.1 к.2 Под 1
	Time: 12.12.24 (21:37:29)
	Status: Not available
	Reason: 
	 - Заданный узел не ответил в течении заданного 
	 - {{timeout_second_request}}
*/

// При статусе 401
/*
⚠️ Warning
	Addr: Груднова д.1 к.2 Под 1
	Time: 12.12.24 (21:37:29)
	Status: 401
	Reason: 
	 - 401 Unauthorized - неверная авторизация для проверки статуса
*/

// При статусе 404
/*
⚠️ Warning
	Addr: Груднова д.1 к.2 Под 1
	Time: 12.12.24 (21:37:29)
	Status: 404
	Reason:
	 - 404 Site or Page Not Found - указанная страница не существует
	 - Проверьте правельность заполнения URL
*/

// При статусе 400
/*
⚠️ Warning
	Addr: Груднова д.1 к.2 Под 1
	Time: 12.12.24 (21:37:29)
	Status: 400
	Reason:
	 - 400 Bad Request - Неверный переданный запрос
	 - Проверьте правельность заполнения CGI в запросе
*/

// При статусе восстановлен
/*
✅ Restored
	Addr:   Груднова д.1 к.2 Под 1
	Time:   12.12.24 (21:37:29)
	Status: 200
		SIP:      AccountReg1=1
		Error:    10
		Label:    200 key
		ScanCode: 12345=>on
		DoorCode: 15926=>on
		Lock:     included
*/
```
