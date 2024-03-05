package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/pkg/browser"
)

var json_data = `
{
  "squadName": "Super hero squad",
  "homeTown": "Metro City",
  "formed": 2016,
  "secretBase": "Super tower",
  "active": true,
  "powers": [
          "Million tonne punch",
          "Damage resistance",
          "Superhuman reflexes"
        ],
  "members": [
    {
      "name": "Molecule Man",
      "age": 29,
      "secretIdentity": "Dan Jukes",
      "powers": [
        "Radiation resistance",
        "Turning tiny",
        "Radiation blast"
      ]
    },
	"lol",
    {
      "name": "Madame Uppercut",
      "age": 39,
      "secretIdentity": "Jane Wilson"
    },
    {
      "name": "Eternal Flame",
      "age": 1000000,
      "secretIdentity": "Unknown",
      "powers": [
        "Immortality",
        "Heat Immunity",
        "Inferno",
        "Teleportation",
        "Interdimensional travel"
      ]
    }
  ]
}
`

type j_data struct {
	_key      string
	_map      map[string]interface{}
	arraysMap map[string][]interface{}
	jDataMap  map[string][]j_data
	level     int
	index     int
}

func recursion_map(data map[string]interface{}, key string, level, index int) j_data {
	j := j_data{
		_map:      make(map[string]interface{}),
		arraysMap: make(map[string][]interface{}),
		jDataMap:  make(map[string][]j_data),
		_key:      key,
		level:     level,
		index:     index,
	}

	for k, v := range data {
		switch v.(type) {
		case []interface{}:
			var i int
			for _, val := range v.([]interface{}) {
				if m, ok := val.(map[string]interface{}); ok {
					o := recursion_map(m, k, level+1, i)
					j.jDataMap[k] = append(j.jDataMap[k], o)
				} else {
					j.arraysMap[k] = append(j.arraysMap[k], val)
				}

				i++
			}
		case map[string]interface{}:
			o := recursion_map(v.(map[string]interface{}), k, level+1, 0)
			j.jDataMap[k] = append(j.jDataMap[k], o)
		default:
			j._map[k] = v
		}
	}

	return j
}

func recursion_print(j j_data) {
	for i := 0; i < j.level; i++ {
		fmt.Print("  ")
	}

	if j.level != 0 {
		fmt.Printf("%s:: index %d: level %d:\n", j._key, j.index, j.level)
	} else {
		fmt.Printf("root:: index %d: level %d:\n", j.index, j.level)
	}

	j.level += 1

	for k, v := range j._map {
		for i := 0; i < j.level; i++ {
			fmt.Print("  ")
		}

		fmt.Printf("%s: %v\n", k, v)
	}

	for k, v := range j.arraysMap {
		for i := 0; i < j.level; i++ {
			fmt.Print("  ")
		}

		fmt.Printf("%s: %v\n", k, v)
	}

	for _, v := range j.jDataMap {
		for _, val := range v {
			recursion_print(val)
		}
	}
}

type Mermaid struct {
	MMD string
}

func main() {
	var filePath, port string
	flag.StringVar(&filePath, "file", "", "file path")
	flag.StringVar(&port, "port", "8888", "port")
	flag.Parse()

	var _json []byte
	var err error
	if filePath != "" {
		_json, err = os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading file: %s", err)
			return
		}
	} else {
		_json = []byte(json_data)
	}

	var l1 map[string]interface{}
	err = json.Unmarshal(_json, &l1)
	if err != nil {
		fmt.Printf("Error L1: %s", err)
		return
	}

	j := recursion_map(l1, "", 0, 0)
	recursion_print(j)

	tmp := template.Must(template.ParseFiles("index.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmp.Execute(w, Mermaid{MMD: "graph TD\nA[Christmas] -->|Get money| B(Go shopping)\nB --> C{Let me think}\nC -->|One| D[Laptop]\nC -->|Two| E[iPhone]\nC -->|Three| F[fa:fa-car Car]\n"})
	})
	browser.OpenURL("http://localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}
