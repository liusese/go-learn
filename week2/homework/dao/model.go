/*
我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？
*/
package dao

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/pquerna/ffjson/ffjson"
)

type Person struct {
	ID     int    `json:"id"`     // pk
	Name   string `json:"name"`   // not null  unique  validator require
	Gender string `json:"gender"` // not null
}

func (p Person) String() string {
	bPerson, _ := ffjson.Marshal(&p)
	return string(bPerson)
}

// Connection
func TryConn(driver, source string) (*sql.DB, error) {
	return sql.Open(driver, source)
}

// Read
func Read(conn *sql.DB, person *Person, columns ...string) (err error) {
	var (
		querySql   strings.Builder
		args       []interface{}
		conditions []string
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
	var (
		result  sql.Result
		execSql string
		args    []interface{}
	)

	if err = Read(conn, p, "name"); nil != err && err == sql.ErrNoRows {
		execSql = "INSERT INTO person(name, gender) VALUES (?, ?)"
		args = append(args, (*p).Name, (*p).Gender)
		fmt.Println(execSql, args)

		result, err = conn.Exec(execSql, args...)
		if nil == err {
			num, err = result.RowsAffected()
		}
	} else if nil == err {
		err = errors.New("Errors: person alread exist")
	}

	if nil != err {
		fmt.Println(errors.Errorf("Add person failure, sql(%+v), args(%+v), err(%+v)", execSql, args, err))
		err = errors.New("Errors: unexpect error")
	}

	return
}

// Delete
func Delete(conn *sql.DB, id int) (num int64, err error) {
	var (
		result  sql.Result
		execSql string
		args    []interface{}
	)

	p := Person{ID: id}
	if err = Read(conn, &p); nil == err {
		execSql = "DELETE FROM person WHERE id=?"
		args = append(args, id)
		fmt.Println(execSql, args)

		result, err = conn.Exec(execSql, args...)
		if nil == err {
			num, _ = result.RowsAffected()
		}
	}

	if nil != err {
		if err == sql.ErrNoRows {
			err = errors.New("Errors: no data")
		} else {
			fmt.Println(errors.Errorf("Del person failure, sql(%+v), args(%+v), err(%+v)", execSql, args, err))
			err = errors.Errorf("Errors: unexpect error")
		}
	}

	return
}

// Update
func Update(conn *sql.DB, p *Person) (num int64, err error) {
	var (
		result  sql.Result
		execSql string
		args    []interface{}
	)

	m := Person{ID: (*p).ID}
	if err = Read(conn, &m); nil == err {
		execSql = "UPDATE person SET name=?, gender=? WHERE id=?"
		args = append(args, p.Name, p.Gender, m.ID)
		fmt.Println(execSql, args)

		result, err = conn.Exec(execSql, args...)
		if nil == err {
			num, err = result.RowsAffected()
		}
	}

	if nil != err {
		if err == sql.ErrNoRows {
			err = errors.New("Errors: no data")
		} else {
			fmt.Println(errors.Errorf("Update person failure, sql(%+v), args(%+v), err(%+v)", execSql, args, err))
			err = errors.Errorf("Errors: unexpect error")
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
			var person Person
			_ = rows.Scan(&person.ID, &person.Name, &person.Gender)
			l = append(l, &person)
		}

		num = int64(len(l))
	} else {
		if err == sql.ErrNoRows {
			err = errors.New("Errors: no data")
		} else {
			fmt.Println(errors.Errorf("Fetch person failure, sql(%+v), args(%+v), err(%+v)", querySql, args, err))
			err = errors.Errorf("Errors: unexpect error")
		}
	}

	return
}
