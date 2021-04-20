package main

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	driver     = "mysql"
	dataSource = "root:123456@tcp(10.10.50.198:13306)/test?charset=utf8&loc=Local"
)

func TestConnection(t *testing.T) {
	Convey("test conn", t, func() {
		conn, err := tryConn(driver, dataSource)
		So(conn, ShouldNotBeNil)
		So(err, ShouldBeNil)
		fmt.Printf("%+v\n", conn)
		if nil != conn {
			conn.Close()
		}
	})
}

func TestRead(t *testing.T) {
	Convey("test read", t, func() {
		conn, err := tryConn(driver, dataSource)
		So(conn, ShouldNotBeNil)
		So(err, ShouldBeNil)
		fmt.Printf("%+v\n", conn)
		defer conn.Close()

		p := Person{
			Name: "Lucy",
		}

		err = Read(conn, &p, "name", "gender")
		So(err, ShouldBeNil)
		fmt.Printf("%+v\n", p)
	})
}

func TestAdd(t *testing.T) {
	Convey("test add", t, func() {
		conn, err := tryConn(driver, dataSource)
		So(conn, ShouldNotBeNil)
		So(err, ShouldBeNil)
		fmt.Printf("%+v\n", conn)
		defer conn.Close()

		p := Person{
			Name:   "Jim",
			Gender: "male",
		}

		num, err := Add(conn, &p)
		// _ = num
		So(num, ShouldEqual, 1)
		So(err, ShouldBeNil)
		fmt.Printf("%+v\n", p)
	})
}

func TestDel(t *testing.T) {
	Convey("test delete", t, func() {
		conn, err := tryConn(driver, dataSource)
		So(conn, ShouldNotBeNil)
		So(err, ShouldBeNil)
		fmt.Printf("%+v\n", conn)
		defer conn.Close()

		num, err := Delete(conn, 4)
		_ = num
		// So(num, ShouldEqual, 1)
		So(err, ShouldBeNil)
	})
}

func TestUpdate(t *testing.T) {
	Convey("test update", t, func() {
		conn, err := tryConn(driver, dataSource)
		So(conn, ShouldNotBeNil)
		So(err, ShouldBeNil)
		fmt.Printf("%+v\n", conn)
		defer conn.Close()

		p := Person{
			ID:     10,
			Name:   "Lucy001",
			Gender: "female",
		}

		num, err := Update(conn, &p)
		_ = num
		// So(num, ShouldEqual, 1)
		So(err, ShouldBeNil)
	})
}

func TestFetch(t *testing.T) {
	Convey("test fetch", t, func() {
		conn, err := tryConn(driver, dataSource)
		So(conn, ShouldNotBeNil)
		So(err, ShouldBeNil)
		fmt.Printf("%+v\n", conn)
		defer conn.Close()

		query := map[string]interface{}{
			"gender": "female",
		}

		l, num, err := Fetch(conn, query)
		// So(num, ShouldEqual, 1)
		So(err, ShouldBeNil)
		fmt.Printf("%+v\n%#v\n", num, l)
	})
}
