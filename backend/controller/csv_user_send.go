package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

//UserFromCSV ...
func (s *server) UserFromCSV() http.HandlerFunc {
	type label struct {
		Email    int
		Password int
		Name     int
	}
	type User struct {
		Email    string
		Password string
		Name     string
		Okay     bool
	}
	type record struct {
		Accept  int
		Discard int
	}
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(32 << 20) // limit your max input length!
		var buf bytes.Buffer
		var size int
		u := User{}
		// In your case file would be fileupload
		file, _, err := r.FormFile("file")
		if err != nil {
			panic(err)
		}
		defer file.Close()
		// Copy the file data to my buffer
		io.Copy(&buf, file)
		// Buf to string
		contents := buf.String()
		//Create main array of string csv
		arr := strings.Split(contents, "\n")
		//Create dictionary labels
		labels := &label{
			Email:    -1,
			Password: -1,
			Name:     -1,
		}
		records := &record{
			Accept:  0,
			Discard: 0,
		}
		str := "{"
		for i := range arr {
			// Parse string of csv
			data := strings.Split(arr[i], ",")
			// Delete \r
			for j := range data {
				data[j] = strings.Replace(data[j], "\r", "", -1)
			}
			// First cycle for designation headers
			if i == 0 {
				size = len(data)
				for j := range data {
					fmt.Printf("=%s=", j)
					switch os := data[j]; os {
					case "Email":
						labels.Email = j
					case "Password":
						labels.Password = j
					case "Name":
						labels.Name = j
					}
				}
				continue
			}
			// If string length is less than specified..
			if len(data) < size {
				break
			}
			// Password check
			if labels.Password == -1 {
				u.Email = data[labels.Email]
				u.Password = "P@ssw0rd"
				u.Name = data[labels.Name]

			} else {
				u.Email = data[labels.Email]
				u.Password = data[labels.Password]
				u.Name = data[labels.Name]
			}
			// If user was create
			if x, err := s.store.User().FindByEmail(u.Email); err == nil {
				fmt.Println("Dublicate")
				fmt.Println(x)
				records.Discard++
				u.Okay = false
			} else {
				fmt.Println(err)
				u.Okay = true
				records.Accept++
			}

			b, err := json.Marshal(u)
			if err != nil {
				fmt.Println(err)
				return
			}
			if i > 1 {
				str += ",\n"
			}
			str += "\"" + strconv.Itoa(i) + "\":" + string(b)
		}
		str += "\n}"
		w.Write([]byte(fmt.Sprintf(str)))
		// Reset the buffer in case I want to use it again
		// Reduces memory allocations in more intense projects
		buf.Reset()
		//w.Write([]byte(fmt.Sprintf("{\"status\":\"OK\",\"accept\":\"%v\",\"discard\":\"%v\"}", records.Accept, records.Discard)))
		s.respond(w, r, http.StatusOK, nil)
	}
}