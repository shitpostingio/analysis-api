package cache

import (
	"encoding/json"
	"fmt"
	"github.com/shitpostingio/analysis-commons/structs"
	"log"
)

// PutNSFW puts a NSFWResponse in cache.
func PutNSFW(key string, data structs.NSFWResponse) error {

	bytes, err := json.Marshal(&data)
	if err != nil {
		log.Println("PutNSFW: unable to marshal data ", err)
		return err
	}

	return nClient.Set(key, string(bytes), 0).Err()

}

// GetNSFW gets a NSFWResponse from the cache.
func GetNSFW(key string) (data structs.NSFWResponse, err error) {

	if nClient.Exists(key).Val() == 0 {
		return data, fmt.Errorf("GetNSFW: Key %s not found", key)
	}

	result, err := nClient.Get(key).Result()
	if err != nil {
		log.Println(fmt.Sprintf("GetNSFW: Unable to perform GET for key %s: %s", key, err))
		return
	}

	err = json.Unmarshal([]byte(result), &data)
	return

}
