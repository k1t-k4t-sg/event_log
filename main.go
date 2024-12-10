/*
	Beward - CGI (GO)
	Regstatus
*/
package main

import (
	"fmt"
	"log"
	"io/ioutil"
	//"errors"
	"net/http"
    "time"
)

func main(){
	CheckingStatus("http://login:pass@192.168.1.249:27027/cgi-bin/sip_cgi?action=regstatus", 5)
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
		log.Fatalln(err, "timeout")
	}
    _, err = ioutil.ReadAll(resp.Body)

	/*
		!!! Добавить обработку ошибки чтения
	*/
    if err != nil {
        log.Fatalln(err, "ioutil.ReadAll(resp.Body) Ошибка чтения body")
    }

	switch resp.Status {
		// Succes
		case "200 OK":
			fmt.Println(resp.Status, "Succes")
		// 401 Unauthorized - неверная авторизация для проверки статуса
		case "401 Unauthorized":
			fmt.Println(resp.Status)
		// 404 Site or Page Not Found - указанная страница не существует
		case "404 Site or Page Not Found":
			fmt.Println(resp.Status)
		// 400 Bad Request - Неверный переданный запрос
		case "400 Bad Request":
			fmt.Println(resp.Status)
    }
}




 






















func StatusTimeout(){

}
func Status401(){

}
func Status404(){

}
func Status400(){

}
func StatusNotDefined(){
	
}