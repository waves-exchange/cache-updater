package entities

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/ventuary-lab/cache-updater/swagger-types/models"
	"os"
	"regexp"
	"encoding/json"
)

func unwrapDefaultRegex (rawregex *string, defaultRegex string) string {
	if rawregex == nil {
		return defaultRegex
	} else {
		return *rawregex
    }
}

func GetDBCredentials () (string, string, string, string, string) {
	envLoadErr := godotenv.Load(".env")
	if envLoadErr != nil {
		_ = godotenv.Load(".env.example")
	}

	dbhost := "localhost"
	dbport := "5432"
	if os.Getenv("DB_HOST") != "" {
		dbhost = os.Getenv("DB_HOST")
	}
	if os.Getenv("DB_PORT") != "" {
		dbport = os.Getenv("DB_PORT")
	}
	dbuser := os.Getenv("DB_USERNAME")
	dbpass := os.Getenv("DB_PASS")
	dbdatabase := os.Getenv("DB_NAME")

	return dbhost, dbport, dbuser, dbpass, dbdatabase
}

func MapNodeDataToDict (addressData []byte) map[string]string {
	var records []map[string]interface{}
	json.Unmarshal([]byte(addressData), &records)

	nodeData := map[string]string{}

	for i := 0; i < len(records); i++ {
		record := records[i]

		key := record["key"].(string)
		valueType := record["type"].(string)
		rawValue := record["value"]

		if valueType == "integer" {
			nodeData[key] = fmt.Sprintf("%v", int(rawValue.(float64)))
		} else if valueType == "string" {
			nodeData[key] = rawValue.(string)
		}
	}

	return nodeData
}

func UpdateAll (nodeData *map[string]string, GetKeys func(*string)[]string) ([]string, map[string]map[string]string, error) {
	var ids []string
	regexKeys := GetKeys(nil)
	heightKey := regexKeys[0]
	heightRegex, heightRegexErr := regexp.Compile(heightKey)
	var nodeKeys []string
	resolveData := make(map[string]map[string]string)

	for k, _ := range *nodeData {
		for _, regexKey := range regexKeys {
			compiledRegex := regexp.MustCompile(regexKey)

			if len(compiledRegex.FindSubmatch([]byte(k))) == 0 {
				continue
			}
		}
		nodeKeys = append(nodeKeys, k)
	}

	for _, k := range nodeKeys {
		heightRegexSubmatches := heightRegex.FindSubmatch([]byte(k))

		if len(heightRegexSubmatches) < 2 {
			continue
		}

		matchedAddress := string(heightRegexSubmatches[1])

		if matchedAddress != "" {
			ids = append(ids, matchedAddress)
			resolveData[matchedAddress] = map[string]string{}
			validKeys := GetKeys(&matchedAddress)

			for _, validKey := range validKeys {
				for _, k := range nodeKeys {
					if k == validKey {
						resolveData[matchedAddress][k] = (*nodeData)[k]
					}
				}
			}
		}
	}

	return ids, resolveData, heightRegexErr
}

func MapStateChangesDataToDict (stateChanges *models.StateChangesStateChanges) map[string]string {
	res := make(map[string]string)
	//res := make(map[string]interface{})

	for _, change := range stateChanges.Data {
		key := *change.Key

		// Enum: [integer boolean binary string]
		valueType := *change.Type

		if valueType == "integer" {
			res[key] = fmt.Sprintf("%v", change.Value.(float64))
		}
		if valueType == "string" {
			res[key] = change.Value.(string)
		}
		if valueType == "boolean" {
			if change.Value.(bool) {
				res[key] = "true"
			} else {
				res[key] = "false"
			}
		}
	}

	return res
}

func DefaultErrorHandler (err error) {
	if err != nil {
		fmt.Printf("Error occured. %v \n", err)
	}
}

