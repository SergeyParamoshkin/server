package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/gorilla/mux"
)

type Response struct {
	Status int
	Value  string
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

func TodoIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		panic(err)
	}
}

func TodoShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var todoId int
	var err error
	if todoId, err = strconv.Atoi(vars["todoId"]); err != nil {
		panic(err)
	}
	todo := RepoFindTodo(todoId)
	if todo.Id > 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(todo); err != nil {
			panic(err)
		}
		return
	}

	// If we didn't find it, 404
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
		panic(err)
	}

}

/*
Test with this curl command:

curl -H "Content-Type: application/json" -d '{"name":"New Todo"}' http://localhost:8080/todos

*/
func TodoCreate(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &todo); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	t := RepoCreateTodo(todo)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		panic(err)
	}
}

/*
func CreateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if len(vars) == 0 {
		log.Println("bad request!!!", "http request:[", r, "]")
		if err := json.NewEncoder(w).Encode(Response{http.StatusBadRequest, ""}); err != nil {
			panic(err)
		}
		return
	}
	log.Println("userName:[", vars["userName"], "]")
	out, err := exec.Command("env").Output()
	if err != nil {
		log.Fatal(err.Error())
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(Response{http.StatusCreated, string(out)}); err != nil {
		panic(err)
	}
}
*/
func CreateUser(w http.ResponseWriter, r *http.Request) {
	type User struct {
		Name    string
		Comment string
		Shell   string
	}
	var (
		err       error
		out, body []byte
		user      User
	)
	w.Header().Set("Content-Type", ContentTypeJson)

	body, err = ioutil.ReadAll(io.LimitReader(r.Body, LimitReader))
	if err != nil {
		log.Println("ioutil.ReadAll error:[", err.Error(), "]")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = r.Body.Close(); err != nil {
		log.Println("Body.Close error:[", err.Error(), "]")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(body, &user); err != nil {
		log.Println("bad request!!!", "http request:[", r, "]", "body:[", string(body), "]", "error:[", err.Error(), "]")
		w.WriteHeader(http.StatusBadRequest)
		if err = json.NewEncoder(w).Encode(Response{http.StatusBadRequest, err.Error()}); err != nil {
			log.Println(err.Error())
			return
		}
		return
	} else {
		if len(user.Name) > 0 {
			out, err = exec.Command("/bin/bash", "./plugins/useradd/useradd.sh", user.Name).Output()
			if err != nil {
				log.Fatal(err.Error())
			}
			w.WriteHeader(http.StatusCreated)
			if err := json.NewEncoder(w).Encode(Response{http.StatusCreated, string(out)}); err != nil {
				log.Println(err.Error())
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			if err = json.NewEncoder(w).Encode(Response{http.StatusBadRequest, "user name empty"}); err != nil {
				log.Println(err.Error())
			}
		}
	}
}
