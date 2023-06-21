package constants

import (
	"log"
	"os"
)

var (
	STORECREDENTIALSFILE = ""
	STOREBUCKETNAME      = ""
)

func InitConstants() {

	storeCredentialsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if len(storeCredentialsFile) > 0 {
		STORECREDENTIALSFILE = storeCredentialsFile
	}
	storeBucketName := os.Getenv("STORE_BUCKET")
	if len(storeBucketName) > 0 {
		STOREBUCKETNAME = storeBucketName
	}

	log.Println("GOOGLE_APPLICATION_CREDENTIALS: ", STORECREDENTIALSFILE)
}
