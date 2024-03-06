package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
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

func removeSpecialChars(s string) string {
	chars := []string{"(", ")", "@"}
	r := strings.Join(chars, "")
	re := regexp.MustCompile("[" + r + "]+")
	return re.ReplaceAllString(s, "")
}

func recursion_mermaid(j j_data, perviousKey string) string {
	var mmd string
	m_key := fmt.Sprintf("%d%d%s", j.level, j.index, j._key)
	if j.level == 0 {
		mmd += "graph LR\n"
	} else if perviousKey != "" {
		mmd += fmt.Sprintf("%s --> %s\n", perviousKey, m_key)
	}

	mmd += fmt.Sprintf("%s[\n<p style='text-align: left;'>", m_key)

	for k, v := range j._map {
		mmd += fmt.Sprintf("%s: %v\n", k, removeSpecialChars(fmt.Sprintf("%v", v)))
	}

	mmd += "</>]\n"

	for k, v := range j.arraysMap {
		mmd += fmt.Sprintf("%s --> %s%d%d[%s]\n", m_key, k, j.level, j.index, k)
		for n, val := range v {
			mmd += fmt.Sprintf("%s%d%d --> %s%d%d%d[%v]\n", k, j.level, j.index, k, j.level, j.index, n, removeSpecialChars(fmt.Sprintf("%v", val)))
		}
	}

	for _, v := range j.jDataMap {
		for i, val := range v {
			mm_key := fmt.Sprintf("%d%d%s", j.level, j.index, val._key)

			if i == 0 {
				mmd += fmt.Sprintf("%s --> %s[%s]\n", m_key, mm_key, val._key)
			}

			mmd += recursion_mermaid(val, mm_key)
		}
	}

	return mmd
}

type Mermaid struct {
	MMD string
}

func main() {
	var port string
	flag.StringVar(&port, "port", "8888", "port")
	flag.Parse()

	tmp := template.Must(template.ParseFiles("index.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		file := r.URL.Query().Get("file")
		l1, err := newFunction(file)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(fmt.Sprintf("{\"error\": \"%s\"}", err)))
			return
		}

		j := recursion_map(l1, "", 0, 0)
		m := recursion_mermaid(j, "")
		recursion_print(j)
		fmt.Print(m)
		tmp.Execute(w, Mermaid{MMD: m})
	})
	http.HandleFunc("/react", func(w http.ResponseWriter, r *http.Request) {
		result := api.Build(api.BuildOptions{
			EntryPoints: []string{"my-react-app/src/main.jsx"},
			Bundle:      true,
			Outfile:     "my-react-app/dist/bundle.js",
			Write:       true,
		})
		if len(result.Errors) > 0 {
			fmt.Println(result.Errors)
			return
		}
		http.ServeFile(w, r, "react.html")
	})
	http.HandleFunc("/dist/bundle.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "my-react-app/dist/bundle.js")
	})
	http.HandleFunc("/dist/bundle.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "my-react-app/dist/bundle.css")
	})
	browser.OpenURL("http://localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}

func newFunction(filePath string) (map[string]interface{}, error) {
	var _json []byte
	var err error
	if filePath != "" {
		_json, err = os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("error reading file: %s", err)
		}
	} else {
		_json = []byte(json_data)
	}

	var l1 map[string]interface{}
	err = json.Unmarshal(_json, &l1)
	if err != nil {
		return nil, fmt.Errorf("error L1: %s", err)
	}

	return l1, nil
}
