package entities;

import (
	"os"
	"regexp"

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

func CollectionUpdateAll(
	nodeData *map[string]string,
	GetKeys func(*string) []string,
	MapItemToModel func(string, map[string]string) interface{},
) interface{} {
	var ids []string
	var result []interface{}
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

	if heightRegexErr != nil {
		return result
	}

	for _, id := range ids {
		mappedModel := MapItemToModel(id, resolveData[id])
		result = append(result, mappedModel)
	}

	return result
}
