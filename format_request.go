package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func formatMessageHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var message map[string]interface{}
	err := decoder.Decode(&message)
	if err != nil {
		fmt.Println("error while decoding: ", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Create a channel to signal completion
	done := make(chan struct{})

	// Run the message formatting in a goroutine
	go func() {
		defer close(done)

		formattedMessage := formatMessage(message)

		// Marshal the formatted message into JSON
		response, err := json.Marshal(formattedMessage)
		if err != nil {
			fmt.Println("Error marshaling response:", err)
			return
		}

		// webhookURL := "https://webhook.site/" // Replace with your actual webhook URL
		// resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(response))
		// if err != nil {
		// 	fmt.Println("Error sending data to webhook:", err)
		// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		// 	return
		// }
		// defer resp.Body.Close()

		// Check the response status from the webhook
		// if resp.StatusCode != http.StatusOK {
		// 	fmt.Println("Webhook returned non-OK status:", resp.Status)
		// 	http.Error(w, "Webhook Error", http.StatusInternalServerError)
		// 	return
		// }

		// Send the JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}()

	// Wait for the goroutine to finish before responding
	<-done
}

func main() {
	http.HandleFunc("/format-message", formatMessageHandler)
	fmt.Println("listening on port:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func formatMessage(message map[string]interface{}) map[string]interface{} {
	formattedMessage := make(map[string]interface{})

	formattedMessage["event"] = message["ev"]
	formattedMessage["event_type"] = message["et"]
	formattedMessage["app_id"] = message["id"]
	formattedMessage["user_id"] = message["uid"]
	formattedMessage["message_id"] = message["mid"]
	formattedMessage["page_title"] = message["t"]
	formattedMessage["page_url"] = message["p"]
	formattedMessage["browser_language"] = message["l"]
	formattedMessage["screen_size"] = message["sc"]

	formattedMessage["attributes"] = map[string]interface{}{}
	formattedMessage["traits"] = map[string]interface{}{}

	for key, value := range message {
		if len(key) >= 5 {
			if key[:4] == "atrk" {
				attrKey := key[4:]
				attrType := "atrt" + attrKey
				attrVal := "atrv" + attrKey
				formattedMessage["attributes"].(map[string]interface{})[value.(string)] = map[string]interface{}{
					"value": message[attrVal],
					"type":  message[attrType],
				}
			} else if key[:5] == "uatrk" {
				traitKey := key[5:]
				traitType := "uatrt" + traitKey
				traitVal := "uatrv" + traitKey
				formattedMessage["traits"].(map[string]interface{})[value.(string)] = map[string]interface{}{
					"value": message[traitVal],
					"type":  message[traitType],
				}
			}
		}
	}

	return formattedMessage
}
