package brain

import (
	"github.com/Efabien/string"
	"github.com/Efabien/cognitive_types"
	"strings"
)

func GenerateVault(keywords cognitivetypes.Keywords ,rawIntents cognitivetypes.Raw, scope int, degree int) *cognitivetypes.Vault {
	intents := tool.Precompute(rawIntents)
	weigths := getWeigths(intents)
	return &cognitivetypes.Vault {
		Intents: intents,
		Keywords: keywords,
		Scope: scope,
		Degree: degree,
		Weigths: weigths,
	}
}



func depth(what string, keywords cognitivetypes.Keywords) (result int) {
	keyword := keywords[what]
	for _, value := range keyword {
		currentLength := len(value)
		if currentLength > result {
			result = currentLength
		}
	}
	return
}

func ExtractAll(txt string, keywords cognitivetypes.Keywords, degree int)(result[]map[string][]string) {
	for key, _ := range keywords {
		result = append(result, Extract(txt, key, keywords, degree))
	}
	return
}

func Extract(txt string, what string, keywords cognitivetypes.Keywords, degree int) map[string][]string {
	result := make(map[string][]string)
	data := strings.Fields(txt)
	for key, value := range keywords[what] {
		for _, element := range value {
			for i := 1; i <= depth(what, keywords); i ++ {
				tool.PortionReading(data, i, func(array []string, from int, to int) {
					if tool.ExactMatch(element, array, degree) {
						result[what] = append(result[what], key)
					}
				})
			}
		}
	}
	return result
}

func Detect(input string, intents cognitivetypes.Intents, scope int, degree int, weigths []cognitivetypes.IntentWeigths) (res []cognitivetypes.Detection) {
	data := strings.Fields(input)
	for intent, value := range intents {
		texts := value.Texts
		var result cognitivetypes.Detection
		result.Intent = intent
		for _, element := range texts {
			for i := 1; i <= scope; i++ {
				tool.PortionReading(data, i, func(array[]string, from int, to int) {
					tool.PortionReading(element, i, func(proc[]string, start int, end int) {
						if tool.ExactMatch(array, proc, degree) {
							result.Matchs = append(result.Matchs, array...)
							var weigthsSums float32
							for _, word := range array {
								weigthsSums += GetWordWeigth(intent, word, weigths)
							}
							result.Score += (1 / float32(len(array)) * float32(len(texts))) * weigthsSums
						}
					})
				})
			}
		}
		res = append(res, result)
	}
	return
}

func GetWordWeigth(intent string, word string, weigths []cognitivetypes.IntentWeigths)(res float32) {
	var data cognitivetypes.IntentWeigths
	var result cognitivetypes.WordWeigth
	for _, w := range weigths {
		if w.Intent == intent {
			data = w
			break
		}
	}
	for _, item := range data.Weigths {
		if item.Word == word {
			result = item
			break
		}
	}
	return result.Weigth
}

func getWeigths(intents cognitivetypes.Intents)(result []cognitivetypes.IntentWeigths) {
	for key, intent := range intents {
		res := cognitivetypes.IntentWeigths{}
		res.Intent = key
		weigths := []cognitivetypes.WordWeigth{}
		for _, sentence := range intent.Texts {
			toInsert := cognitivetypes.WordWeigth{}
			for _, word := range sentence {
				toInsert.Word = word
				toInsert.Weigth = calculateW(word, intent.Texts, intents)
				weigths = append(weigths, toInsert)
			}
		}
		res.Weigths = weigths
		result = append(result, res)
	}
	return
}

func calculateW(word string, texts[][]string, intents cognitivetypes.Intents) (w float32) {
	var wordFrequency float32
	var docFrenquency float32
	for _, sentence := range texts {
		filtered := tool.Filter(sentence, func(item string, index int)bool {
			return item == word
		})
		wordFrequency += float32(len(filtered))
	}
	for _, intent := range intents {
		if isWordInIntent(word, intent) {
			docFrenquency ++
		}
	}
	return wordFrequency * (1 / docFrenquency)
}

func isWordInIntent(word string, intent cognitivetypes.Intent) (result bool) {
	for _, sentence := range intent.Texts {
		if tool.Some(sentence, func(item string, index int)bool {
			return item == word
		}) {
			result = true
			break
		}
	}
	return
}