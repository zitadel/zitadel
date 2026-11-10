package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
)

type serializableData struct {
	ContextInfo map[string]interface{} `json:"contextInfo,omitempty"`
	Args        map[string]interface{} `json:"args,omitempty"`
}

type response struct {
	Recipient string `json:"recipient,omitempty"`
}

func main() {
	port := flag.String("port", "3333", "used port for the sink")
	email := flag.String("email", "/email", "path for a sent email")
	emailKey := flag.String("email-key", "recipientEmailAddress", "value in the sent context info of the email used as key to retrieve the notification")
	sms := flag.String("sms", "/sms", "path for a sent sms")
	smsKey := flag.String("sms-key", "recipientPhoneNumber", "value in the sent context info of the sms used as key to retrieve the notification")
	notification := flag.String("notification", "/notification", "path to receive the notification")
	flag.Parse()

	messages := make(map[string]serializableData)

	http.HandleFunc(*email, func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		serializableData := serializableData{}
		if err := json.Unmarshal(data, &serializableData); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		email, ok := serializableData.ContextInfo[*emailKey].(string)
		if !ok {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println(email + ": " + string(data))
		messages[email] = serializableData
		io.WriteString(w, "Email!\n")
	})

	http.HandleFunc(*sms, func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		serializableData := serializableData{}
		if err := json.Unmarshal(data, &serializableData); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		phone, ok := serializableData.ContextInfo[*smsKey].(string)
		if !ok {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println(phone + ": " + string(data))
		messages[phone] = serializableData
		io.WriteString(w, "SMS!\n")
	})

	http.HandleFunc(*notification, func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response := response{}
		if err := json.Unmarshal(data, &response); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		msg, ok := messages[response.Recipient]
		if !ok {
			http.Error(w, "No messages found for recipient: "+response.Recipient, http.StatusNotFound)
			return
		}
		serializableData, err := json.Marshal(msg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, string(serializableData))
	})

	fmt.Println("Starting server on", *port)
	fmt.Println(*email, " for email handling")
	fmt.Println(*sms, " for sms handling")
	fmt.Println(*notification, " for retrieving notifications")
	http.Handle("/healthy", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { return }))
	fmt.Println("/healthy returns 200 OK")
	err := http.ListenAndServe(":"+*port, nil)
	if err != nil {
		panic("Server could not be started: " + err.Error())
	}
}
