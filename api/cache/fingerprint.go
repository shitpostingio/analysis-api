package cache

import (
	"encoding/json"
	"fmt"
	"github.com/shitpostingio/analysis-commons/structs"
	"log"
)

// PutFingerprint puts a FingerprintResponse in cache.
func PutFingerprint(key string, data structs.FingerprintResponse) error {

	bytes, err := json.Marshal(&data)
	if err != nil {
		log.Println("PutFingerprint: unable to marshal data ", err)
		return err
	}

	return fClient.Set(key, string(bytes), 0).Err()

}

// GetFingerprint gets a FingerprintResponse from the cache.
func GetFingerprint(key string) (data structs.FingerprintResponse, err error) {

	if fClient.Exists(key).Val() == 0 {
		return data, fmt.Errorf("GetFingerprint: Key %s not found", key)
	}

	result, err := fClient.Get(key).Result()
	if err != nil {
		log.Println(fmt.Sprintf("GetFingerprint: Unable to perform GET for key %s: %s", key, err))
		return
	}

	err = json.Unmarshal([]byte(result), &data)
	return

}
