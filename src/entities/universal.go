package entities;

import (
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"net/http"
	"os"
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

//func CollectionUpdateAll(
//	nodeData *map[string]string,
//	GetKeys func(*string) []string,
//	MapItemToModel func(string, map[string]string) interface{},
//) []interface{} {
//	var ids []string
//	regexKeys := GetKeys(nil)
//	heightKey := regexKeys[0]
//	heightRegex, heightRegexErr := regexp.Compile(heightKey)
//	var nodeKeys []string
//	resolveData := make(map[string]map[string]string)
//
//	for k, _ := range *nodeData {
//		for _, regexKey := range regexKeys {
//			compiledRegex := regexp.MustCompile(regexKey)
//
//			if len(compiledRegex.FindSubmatch([]byte(k))) == 0 {
//				continue
//			}
//		}
//		nodeKeys = append(nodeKeys, k)
//	}
//
//	for _, k := range nodeKeys {
//		heightRegexSubmatches := heightRegex.FindSubmatch([]byte(k))
//
//		if len(heightRegexSubmatches) < 2 {
//			continue
//		}
//
//		matchedAddress := string(heightRegexSubmatches[1])
//
//		if matchedAddress != "" {
//			ids = append(ids, matchedAddress)
//			resolveData[matchedAddress] = map[string]string{}
//			validKeys := GetKeys(&matchedAddress)
//
//			for _, validKey := range validKeys {
//				for _, k := range nodeKeys {
//					if k == validKey {
//						resolveData[matchedAddress][k] = (*nodeData)[k]
//					}
//				}
//			}
//		}
//	}
//
//	if heightRegexErr != nil {
//		return make([]interface{}, 0)
//	}
//
//	result := make([]interface{}, len(ids))
//	for i, id := range ids {
//		mappedModel := MapItemToModel(id, resolveData[id])
//		result[i] = mappedModel
//	}
//
//	return result
//}

func FetchLastBlock () []byte {
	nodeUrl := os.Getenv("NODE_URL")
	connectionUrl := fmt.Sprintf("%v/blocks/headers/last", nodeUrl)
	response, err := http.Get(connectionUrl)

	if err != nil {
		fmt.Printf("Error occured on fetch... %v \n", err)
		return make([]byte, 0)
	}

	defer response.Body.Close()

	byteValue, readErr := ioutil.ReadAll(response.Body)

	if readErr != nil {
		return make([]byte, 0)
	}

	return byteValue
}

func FetchBlocksRange (heightMin, heightMax string) []byte {
	nodeUrl := os.Getenv("NODE_URL")
	connectionUrl := fmt.Sprintf("%v/blocks/headers/seq/%v/%v", nodeUrl, heightMin, heightMax)
	response, err := http.Get(connectionUrl)

	if err != nil {
		fmt.Printf("Error occured on fetch... %v \n", err)
		return make([]byte, 0)
	}

	defer response.Body.Close()

	byteValue, readErr := ioutil.ReadAll(response.Body)

	if readErr != nil {
		return make([]byte, 0)
	}

	return byteValue
}