package entities;

import (
	"os"

	"github.com/joho/godotenv"
)
type DAappNumberRecord struct {
	Key, Type string
	Value *int
}

type DAppStringRecord struct {
	Key, Type string
	Value *string
}

func StrArrayContains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

func unwrapDefaultRegex (rawregex *string, defaultRegex string) string {
	if rawregex == nil {
		return defaultRegex
	} else {
		return *rawregex
    }
}

func GetDBCredentials () (string, string, string) {
	envLoadErr := godotenv.Load(".env")
	if envLoadErr != nil {
		_ = godotenv.Load(".env.example")
	}

	dbuser := os.Getenv("DB_USERNAME");
	dbpass := os.Getenv("DB_PASS");
	dbdatabase := os.Getenv("DB_NAME");

	return dbuser, dbpass, dbdatabase
}

// func UpdateCollection (nodeData *map[string]string, constructor map[string]interface{}) {
// 	ids := []string{}
// 	result := []map[string]interface{}
// 	regexKeys := this.GetKeys(nil)
// 	heightKey := regexKeys[0]
// 	heightRegex, heightRegexErr := regexp.Compile(heightKey)
// 	nodeKeys := []string{}
// 	resolveData := make(map[string](map[string]string))

// 	for k, _ := range *nodeData {
// 		for _, regexKey := range regexKeys {
// 			compiledRegex := regexp.MustCompile(regexKey)

// 			if len(compiledRegex.FindSubmatch([]byte(k))) == 0 {
// 				continue;
// 			}
// 		}
// 		nodeKeys = append(nodeKeys, k)
// 	}

// 	for _, k := range nodeKeys {
// 		heightRegexSubmatches := heightRegex.FindSubmatch([]byte(k))

// 		if len(heightRegexSubmatches) < 2 {
// 			continue
// 		}

// 		matchedAddress := string(heightRegexSubmatches[1])

// 		if matchedAddress != "" {
// 			ids = append(ids, matchedAddress)
// 			resolveData[matchedAddress] = map[string]string{}
// 			validKeys := this.GetKeys(&matchedAddress)

// 			for _, validKey := range validKeys {
// 				for _, k := range nodeKeys {
// 					if k == validKey {
// 						resolveData[matchedAddress][k] = (*nodeData)[k]
// 					}
// 				}
// 			}
// 		}
// 	}

// 	if heightRegexErr != nil {
// 		return result
// 	}

// 	raw := constructor{}

// 	for _, id := range ids {
// 		mappedModel := raw.MapItemToModel(id, resolveData[id])
// 		result = append(result, *mappedModel)
// 	}

// 	return result
// }