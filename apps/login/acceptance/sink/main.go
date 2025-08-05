package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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
	configureZitadel := flag.Bool("configure-zitadel", false, "if set, the sink will configure the Zitadel instance with the given email and sms paths")
	zitadelAPIUrl := flag.String("zitadel-api-url", "http://localhost:8080", "Zitadel API URL to configure the sink")
	zitadelExternalDomain := flag.String("zitadel-external-domain", "localhost", "Zitadel external domain to configure the sink")
	zitadelAPITokenFile := flag.String("zitadel-api-token-file", "", "File containing the Zitadel API token to configure the sink")
	mockServiceURL := flag.String("mock-service-url", "http://localhost:3333", "URL of the mock service to be used in tests")
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
		w.Write([]byte("Email!\n"))
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
		w.Write([]byte("SMS!\n"))
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
		w.Write(serializableData)
	})
	if *configureZitadel {
		zitadelAPIToken, err := os.ReadFile(*zitadelAPITokenFile)
		if err != nil {
			panic("Could not read Zitadel API token file: " + err.Error())
		}
		cleanToken := strings.TrimSpace(string(zitadelAPIToken))
		ensureProvider(*zitadelAPIUrl, cleanToken, *zitadelExternalDomain, *mockServiceURL, *email)
		ensureProvider(*zitadelAPIUrl, cleanToken, *zitadelExternalDomain, *mockServiceURL, *sms)
	}

	fmt.Println("Starting server on", *port)
	fmt.Println(*email, " for email handling")
	fmt.Println(*sms, " for sms handling")
	fmt.Println(*notification, " for retrieving notifications")
	http.Handle("/healthy", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	fmt.Println("/healthy returns 200 OK")
	err := http.ListenAndServe(":"+*port, nil)
	if err != nil {
		panic("Server could not be started: " + err.Error())
	}
}

func ensureProvider(zitadelAPIUrl string, zitadelAPIToken string, zitadelAPIExternalDomain string, mockServiceUrl string, path string) {
	fmt.Println("Ensuring Zitadel provider for", path)
	ensureProviderURL := fmt.Sprintf("%s/admin/v1%s/http", zitadelAPIUrl, path)
	payload := "{\"endpoint\": \"" + mockServiceUrl + path + "\", \"description\": \"test\"}"
	newProvider := &struct {
		ID string `json:"id"`
	}{}
	header := map[string]string{
		"Authorization": "Bearer " + zitadelAPIToken,
	}
	post(ensureProviderURL, zitadelAPIExternalDomain, header, payload, newProvider)
	activateProviderURL := fmt.Sprintf("%s/admin/v1%s/%s/_activate", zitadelAPIUrl, path, newProvider.ID)
	post(activateProviderURL, zitadelAPIExternalDomain, header, payload, nil)
}

func post(url string, host string, header map[string]string, payload string, parseResponse any) {
	fmt.Println("POSTing to", url)
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(payload))
	if err != nil {
		panic("Could not create request: " + err.Error())
	}
	req.Host = host
	for key, value := range header {
		req.Header[key] = []string{value}
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic("Could not configure Zitadel: " + err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		errorResp, err := io.ReadAll(resp.Body)
		if err != nil {
			panic("Could not read error response from Zitadel: " + err.Error())
		}
		panic(fmt.Sprintf("Zitadel configuration failed with status %d: %s, request url: %s, request headers: %+v", resp.StatusCode, string(errorResp), req.URL, req.Header))
	}
	if parseResponse == nil {
		return
	}
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		panic("Could not read response from Zitadel: " + err.Error())
	}
	if err := json.Unmarshal(response, parseResponse); err != nil {
		panic("Could not parse response from Zitadel: " + err.Error())
	}
	fmt.Println("Zitadel response:", string(response))
}
