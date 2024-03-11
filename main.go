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
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/nulab/autog"
	"github.com/nulab/autog/graph"
	"github.com/pkg/browser"
	"github.com/sithumonline/json_x/data"
	"golang.org/x/exp/slices"
)

type r_data struct {
	InitialNodes string `json:"initialNodes,omitempty"`
	InitialEdges string `json:"initialEdges,omitempty"`
}

type j_data struct {
	_key      string
	_map      map[string]interface{}
	arraysMap map[string][]interface{}
	jDataMap  map[string][]j_data
	level     int
	index     int
	id        string
}

func recursion_map(data map[string]interface{}, key string, level, index int) j_data {
	id, _ := gonanoid.Generate("abcde", 21)
	j := j_data{
		_map:      make(map[string]interface{}),
		arraysMap: make(map[string][]interface{}),
		jDataMap:  make(map[string][]j_data),
		_key:      key,
		level:     level,
		index:     index,
		id:        id,
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

type flow struct {
	ID   string `json:"id,omitempty"`
	Data struct {
		Label map[string]interface{} `json:"label,omitempty"`
	} `json:"data,omitempty"`
	Type     string `json:"type,omitempty"`
	Position struct {
		X float64 `json:"x,omitempty"`
		Y float64 `json:"y,omitempty"`
	} `json:"position,omitempty"`
	Parent         string `json:"parent,omitempty"`
	TargetPosition string `json:"targetPosition,omitempty"`
	SourcePosition string `json:"sourcePosition,omitempty"`
	Width          int    `json:"width,omitempty"`
	Height         int    `json:"height,omitempty"`
	H              int    `json:"$H,omitempty"`
	X              int    `json:"x,omitempty"`
	Y              int    `json:"y,omitempty"`
}

func recursion_flow(j j_data, previousKey string) ([]flow, [][]string) {
	var fl []flow
	var adj [][]string

	f := flow{
		ID:     j.id,
		Type:   "jsonVis",
		Parent: previousKey,
	}
	f.Data.Label = make(map[string]interface{})

	for k, v := range j._map {
		f.Data.Label[k] = v
	}

	for k, v := range j.arraysMap {
		id, _ := gonanoid.Generate("abcde", 21)
		flow1 := flow{
			ID:     id,
			Type:   "jsonVis",
			Parent: j.id,
		}
		flow1.Data.Label = make(map[string]interface{})
		flow1.Data.Label[k] = k
		fl = append(fl, flow1)
		adj = append(adj, []string{j.id, id})
		for _, val := range v {
			id2, _ := gonanoid.Generate("abcde", 21)
			flow2 := flow{
				ID:     id2,
				Type:   "jsonVis",
				Parent: id,
			}
			flow2.Data.Label = make(map[string]interface{})
			flow2.Data.Label[fmt.Sprintf("%v", val)] = val
			fl = append(fl, flow2)
			adj = append(adj, []string{id, id2})
		}
	}

	for _, v := range j.jDataMap {
		for _, val := range v {
			adj = append(adj, []string{j.id, val.id})
			gg, mm := recursion_flow(val, j.id)
			fl = append(fl, gg...)
			adj = append(adj, mm...)
		}
	}

	fl = append(fl, f)
	return fl, adj
}

type edges struct {
	Id     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
}

func autoGraph(f []flow, dd [][]string) []flow {
	// obtain a graph.Source (here by converting the input to EdgeSlice)
	src := graph.EdgeSlice(dd)

	// run the default autolayout pipeline
	layout := autog.Layout(src)

	for _, l := range layout.Nodes {
		i := slices.IndexFunc(f, func(d flow) bool {
			return d.ID == l.ID
		})
		if i == -1 {
			continue
		}

		f[i].Position.X = l.X
		f[i].Position.Y = l.Y
	}

	return f
}

func getEdges(dd [][]string) []edges {
	var e []edges
	for _, d := range dd {
		id, _ := gonanoid.Generate("abcde", 21)
		e = append(e, edges{Id: id, Source: d[0], Target: d[1]})
	}

	return e
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
		tmp_r.Execute(f, r_data{InitialNodes: data.Nodes, InitialEdges: data.Edges})
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
	http.HandleFunc("/flow-json", func(w http.ResponseWriter, r *http.Request) {
		file := r.URL.Query().Get("file")
		l1, err := newFunction(file)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(fmt.Sprintf("{\"error\": \"%s\"}", err)))
			return
		}

		j := recursion_map(l1, "", 0, 0)
		fl, adj := recursion_flow(j, "")
		fl2 := autoGraph(fl, adj)
		e := getEdges(adj)
		w.Header().Set("Content-Type", "application/json")
		type response struct {
			Nodes []flow  `json:"nodes,omitempty"`
			Edges []edges `json:"edges,omitempty"`
		}
		json.NewEncoder(w).Encode(response{Nodes: fl2, Edges: e})
	})
	http.HandleFunc("/flow-html", func(w http.ResponseWriter, r *http.Request) {
		file := r.URL.Query().Get("file")
		l1, err := newFunction(file)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(fmt.Sprintf("{\"error\": \"%s\"}", err)))
			return
		}

		j := recursion_map(l1, "", 0, 0)
		fl, adj := recursion_flow(j, "")
		fl2 := autoGraph(fl, adj)
		e := getEdges(adj)

		fl2d, err := json.Marshal(fl2)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(fmt.Sprintf("{\"error\": \"%s\"}", err)))
			return
		}

		ed, err := json.Marshal(e)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(fmt.Sprintf("{\"error\": \"%s\"}", err)))
			return
		}

		f, err := os.Create("my-react-app/src/App.jsx")
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(fmt.Sprintf("{\"error\": \"%s\"}", err)))
			return
		}
		tmp_r.Execute(f, r_data{InitialNodes: string(fl2d), InitialEdges: string(ed)})
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
		_json = []byte(data.Json_data)
	}

	var l1 map[string]interface{}
	err = json.Unmarshal(_json, &l1)
	if err != nil {
		return nil, fmt.Errorf("error L1: %s", err)
	}

	return l1, nil
}
