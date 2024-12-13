/*
	Beward - CGI (GO)
	Regstatus
*/
package main

import (
	//"fmt"
	"log"
	"io/ioutil"
	//"errors"
	"net/http"
    "time"
	"os"
	"strings"
	"encoding/json"
	"regexp"
)

/*
	xType := fmt.Sprintf("%T", resp)
	fmt.Println(xType)

*/
var DEVICES_OBJ map[string]Devices

var matched_vendor	= regexp.MustCompile(`Beward`)
var matched_url		= regexp.MustCompile(`(http|s)(:\/\/)(\w*)(:)(\w*@)(([0-9]{1,3}[\.]){3}[0-9]{1,3})(:[0-9]{1,5})`)

type Alert struct {
	Addr		string
	Time 		time.Duration
	Status		int
	LogAlert 	bool
}

/*
	Devices.json
*/
type Devices struct {
	Vendor 	string `json:"Vendor"`
	Adress 	string `json:"Adress"`
	Url 	string `json:"URL"`
}

func main(){
	ListDevicesInit()
	
	for{
		for _, value := range DEVICES_OBJ {
			if	!matched_vendor.MatchString(value.Vendor) {
				continue
			}
			if	!matched_url.MatchString(value.Url) {
				continue
			}
			CheckingStatus(value.Url+"/cgi-bin/sip_cgi?action=regstatus", 4)
		}
		time.Sleep(60 * time.Second)
	}
}

/*
	Проверка статуса устройства
*/
func CheckingStatus(ip string, timeout time.Duration) {

	/*
		!!! Установить значение таймаута с conf.json
	*/
	client := &http.Client{
		Timeout: timeout * time.Second,
	}
	r, err := http.NewRequest("GET", ip, nil)
	if err != nil {
		log.Fatalln(err, `r, err := http.NewRequest("GET", ip, nil)`)
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

    resp, err := client.Do(r)
	if err != nil {
		StatusTimeout(err)
		return
	}

    code, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Println(err, "ioutil.ReadAll(resp.Body) Ошибка чтения body")
		return
    }

	switch resp.StatusCode {
		case 200:
			Status200(resp)
			AccountReg(resp.Request.Host, string(code))
		case 401:
			Status401(resp)
		case 404:
			Status404(resp)
		case 400:
			Status400(resp)
		default:
			StatusNotDefined(resp)
    }
}


/*
	Загрузка Devices.json и преоброзование в 
	DEVICES_OBJ map[string]Devices
*/
func ListDevicesInit() {

	devices, err := os.Open("Devices.json")
	if err != nil {
		log.Fatalln(
			LevelLog("F"),
			"ListDevices()",
			"Ошибка открытия файла:",
			err)
		return
	}
	defer devices.Close()

	byte_devices, err := ioutil.ReadAll(devices)
	if err != nil {
		log.Fatalln(
			LevelLog("F"),
			"ioutil.ReadAll(devices) Ошибка чтения Devices.json",
			err)
		return
	}

	if err = json.Unmarshal(byte_devices, &DEVICES_OBJ); err != nil {
		log.Fatalln(
			LevelLog("F"),
			"json.Unmarshal(byte_devices) Ошибка преоброзования в Devices struct",
			err)
		return
	}
}



 
















/*
	Восстановление оборудования после падения
	проверка основных параметров 
		Status: 200 (success)
		SIP:      AccountReg1=1
		Error:    10
		Label:    200 key
		ScanCode: 12345=>on
		DoorCode: 15926=>on
		Lock:     included

*/
func Recovery(){
	var Status 	 	string
	var SIP 	 	string
	var Error 	 	string
	var Label 	 	string
	var ScanCode  	string
	var DoorCode  	string
	var Lock 	 	string

	log.Println(Status, SIP, Error, Label, ScanCode, DoorCode, Lock)
}

/*
	Проверка состояния SIP сервера
*/
func AccountReg(ip, status string){
	//log.Println(status)

	statuses := strings.Split(status, "\n")

	if statuses[0] == "AccountReg1=1" {
		log.Println(LevelLog("I"), ip, "AccountReg1=1")
		return
	}

	if statuses[0] == "AccountReg1=0" {
		log.Println(LevelLog("W"), ip, "AccountReg1=0")
		return
	}

}
func StatusTimeout(err error){
	log.Println(LevelLog("E"), "Заданный узел не ответил в течении заданного {{timeout_second_request}}")
	log.Println(LevelLog("E"), err)
	return
}
func Status200(resp *http.Response){
	log.Println(
		LevelLog("I"),
		resp.Request.Host,
		resp.Status,
		"Succes",
	)
}
func Status401(resp *http.Response){
	log.Println(
		LevelLog("W"),
		resp.Request.Host,
		resp.Status,
		"401 Unauthorized - неверная авторизация для проверки статуса",
	)
}
func Status404(resp *http.Response){
	log.Println(
		LevelLog("W"),
		resp.Request.Host,
		resp.Status, 
		"404 Site or Page Not Found - указанная страница не существует",
	)
}
func Status400(resp *http.Response){
	log.Println(
		LevelLog("W"), 
		resp.Request.Host, 
		resp.Status, 
		"400 Bad Request - Неверный переданный запрос",
	)
}
func StatusNotDefined(resp *http.Response){
	log.Println(
		LevelLog("W"), 
		resp.Request.Host, 
		resp.Status, 
		"StatusNotDefined() - Обработка исключения",
	)
}

/*
	Уровень логирования
*/
func LevelLog(level string) string {
	switch level {
		case "F":
			return "<	/FATAL >"
		case "E":
			return "<	/ERROR >"
		case "W":
			return "< /WARNING >"
		case "I":
			return "<	 /INFO >"
		case "D":
			return "<	/DEBUG >"
		case "T":
			return "<	/TRACE >"
		default:
			return ""
	}
}

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