package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"../model"
)

//UserFromCSV ...
func (s *server) StandFromCsv() http.HandlerFunc {
	type label struct {
		Email      int
		Address    int
		Digipass   int
		Datacenter int
		Digiuser   int
		Esxipass   int
		Esxiuser   int
		Module     int
		Digi       int
		PortT      int
	}
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(64 << 20) // limit your max input length!
		var buf bytes.Buffer
		var size int
		// in your case file would be fileupload
		file, _, err := r.FormFile("file")
		if err != nil {
			panic(err)
		}
		defer file.Close()
		// Copy the file data to my buffer
		io.Copy(&buf, file)
		contents := buf.String()
		arr := strings.Split(contents, "\n")

		labels := &label{}
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
					case "Address":
						labels.Address = j
					case "Digipass":
						labels.Digipass = j
					case "Datacenter":
						labels.Datacenter = j
					case "Digiuser":
						labels.Digiuser = j
					case "Esxipass":
						labels.Esxipass = j
					case "Esxiuser":
						labels.Esxiuser = j
					case "Module":
						labels.Module = j
					case "Digi":
						labels.Digi = j
					case "Port":
						labels.PortT = j
					}
				}
				continue
			}
			// If string length is less than specified..
			if len(data) < size {
				break
			}
			u := &model.Stand{
				Email:      data[labels.Email],
				Address:    data[labels.Address],
				Digipass:   data[labels.Digipass],
				Datacenter: data[labels.Datacenter],
				Digiuser:   data[labels.Digiuser],
				Esxipass:   data[labels.Esxipass],
				Esxiuser:   data[labels.Esxiuser],
				Module:     data[labels.Module],
				Digi:       data[labels.Digi],
				PortT:      data[labels.PortT],
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
		// I reset the buffer in case I want to use it again
		// reduces memory allocations in more intense projects
		buf.Reset()
		s.respond(w, r, http.StatusOK, nil)
	}
}
