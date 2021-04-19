/*
我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？
*/
package main

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Person struct {
	ID     int    `json:"id"`     // pk
	Name   string `json:"name"`   // not null  unique
	Gender string `json:"gender"` // not null
}

// Connection
func tryConn(driver, source string) (*sql.DB, error) {
	return sql.Open(driver, source)
}

// Read
func Read(conn *sql.DB, person *Person, columns ...string) (err error) {
	var (
		querySql   strings.Builder
		conditions []string
		args       []interface{}
	)

	if person == nil {
		err = sql.ErrNoRows
	} else {
		if len(columns) == 0 {
			querySql.WriteString("SELECT * FROM person WHERE id=?")
			args = append(args, person.ID)
		} else {
			querySql.WriteString("SELECT * FROM person WHERE ")
			ctype := reflect.TypeOf(*person)
			cValue := reflect.ValueOf(*person)
			for _, column := range columns {
				for i := 0; i < ctype.NumField(); i++ {
					if strings.EqualFold(strings.ToLower(ctype.Field(i).Name), strings.ToLower(column)) {
						conditions = append(conditions, fmt.Sprintf("%s=?", strings.ToLower(column)))
						args = append(args, cValue.Field(i).Interface())
					}
				}
			}
		}
		querySql.WriteString(strings.Join(conditions, " AND "))
		fmt.Println(querySql.String(), args)

		row := conn.QueryRow(querySql.String(), args...)
		if nil != row.Err() {
			err = row.Err()
		} else {
			err = row.Scan(&person.ID, &person.Name, &person.Gender)
		}
	}

	return
}

// Add
func Add(conn *sql.DB, p *Person) (num int64, err error) {
	var result sql.Result

	if err = Read(conn, p); nil != err && err == sql.ErrNoRows {
		result, err = conn.Exec("INSERT INTO person(name, gender) VALUES (?, ?)", (*p).Name, (*p).Gender)
		if nil == err {
			num, err = result.RowsAffected()
		}
	} else {
		log.Println("data already exist")
	}

	return
}

// Delete
func Delete(conn *sql.DB, id int) (num int64, err error) {
	var result sql.Result

	p := Person{ID: id}
	if err = Read(conn, &p); nil == err {
		result, err = conn.Exec("DELETE FROM person WHERE id=?", id)
		if nil == err {
			num, _ = result.RowsAffected()
		}
	}

	return
}

// Update
func Update(conn *sql.DB, p *Person) (num int64, err error) {
	var result sql.Result

	m := Person{ID: (*p).ID}
	if err = Read(conn, &m); nil == err {
		result, err = conn.Exec("UPDATE person SET name=?, gender=? WHERE id=?", p.Name, p.Gender, m.ID)
		if nil == err {
			num, err = result.RowsAffected()
		}
	}

	return
}

// Fetch
func Fetch(conn *sql.DB, query map[string]interface{}) (l []*Person, num int64, err error) {
	var (
		querySql   strings.Builder
		conditions []string
		args       []interface{}
	)

	querySql.WriteString("SELECT * FROM person WHERE ")
	for k, v := range query {
		conditions = append(conditions, fmt.Sprintf("%s=?", k))
		args = append(args, v)
	}
	querySql.WriteString(strings.Join(conditions, " AND "))
	fmt.Println(querySql.String(), args)

	rows, err := conn.Query(querySql.String(), args...)
	if nil == err {
		defer rows.Close()

		for rows.Next() {
			var (
				id     int
				name   string
				gender string
			)
			rows.Scan(&id, &name, &gender)
			l = append(l, &Person{
				ID:     id,
				Name:   name,
				Gender: gender,
			})
		}

		num = int64(len(l))
	}

	return
}