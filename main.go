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
	txt "text/template"

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

var initialNodes = `
[
  {
    "id": "uokKcnPMRv9nKWd-F-_YY",
    "data": {
      "label": {
        "squadName": "Super hero squad",
        "homeTown": "Metro City",
        "formed": 2016,
        "secretBase": "Super tower",
        "active": true
      }
    },
    "position": {
      "x": 462,
      "y": 12
    },
    "type": "jsonVis",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 511,
    "x": 462,
    "y": 12
  },
  {
    "id": "aXp1nbws2IiRJsjHO5E1C",
    "data": {
      "label": "members"
    },
    "type": "jsonVis",
    "position": {
      "x": 462,
      "y": 262
    },
    "parent": "uokKcnPMRv9nKWd-F-_YY",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 513,
    "x": 462,
    "y": 262
  },
  {
    "id": "L3fzrf0cmkcLFzeCdk5MG",
    "data": {
      "label": {
        "name": "Molecule Man",
        "age": 29,
        "secretIdentity": "Dan Jukes"
      }
    },
    "type": "jsonVis",
    "position": {
      "x": 1943.25,
      "y": 512
    },
    "parent": "aXp1nbws2IiRJsjHO5E1C",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 515,
    "x": 1943.25,
    "y": 512
  },
  {
    "id": "MF7EROOWXVAX8E3CqlBda",
    "data": {
      "label": "powers"
    },
    "type": "jsonVis",
    "position": {
      "x": 1943.25,
      "y": 762
    },
    "parent": "L3fzrf0cmkcLFzeCdk5MG",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 517,
    "x": 1943.25,
    "y": 762
  },
  {
    "id": "cyWW5lEqFvWUpzyWC8iYQ",
    "data": {
      "label": "Radiation resistance"
    },
    "type": "jsonVis",
    "position": {
      "x": 1812,
      "y": 1012
    },
    "parent": "MF7EROOWXVAX8E3CqlBda",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 519,
    "x": 1812,
    "y": 1012
  },
  {
    "id": "NR5o-GQJoBLDV577J1D7L",
    "data": {
      "label": "Turning tiny"
    },
    "type": "jsonVis",
    "position": {
      "x": 2112,
      "y": 1012
    },
    "parent": "MF7EROOWXVAX8E3CqlBda",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 521,
    "x": 2112,
    "y": 1012
  },
  {
    "id": "Ev8wFKpv2Lc-I1-FxoxUD",
    "data": {
      "label": "Radiation blast"
    },
    "type": "jsonVis",
    "position": {
      "x": 1512,
      "y": 1012
    },
    "parent": "MF7EROOWXVAX8E3CqlBda",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 523,
    "x": 1512,
    "y": 1012
  },
  {
    "id": "9kbMx45W4QKeW1DynV-AN",
    "data": {
      "label": "lol"
    },
    "type": "jsonVis",
    "position": {
      "x": 12,
      "y": 512
    },
    "parent": "aXp1nbws2IiRJsjHO5E1C",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 525,
    "x": 12,
    "y": 512
  },
  {
    "id": "Z2K5-GXcYDy9MewoZRf27",
    "data": {
      "label": {
        "name": "Madame Uppercut",
        "age": 39,
        "secretIdentity": "Jane Wilson"
      }
    },
    "type": "jsonVis",
    "position": {
      "x": 312,
      "y": 512
    },
    "parent": "aXp1nbws2IiRJsjHO5E1C",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 527,
    "x": 312,
    "y": 512
  },
  {
    "id": "iDXyn5OcnG17jZ1WmDVsq",
    "data": {
      "label": {
        "name": "Eternal Flame",
        "age": 1000000,
        "secretIdentity": "Unknown"
      }
    },
    "type": "jsonVis",
    "position": {
      "x": 612,
      "y": 512
    },
    "parent": "aXp1nbws2IiRJsjHO5E1C",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 529,
    "x": 612,
    "y": 512
  },
  {
    "id": "0zyljObQeS9W9vpBrdyeE",
    "data": {
      "label": "powers"
    },
    "type": "jsonVis",
    "position": {
      "x": 612,
      "y": 762
    },
    "parent": "iDXyn5OcnG17jZ1WmDVsq",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 531,
    "x": 612,
    "y": 762
  },
  {
    "id": "O6tpqkcq2EzqY89khlusu",
    "data": {
      "label": "Immortality"
    },
    "type": "jsonVis",
    "position": {
      "x": 312,
      "y": 1012
    },
    "parent": "0zyljObQeS9W9vpBrdyeE",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 533,
    "x": 312,
    "y": 1012
  },
  {
    "id": "AJ_GuX3R5VNkzdktk52q5",
    "data": {
      "label": "Heat Immunity"
    },
    "type": "jsonVis",
    "position": {
      "x": 912,
      "y": 1012
    },
    "parent": "0zyljObQeS9W9vpBrdyeE",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 535,
    "x": 912,
    "y": 1012
  },
  {
    "id": "7uofEQ_p5J5EtvLeqAkej",
    "data": {
      "label": "Inferno"
    },
    "type": "jsonVis",
    "position": {
      "x": 1212,
      "y": 1012
    },
    "parent": "0zyljObQeS9W9vpBrdyeE",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 537,
    "x": 1212,
    "y": 1012
  },
  {
    "id": "5ygEORyDxPRH529qqxLcL",
    "data": {
      "label": "Teleportation"
    },
    "type": "jsonVis",
    "position": {
      "x": 612,
      "y": 1012
    },
    "parent": "0zyljObQeS9W9vpBrdyeE",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 539,
    "x": 612,
    "y": 1012
  },
  {
    "id": "fIv7iBMg4nBQOxCoZGz1y",
    "data": {
      "label": "Interdimensional travel"
    },
    "type": "jsonVis",
    "position": {
      "x": 12,
      "y": 1012
    },
    "parent": "0zyljObQeS9W9vpBrdyeE",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 541,
    "x": 12,
    "y": 1012
  }
]
`

var initialEdges = `
[
  {
    "id": "uokKcnPMRv9nKWd-F-_YY",
    "data": {
      "label": {
        "squadName": "Super hero squad",
        "homeTown": "Metro City",
        "formed": 2016,
        "secretBase": "Super tower",
        "active": true
      }
    },
    "position": {
      "x": 462,
      "y": 12
    },
    "type": "jsonVis",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 511,
    "x": 462,
    "y": 12
  },
  {
    "id": "aXp1nbws2IiRJsjHO5E1C",
    "data": {
      "label": "members"
    },
    "type": "jsonVis",
    "position": {
      "x": 462,
      "y": 262
    },
    "parent": "uokKcnPMRv9nKWd-F-_YY",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 513,
    "x": 462,
    "y": 262
  },
  {
    "id": "L3fzrf0cmkcLFzeCdk5MG",
    "data": {
      "label": {
        "name": "Molecule Man",
        "age": 29,
        "secretIdentity": "Dan Jukes"
      }
    },
    "type": "jsonVis",
    "position": {
      "x": 1943.25,
      "y": 512
    },
    "parent": "aXp1nbws2IiRJsjHO5E1C",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 515,
    "x": 1943.25,
    "y": 512
  },
  {
    "id": "MF7EROOWXVAX8E3CqlBda",
    "data": {
      "label": "powers"
    },
    "type": "jsonVis",
    "position": {
      "x": 1943.25,
      "y": 762
    },
    "parent": "L3fzrf0cmkcLFzeCdk5MG",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 517,
    "x": 1943.25,
    "y": 762
  },
  {
    "id": "cyWW5lEqFvWUpzyWC8iYQ",
    "data": {
      "label": "Radiation resistance"
    },
    "type": "jsonVis",
    "position": {
      "x": 1812,
      "y": 1012
    },
    "parent": "MF7EROOWXVAX8E3CqlBda",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 519,
    "x": 1812,
    "y": 1012
  },
  {
    "id": "NR5o-GQJoBLDV577J1D7L",
    "data": {
      "label": "Turning tiny"
    },
    "type": "jsonVis",
    "position": {
      "x": 2112,
      "y": 1012
    },
    "parent": "MF7EROOWXVAX8E3CqlBda",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 521,
    "x": 2112,
    "y": 1012
  },
  {
    "id": "Ev8wFKpv2Lc-I1-FxoxUD",
    "data": {
      "label": "Radiation blast"
    },
    "type": "jsonVis",
    "position": {
      "x": 1512,
      "y": 1012
    },
    "parent": "MF7EROOWXVAX8E3CqlBda",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 523,
    "x": 1512,
    "y": 1012
  },
  {
    "id": "9kbMx45W4QKeW1DynV-AN",
    "data": {
      "label": "lol"
    },
    "type": "jsonVis",
    "position": {
      "x": 12,
      "y": 512
    },
    "parent": "aXp1nbws2IiRJsjHO5E1C",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 525,
    "x": 12,
    "y": 512
  },
  {
    "id": "Z2K5-GXcYDy9MewoZRf27",
    "data": {
      "label": {
        "name": "Madame Uppercut",
        "age": 39,
        "secretIdentity": "Jane Wilson"
      }
    },
    "type": "jsonVis",
    "position": {
      "x": 312,
      "y": 512
    },
    "parent": "aXp1nbws2IiRJsjHO5E1C",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 527,
    "x": 312,
    "y": 512
  },
  {
    "id": "iDXyn5OcnG17jZ1WmDVsq",
    "data": {
      "label": {
        "name": "Eternal Flame",
        "age": 1000000,
        "secretIdentity": "Unknown"
      }
    },
    "type": "jsonVis",
    "position": {
      "x": 612,
      "y": 512
    },
    "parent": "aXp1nbws2IiRJsjHO5E1C",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 529,
    "x": 612,
    "y": 512
  },
  {
    "id": "0zyljObQeS9W9vpBrdyeE",
    "data": {
      "label": "powers"
    },
    "type": "jsonVis",
    "position": {
      "x": 612,
      "y": 762
    },
    "parent": "iDXyn5OcnG17jZ1WmDVsq",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 531,
    "x": 612,
    "y": 762
  },
  {
    "id": "O6tpqkcq2EzqY89khlusu",
    "data": {
      "label": "Immortality"
    },
    "type": "jsonVis",
    "position": {
      "x": 312,
      "y": 1012
    },
    "parent": "0zyljObQeS9W9vpBrdyeE",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 533,
    "x": 312,
    "y": 1012
  },
  {
    "id": "AJ_GuX3R5VNkzdktk52q5",
    "data": {
      "label": "Heat Immunity"
    },
    "type": "jsonVis",
    "position": {
      "x": 912,
      "y": 1012
    },
    "parent": "0zyljObQeS9W9vpBrdyeE",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 535,
    "x": 912,
    "y": 1012
  },
  {
    "id": "7uofEQ_p5J5EtvLeqAkej",
    "data": {
      "label": "Inferno"
    },
    "type": "jsonVis",
    "position": {
      "x": 1212,
      "y": 1012
    },
    "parent": "0zyljObQeS9W9vpBrdyeE",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 537,
    "x": 1212,
    "y": 1012
  },
  {
    "id": "5ygEORyDxPRH529qqxLcL",
    "data": {
      "label": "Teleportation"
    },
    "type": "jsonVis",
    "position": {
      "x": 612,
      "y": 1012
    },
    "parent": "0zyljObQeS9W9vpBrdyeE",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 539,
    "x": 612,
    "y": 1012
  },
  {
    "id": "fIv7iBMg4nBQOxCoZGz1y",
    "data": {
      "label": "Interdimensional travel"
    },
    "type": "jsonVis",
    "position": {
      "x": 12,
      "y": 1012
    },
    "parent": "0zyljObQeS9W9vpBrdyeE",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 541,
    "x": 12,
    "y": 1012
  }
]
`

type r_data struct {
	InitialNodes string
	InitialEdges string
}

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
	tmp_r := txt.Must(txt.ParseFiles("my-react-app/src/App.tmp.jsx"))
	http.HandleFunc("/react", func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Create("my-react-app/src/App.jsx")
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(fmt.Sprintf("{\"error\": \"%s\"}", err)))
			return
		}
		tmp_r.Execute(f, r_data{InitialNodes: initialNodes, InitialEdges: initialEdges})
		f.Close()

		result := api.Build(api.BuildOptions{
			EntryPoints: []string{"my-react-app/src/main.jsx"},
			Bundle:      true,
			Outfile:     "my-react-app/dist/bundle.js",
			Write:       true,
			JSX:         api.JSXAutomatic,
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
