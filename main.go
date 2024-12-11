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
)

/*
	xType := fmt.Sprintf("%T", resp)
	fmt.Println(xType)

*/

func main(){
	for{
		time.Sleep(5 * time.Second)
		CheckingStatus(os.Getenv("REM_IP")+"/cgi-bin/sip_cgi?action=regstatus", 1)
	}
}

/*
	401 Unauthorized - неверная авторизация для проверки статуса
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

    _, err = ioutil.ReadAll(resp.Body)

	/*
		!!! Добавить обработку ошибки чтения
	*/
    if err != nil {
        log.Fatalln(err, "ioutil.ReadAll(resp.Body) Ошибка чтения body")
    }

	switch resp.StatusCode {
		case 200:
			Status200(resp)
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