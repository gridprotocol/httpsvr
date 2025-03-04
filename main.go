package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

var CmdName string

func main() {
	// info
	http.HandleFunc("/", handlerInfo)
	// alter payee
	http.HandleFunc("/alterpayee", handlerAlterPayee)

	fmt.Println("Starting server on :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}

// alter payee to a grid user address
func handlerAlterPayee(w http.ResponseWriter, r *http.Request) {

	// 获取查询参数中的 user 值
	user := r.URL.Query().Get("user")
	if user == "" {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing 'user' query parameter"))
		return
	}

	fmt.Println("grid user address:", user)

	// check user or provider
	cmdName, err := UserOrProvider()
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	// alter payee cmd arg
	arg := "settle alterPayee --really-do-it=true " + user

	cmd := exec.Command(cmdName, arg)
	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(output))
}

// info
func handlerInfo(w http.ResponseWriter, r *http.Request) {
	// check user or provider
	cmdName, err := UserOrProvider()
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	cmd := exec.Command(cmdName, "info")
	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(output))
}

// check mefs-user or mefs-provider
func UserOrProvider() (string, error) {
	e, err := FileExists("./mefs-user")
	if err != nil {
		return "", err
	}

	// return user or provider
	if e {
		return "mefs-user", nil
	} else {
		return "mefs-provider", nil
	}

}

// 文件是否存在
func FileExists(filename string) (bool, error) {
	_, err := os.Stat(filename)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
