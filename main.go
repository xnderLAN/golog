package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	username = os.Getenv("GOLOG_USERNAME")
	password = os.Getenv("GOLOG_PASSWORD")
)

type DataSchema struct {
	ID   int32  `json:"id"`
	Data string `json:"data"`
}

type DataResponse struct {
	Error bool `json:"err"`
}

type LogResponse struct {
	Error bool   `json:"err"`
	Data  string `json:"data"`
}

type LogRequest struct {
	ID int32 `json:"id"`
}

func log_d(w http.ResponseWriter, r *http.Request) {
	var p DataSchema

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	log.Printf("data: %v", p)
	if p.Data == "" || p.ID == 0 {
		Response := []DataResponse{
			{Error: true},
		}
		b, err := json.MarshalIndent(Response, "", "  ")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
	write_resulte := write_log(p.Data, int(p.ID))
	if write_resulte {
		Response := []DataResponse{
			{Error: false},
		}
		b, err := json.MarshalIndent(Response, "", "  ")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}

func write_log(data string, id int) bool {
	tm := time.Now()
	if id != 0 && len(data) > 0 {
		path := "/opt/golog/" + strconv.Itoa(id) + ".log"
		f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			fmt.Println(err)
			return false
		}
		defer f.Close()
		s := tm.String() + " | " + data + "\n"
		if _, err = f.WriteString(s); err != nil {
			fmt.Println(err)
			return false
		}
		return true
	}
	return false
}

func readlog(w http.ResponseWriter, r *http.Request) {
	u, p, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(401)
		return
	}
	fmt.Printf("username: %v Password: %v\n", username, password)
	if u != username || p != password {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.WriteHeader(401)
		return
	}
	var ps DataSchema

	err := json.NewDecoder(r.Body).Decode(&ps)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.WriteHeader(http.StatusBadRequest)

		fmt.Println(err)
		return
	}
	fmt.Printf("data: %v\n", ps)
	data := readlog_file(int(ps.ID))
	if len(data) > 0 {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, data)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func readlog_file(id int) string {
	path := "/opt/golog/" + strconv.Itoa(id) + ".log"

	file_data, err := os.ReadFile(path)

	if err != nil {
		fmt.Printf("Error: %v", err)
		return ""
	}
	return string(file_data)
}

func index(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "why?")

}

func log_c(w http.ResponseWriter, r *http.Request) {

	var p DataSchema

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	log.Printf("data: %v", p)
	if p.Data == "" || p.ID == 0 {
		Response := []DataResponse{
			{Error: true},
		}
		b, err := json.MarshalIndent(Response, "", "  ")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
	write_resulte := write_log(p.Data, int(p.ID))
	if write_resulte {
		Response := []DataResponse{
			{Error: false},
		}
		b, err := json.MarshalIndent(Response, "", "  ")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}

func main() {
	http.HandleFunc("/readlog/", readlog)
	http.HandleFunc("/golog/", log_d)
	http.HandleFunc("/clog/", log_c)
	http.HandleFunc("/", index)

	if err := http.ListenAndServe("127.0.0.1:8086", nil); err != nil {
		log.Fatal(err)
	}
}
