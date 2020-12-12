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
func (s *server) AddToChampCSV() http.HandlerFunc {
	type label struct {
		Email  int
		Module int
		Moduls int
	}
	type champ struct {
		Email  string
		Module string
		Time   []string
		Okay   bool
	}
	type record struct {
		Accept  int
		Discard int
	}
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(64 << 20) // Limit your max input length!
		var buf bytes.Buffer
		var size int
		u := &champ{}
		champ := r.FormValue("champ")
		file, _, err := r.FormFile("file")
		if err != nil {
			panic(err)
		}
		defer file.Close()
		// Copy the file data to my buffer
		io.Copy(&buf, file)
		contents := buf.String()
		arr := strings.Split(contents, "\n")

		labels := &label{
			Email:  -1,
			Module: -1,
		}
		records := &record{
			Accept:  0,
			Discard: 0,
		}
		str := "{"
		for i := range arr {
			data := strings.Split(arr[i], ",")
			for j := range data {
				data[j] = strings.Replace(data[j], "\r", "", -1)
			}
			if i == 0 {
				size = len(data)
				for j := range data {
					switch os := data[j]; os {
					case "Email":
						labels.Email = j
					case "Module":
						labels.Module = j
					}
				}
				continue
			}
			// If string length is less than specified..
			if len(data) < size {
				break
			}
			//mod := "2020-11-05T21:25:20+05:00;2020-11-05T21:00:20+05:00;2020-11-05T21:25:20+05:00;2020-11-06T23:25:20+05:00%2020-11-05T21:25:20+05:00;2020-11-05T21:00:20+05:00;2020-11-05T21:25:20+05:00;2020-11-06T23:25:20+05:00%2020-11-05T21:25:20+05:00;2020-11-05T21:00:20+05:00;2020-11-05T21:25:20+05:00;2020-11-06T23:25:20+05:00"
			var mod []string
			mod = append(mod, "2020-11-05T21:25:20+05:00;2020-11-05T21:00:20+05:00;2020-11-05T21:25:20+05:00;2020-11-06T23:25:20+05:00")
			mod = append(mod, "2020-11-05T21:25:20+05:00;2020-11-05T21:00:20+05:00;2020-11-05T21:25:20+05:00;2020-11-06T23:25:20+05:00")
			mod = append(mod, "2020-11-05T21:25:20+05:00;2020-11-05T21:00:20+05:00;2020-11-05T21:25:20+05:00;2020-11-06T23:25:20+05:00")
			if labels.Module == -1 {
				u.Email = data[labels.Email]
				u.Module = "A,B,C"
				u.Time = mod
			} else {
				u.Email = data[labels.Email]
				u.Module = strings.Replace(data[labels.Module], ";", ",", -1)
				u.Time = mod
			}
			fmt.Println(u)
			err = s.store.User().FindByEmailChamp(u.Email, champ)
			if err == nil {
				fmt.Println("Duplicate")
				records.Discard++
				u.Okay = false
			} else {
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
