package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

// curl -X POST -H 'Authorization: Key 74ade0e0754147b89896566bc052db70' -H "Content-Type: application/json" -d ' { "inputs": [ { "data": { "image": { "url": "https://farm7.staticflickr.com/5769/21094803716_da3cea21b8_o.jpg" } } } ] }' https://api.clarifai.com/v2/models/aaa03c23b3724a16a56b629203edc62c/versions/aa7f35c01e0642fda5cf400f543e7c40/outputs

// PredictInfoStruct is returned by getPredictInfo()
type PredictInfoStruct struct {
	singleJpg     string
	predictionMap map[string]float64
}

func getPredictInfo(singleJpg string, url string, authKey string, chOut chan<- PredictInfoStruct) {
	debug := false
	fmt.Println("Starting getPredictInfo, JPG", singleJpg)

	predictionMap := make(map[string]float64)

	fmt.Println("URL:>", url)

	var jsonStr = `{ "inputs": [ { "data": { "image": { "url":` + `"` + singleJpg + `"` + `} } } ] }`
	var jsonBytes = []byte(jsonStr)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println(singleJpg, " response Status:", resp.Status)
	if debug {

		fmt.Println(strings.Compare(resp.Status, "200 OK"))
		fmt.Println("response Headers:", resp.Header)
	}
	body, _ := ioutil.ReadAll(resp.Body)

	if debug {
		fmt.Println("response Body:", string(body))
	}

	// See https://blog.golang.org/json-and-go

	var f interface{}
	err = json.Unmarshal(body, &f)

	if debug {
		fmt.Println("json.Unmarshal err:", err)
		fmt.Println("======== After json.Unmarshal()=====================")
		fmt.Println("f =", f)

		fmt.Println("=============================")
	}

	m := f.(map[string]interface{})

	for k, v := range m {
		switch vv := v.(type) {
		case string:
			if debug {
				fmt.Println(k, "is string", vv)
			}
		case float64:
			if debug {
				fmt.Println(k, "is float64", vv)
			}
		case []interface{}:
			if debug {
				fmt.Println(k, "is an array:")
			}
			for i, u := range vv {
				if debug {
					fmt.Println(i, u)
				}

				if i < 1 {
					uMap := u.(map[string]interface{})
					outputsData := uMap["data"]
					if debug {
						fmt.Println("outputsData = ", outputsData)
						fmt.Println("============Concepts =================")
					}
					outputsDataMap := outputsData.(map[string]interface{})
					outputsDataConcept := outputsDataMap["concepts"]

					if debug {
						fmt.Println("outputsDataConcept = ", outputsDataConcept)

						fmt.Println("============Concepts array  =================")
					}

					conceptsArray := outputsDataMap["concepts"].([]interface{})
					for i, message := range conceptsArray {
						if debug {
							fmt.Printf("OK: message %d => %s\n", i, message)
						}
						msgMap := message.(map[string]interface{})

						the_key := ""
						the_value := float64(0.0)
						var ok bool
						for k, v := range msgMap {
							if debug {
								if k == "name" || k == "value" {
									fmt.Println("key  = ", k, " value = ", v)
								}
							}

							if k == "name" {
								if the_key, ok = v.(string); ok {
									/* act on str */

								} else {
									/* not string */
									the_key = ""
								}

								if debug {
									fmt.Println("the_key  = ", the_key)
								}

							}

							if k == "value" {
								if the_value, ok = v.(float64); ok {
									/* act on str */

								} else {
									/* not string */
									the_value = 0.0
								}
								if debug {
									fmt.Println("the_value  = ", the_value)
								}
							}
						}
						if debug {
							fmt.Println("the_key  = ", the_key, "the_value  = ", the_value)
						}
						if the_key != "" {
							predictionMap[the_key] = the_value
						}
					}
				}
			}
		default:
			if debug {
				fmt.Println(k, "is of a type I don't know how to handle")
			}
		}
	}

	if debug {
		fmt.Println("predictionMap = ", predictionMap)
	}

	chOut <- PredictInfoStruct{singleJpg, predictionMap}
}

func readHTMLPage(url string) []string {

	fmt.Println("Starting readHTMLpage(), url =", url)

	listOfJpg := make([]string, 0, 10)

	resp, _ := http.Get(url)
	bytes, _ := ioutil.ReadAll(resp.Body)

	resultTemp := strings.Split(string(bytes), "\n")

	for i, u := range resultTemp {
		fmt.Println(i, u)

		if u != "" {
			listOfJpg = append(listOfJpg, u)
		}

		// if i >= 100 {
		//	break // For testing only
		// }

	}

	resp.Body.Close()
	return listOfJpg
}

func printSortAllTags(uniqPredictionMap map[string]int) {

	// fmt.Println("sorted by tag uniqPredictionMap results, len = ", len(uniqPredictionMap), "==============")
	tags := make([]string, 0, len(uniqPredictionMap))
	for tag := range uniqPredictionMap {
		tags = append(tags, tag)
	}
	sort.Strings(tags)

	fmt.Println("uniqPredictionMap results, len = ", len(tags))
	fmt.Println("tag, number of elements")
	for _, tag := range tags {
		fmt.Println(tag, uniqPredictionMap[tag])
	}
}

func printOnePriorityQueue(tag string, uniqPredictionJpgMap map[string]PriorityQueue) {

	tmpPq := uniqPredictionJpgMap[tag]
	fmt.Println("tag", tag, "length =", tmpPq.Len())
	tmpPq.Print()
	fmt.Println("Completed ==============")
}

func printOneURL(jpgURL string, fullMap map[string]map[string]float64) {

	predictionMap := fullMap[jpgURL]
	for key, value := range predictionMap {
		fmt.Println(key, "==>", value)
	}
}

func main() {
	debug := false
	fmt.Println("Starting main")

	imageURL := "https://s3.amazonaws.com/clarifai-data/backend/api-take-home/images.txt"
	predictionURL := "https://api.clarifai.com/v2/models/aaa03c23b3724a16a56b629203edc62c/versions/aa7f35c01e0642fda5cf400f543e7c40/outputs"
	authKey := "Key 74ade0e0754147b89896566bc052db70"
	qSize := 10

	fullMap := make(map[string]map[string]float64) // map['single.jpg]map[picture_tag][value]
	chOut := make(chan PredictInfoStruct)

	listOfJpg := readHTMLPage(imageURL)
	fmt.Println("======== Inside main() ===========")

	// Call go routine getPredictInfo to get predictions for every sngle jpg
	for i, singleJpg := range listOfJpg {
		fmt.Println("Producer loop", i, singleJpg)

		go getPredictInfo(singleJpg, predictionURL, authKey, chOut)

		if i%10 == 0 {
			fmt.Println("Producer loop sleeping i = ", i)
			time.Sleep(3000 * time.Millisecond)
		}

	}

	// example of predictionMap:
	// predictionMap["computer"] = 0.9918324
	// predictionMap["semiconductor"] = 0.98408026

	// Collect results of go routines in the fullMap using channel
	for i := range listOfJpg {
		fmt.Println("Consumer loop", i)

		predictInfo := <-chOut
		singleJpg := predictInfo.singleJpg
		predictionMap := predictInfo.predictionMap

		fullMap[singleJpg] = predictionMap
	}

	if debug {
		fmt.Println("in main() fullMap = ", fullMap)
	}

	uniqPredictionMap := make(map[string]int)
	uniqPredictionJpgMap := make(map[string]PriorityQueue)

	// Create a map where a tag points to a Priority Queue with size = 10

	for singleJpg, tmpPredictionMap := range fullMap {
		for tag, value := range tmpPredictionMap {
			tmpPq := make(PriorityQueue, 0)
			if counter, ok := uniqPredictionMap[tag]; !ok {
				uniqPredictionMap[tag] = 1
				tmpPq = make(PriorityQueue, 0)
			} else {
				uniqPredictionMap[tag] = counter + 1
				tmpPq = uniqPredictionJpgMap[tag]
			}

			item := &Item{
				value:    singleJpg,
				priority: value,
			}
			heap.Push(&tmpPq, item)

			// downsize to qSize
			for tmpPq.Len() > qSize {
				_ = heap.Pop(&tmpPq).(*Item)
			}

			uniqPredictionJpgMap[tag] = tmpPq
		}
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("======================================")
		fmt.Println("Enter <all_tags> to see the entire list of tags ")
		fmt.Println("Enter <tag 'tag_name'> to see top 10 for a given 'tag_name'")
		fmt.Println("Enter <jpg 'image_url'> to see all tags for a given 'image_url'")
		fmt.Println("Enter <quit> to quit")
		fmt.Println("======================================")

		input, _ := reader.ReadString('\n')
		// fmt.Println(input)

		if strings.TrimRight(input, "\n") == "all_tags" {
			printSortAllTags(uniqPredictionMap)
			continue
		} else if strings.HasPrefix(input, "tag ") {
			inputTmp := strings.TrimRight(input, "\n")
			tag := string(inputTmp[4:])
			printOnePriorityQueue(tag, uniqPredictionJpgMap)
			continue
		} else if strings.HasPrefix(input, "jpg ") {
			inputTmp := strings.TrimRight(input, "\n")
			jpgURL := string(inputTmp[4:])
			printOneURL(jpgURL, fullMap)
			continue
		} else if strings.TrimRight(input, "\n") == "quit" {
			fmt.Println("Exiting ")
			break
		} else {
			fmt.Println("Wrong entry, try again ")
		}
	}
	return
}
