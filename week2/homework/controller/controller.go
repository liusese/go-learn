package controller

import (
	"go-learn/week2/homework/dao"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/pquerna/ffjson/ffjson"
)

const (
	driver     = "mysql"
	dataSource = "root:123456@tcp(10.10.50.198:13306)/test?charset=utf8&loc=Local"
)

type Msg struct {
	ErrCode int         `json:"err_code"`
	ErrMsg  string      `json:"err_msg"`
	Data    interface{} `json:"data"`
}

func AddPerson(w http.ResponseWriter, r *http.Request) {
	var person dao.Person

	bBody, _ := ioutil.ReadAll(r.Body)
	_ = ffjson.Unmarshal(bBody, &person)
	// TODO: validator
	w.Header().Set("content-type", "application/json")
	conn, err := dao.TryConn(driver, dataSource)
	if nil != err {
		w.WriteHeader(500)
		msg, _ := ffjson.Marshal(&Msg{
			ErrCode: 12345,
			ErrMsg:  err.Error(),
		})
		w.Write(msg)
		return
	}
	defer conn.Close()

	w.WriteHeader(200)
	_, err = dao.Add(conn, &person)
	if nil != err {
		msg, _ := ffjson.Marshal(&Msg{
			ErrCode: 54321,
			ErrMsg:  err.Error(),
		})
		w.Write(msg)
	} else {
		msg, _ := ffjson.Marshal(&Msg{
			ErrCode: 0,
			ErrMsg:  `OK`,
		})
		w.Write(msg)
	}

	return
}

func DelPerson(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	idStr := v.Get("id")
	id, _ := strconv.Atoi(idStr)
	w.Header().Set("content-type", "application/json")
	conn, err := dao.TryConn(driver, dataSource)
	if nil != err {
		w.WriteHeader(500)
		msg, _ := ffjson.Marshal(&Msg{
			ErrCode: 12345,
			ErrMsg:  err.Error(),
		})
		w.Write(msg)
		return
	}
	defer conn.Close()

	w.WriteHeader(200)
	_, err = dao.Delete(conn, id)
	if nil != err {
		msg, _ := ffjson.Marshal(&Msg{
			ErrCode: 54321,
			ErrMsg:  err.Error(),
		})
		w.Write(msg)
	} else {
		msg, _ := ffjson.Marshal(&Msg{
			ErrCode: 0,
			ErrMsg:  `OK`,
		})
		w.Write(msg)
	}

	return
}

func UpdatePerson(w http.ResponseWriter, r *http.Request) {
	var person dao.Person

	bBody, _ := ioutil.ReadAll(r.Body)
	_ = ffjson.Unmarshal(bBody, &person)
	w.Header().Set("content-type", "application/json")
	conn, err := dao.TryConn(driver, dataSource)
	if nil != err {
		w.WriteHeader(500)
		msg, _ := ffjson.Marshal(&Msg{
			ErrCode: 12345,
			ErrMsg:  err.Error(),
		})
		w.Write(msg)
		return
	}
	defer conn.Close()

	w.WriteHeader(200)
	_, err = dao.Update(conn, &person)
	if nil != err {
		msg, _ := ffjson.Marshal(&Msg{
			ErrCode: 54321,
			ErrMsg:  err.Error(),
		})
		w.Write(msg)
	} else {
		msg, _ := ffjson.Marshal(&Msg{
			ErrCode: 0,
			ErrMsg:  `OK`,
		})
		w.Write(msg)
	}

	return
}

func FetchPerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	conn, err := dao.TryConn(driver, dataSource)
	if nil != err {
		w.WriteHeader(500)
		msg, _ := ffjson.Marshal(&Msg{
			ErrCode: 12345,
			ErrMsg:  err.Error(),
		})
		w.Write(msg)
		return
	}
	defer conn.Close()

	w.WriteHeader(200)
	l, num, err := dao.Fetch(conn, map[string]interface{}{
		"gender": "female",
	})
	if nil != err {
		msg, _ := ffjson.Marshal(&Msg{
			ErrCode: 54321,
			ErrMsg:  err.Error(),
		})
		w.Write(msg)
	} else {
		msg, _ := ffjson.Marshal(&Msg{
			ErrCode: 0,
			ErrMsg:  `OK`,
			Data: map[string]interface{}{
				"num":     num,
				"persons": l,
			},
		})
		w.Write(msg)
	}

	return
}
