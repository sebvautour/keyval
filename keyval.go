package keyval

import (
	"errors"
	"strings"
)

// Request struct is used to pass to the DB
type Request struct {
	Op  string
	Key string
	Val interface{}
}

// Response struct
type Response struct {
	Error error
	Found bool
	Val   interface{}
}

// DB struct is the main structure that has all the methods attached to it
type DB struct {
	Req chan Request
	Res chan Response
}

// Get retrieves a value from the DB
func (db DB) Get(key string) (value interface{}, found bool) {
	db.Req <- Request{"Get", key, nil}
	r := <-db.Res
	return r.Val, r.Found
}

// Set places a key/val in the DB
func (db DB) Set(key string, value interface{}) {
	db.Req <- Request{"Set", key, value}
	_ = <-db.Res
}

// Del removes a key/val pair from the DB
func (db DB) Del(key string) {
	db.Req <- Request{"Del", key, nil}
	_ = <-db.Res
}

// Launch function starts the db as a goroutine
func (db DB) Launch() {

	data := (map[string]interface{}{})
	for r := range db.Req {
		switch strings.ToUpper(r.Op) {
		case "SET":
			data[r.Key] = r.Val
			db.Res <- Response{
				Error: nil,
				Found: true,
				Val:   nil,
			}
		case "GET":
			i, ok := data[r.Key]
			if ok {
				db.Res <- Response{
					Error: nil,
					Found: true,
					Val:   i,
				}
			} else {
				db.Res <- Response{
					Error: errors.New("Key not found"),
					Found: false,
					Val:   i,
				}
			}
		case "DEL":
			delete(data, r.Key)
			db.Res <- Response{
				Error: nil,
				Found: true,
				Val:   nil,
			}
		default:
			db.Res <- Response{
				Error: errors.New("Unknown method"),
				Found: false,
				Val:   nil,
			}
		}
	}

}

// Init func creates and initalizes a DB struct
func Init() *DB {
	var db DB
	db.Req = make(chan Request)
	db.Res = make(chan Response)
	return &db
}
