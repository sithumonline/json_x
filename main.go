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
    "id": "JT9nF7hFMPVyjU0ktYg8u",
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
    "$H": 825,
    "x": 462,
    "y": 12
  },
  {
    "id": "13J4-8NZIVGrfNovVTq26",
    "data": {
      "label": "members"
    },
    "type": "jsonVis",
    "position": {
      "x": 462,
      "y": 262
    },
    "parent": "JT9nF7hFMPVyjU0ktYg8u",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 827,
    "x": 462,
    "y": 262
  },
  {
    "id": "UpZv9xrgQFuH9phiiu9mA",
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
    "parent": "13J4-8NZIVGrfNovVTq26",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 829,
    "x": 1943.25,
    "y": 512
  },
  {
    "id": "I3q-CPQSBdclDZAfORtoS",
    "data": {
      "label": "powers"
    },
    "type": "jsonVis",
    "position": {
      "x": 1943.25,
      "y": 762
    },
    "parent": "UpZv9xrgQFuH9phiiu9mA",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 831,
    "x": 1943.25,
    "y": 762
  },
  {
    "id": "y5gOWdbgjwxQ6ftXSyiLg",
    "data": {
      "label": "Radiation resistance"
    },
    "type": "jsonVis",
    "position": {
      "x": 1812,
      "y": 1012
    },
    "parent": "I3q-CPQSBdclDZAfORtoS",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 833,
    "x": 1812,
    "y": 1012
  },
  {
    "id": "Ce5SqLytWelETkNBueZ0n",
    "data": {
      "label": "Turning tiny"
    },
    "type": "jsonVis",
    "position": {
      "x": 2112,
      "y": 1012
    },
    "parent": "I3q-CPQSBdclDZAfORtoS",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 835,
    "x": 2112,
    "y": 1012
  },
  {
    "id": "h7fZmnnQZnyM8mqGY2m8F",
    "data": {
      "label": "Radiation blast"
    },
    "type": "jsonVis",
    "position": {
      "x": 1512,
      "y": 1012
    },
    "parent": "I3q-CPQSBdclDZAfORtoS",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 837,
    "x": 1512,
    "y": 1012
  },
  {
    "id": "J2XljHjG17fSTm7UsnKuk",
    "data": {
      "label": "lol"
    },
    "type": "jsonVis",
    "position": {
      "x": 12,
      "y": 512
    },
    "parent": "13J4-8NZIVGrfNovVTq26",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 839,
    "x": 12,
    "y": 512
  },
  {
    "id": "nF_bOmOXonb2DmOzqxOXo",
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
    "parent": "13J4-8NZIVGrfNovVTq26",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 841,
    "x": 312,
    "y": 512
  },
  {
    "id": "AfeEyOE2Qj0dNCYnQ_7oe",
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
    "parent": "13J4-8NZIVGrfNovVTq26",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 843,
    "x": 612,
    "y": 512
  },
  {
    "id": "Kc5TZRLv-CbMIC27TMF7y",
    "data": {
      "label": "powers"
    },
    "type": "jsonVis",
    "position": {
      "x": 612,
      "y": 762
    },
    "parent": "AfeEyOE2Qj0dNCYnQ_7oe",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 845,
    "x": 612,
    "y": 762
  },
  {
    "id": "Cl4g9tLNXqE2R8p5t1Dtn",
    "data": {
      "label": "Immortality"
    },
    "type": "jsonVis",
    "position": {
      "x": 312,
      "y": 1012
    },
    "parent": "Kc5TZRLv-CbMIC27TMF7y",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 847,
    "x": 312,
    "y": 1012
  },
  {
    "id": "2owVKiOQlID9xNu9s-4Z9",
    "data": {
      "label": "Heat Immunity"
    },
    "type": "jsonVis",
    "position": {
      "x": 912,
      "y": 1012
    },
    "parent": "Kc5TZRLv-CbMIC27TMF7y",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 849,
    "x": 912,
    "y": 1012
  },
  {
    "id": "yg1h4Vc0F0eP48PwCufJb",
    "data": {
      "label": "Inferno"
    },
    "type": "jsonVis",
    "position": {
      "x": 1212,
      "y": 1012
    },
    "parent": "Kc5TZRLv-CbMIC27TMF7y",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 851,
    "x": 1212,
    "y": 1012
  },
  {
    "id": "qcPxgcI1nTUzlI9m1Nwmi",
    "data": {
      "label": "Teleportation"
    },
    "type": "jsonVis",
    "position": {
      "x": 612,
      "y": 1012
    },
    "parent": "Kc5TZRLv-CbMIC27TMF7y",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 853,
    "x": 612,
    "y": 1012
  },
  {
    "id": "0ulNSgD2n8rIJ9TfcTKkE",
    "data": {
      "label": "Interdimensional travel"
    },
    "type": "jsonVis",
    "position": {
      "x": 12,
      "y": 1012
    },
    "parent": "Kc5TZRLv-CbMIC27TMF7y",
    "targetPosition": "top",
    "sourcePosition": "bottom",
    "width": 150,
    "height": 50,
    "$H": 855,
    "x": 12,
    "y": 1012
  }
]
`

var initialEdges = `
[
  {
    "id": "BCdZgXFYS1wIZTOhwKI9u",
    "source": "JT9nF7hFMPVyjU0ktYg8u",
    "target": "13J4-8NZIVGrfNovVTq26",
    "sections": [
      {
        "id": "BCdZgXFYS1wIZTOhwKI9u_s0",
        "startPoint": {
          "x": 537,
          "y": 62
        },
        "endPoint": {
          "x": 537,
          "y": 262
        },
        "bendPoints": [
          {
            "x": 537,
            "y": 62
          },
          {
            "x": 537,
            "y": 162
          }
        ],
        "incomingShape": "JT9nF7hFMPVyjU0ktYg8u",
        "outgoingShape": "13J4-8NZIVGrfNovVTq26"
      }
    ],
    "container": "root"
  },
  {
    "id": "kbK4RNkHpxnFittOz3dJl",
    "source": "13J4-8NZIVGrfNovVTq26",
    "target": "UpZv9xrgQFuH9phiiu9mA",
    "sections": [
      {
        "id": "kbK4RNkHpxnFittOz3dJl_s0",
        "startPoint": {
          "x": 582,
          "y": 312
        },
        "endPoint": {
          "x": 2018.25,
          "y": 512
        },
        "bendPoints": [
          {
            "x": 582,
            "y": 312
          },
          {
            "x": 1874.625,
            "y": 412
          }
        ],
        "incomingShape": "13J4-8NZIVGrfNovVTq26",
        "outgoingShape": "UpZv9xrgQFuH9phiiu9mA"
      }
    ],
    "container": "root"
  },
  {
    "id": "j5zJHygY-LGU4vkfF5wzj",
    "source": "UpZv9xrgQFuH9phiiu9mA",
    "target": "I3q-CPQSBdclDZAfORtoS",
    "sections": [
      {
        "id": "j5zJHygY-LGU4vkfF5wzj_s0",
        "startPoint": {
          "x": 2018.25,
          "y": 562
        },
        "endPoint": {
          "x": 2018.25,
          "y": 762
        },
        "bendPoints": [
          {
            "x": 2018.25,
            "y": 562
          },
          {
            "x": 2018.25,
            "y": 662
          }
        ],
        "incomingShape": "UpZv9xrgQFuH9phiiu9mA",
        "outgoingShape": "I3q-CPQSBdclDZAfORtoS"
      }
    ],
    "container": "root"
  },
  {
    "id": "cWPGWj7NINYn5T3tsjYmG",
    "source": "I3q-CPQSBdclDZAfORtoS",
    "target": "y5gOWdbgjwxQ6ftXSyiLg",
    "sections": [
      {
        "id": "cWPGWj7NINYn5T3tsjYmG_s0",
        "startPoint": {
          "x": 2018.25,
          "y": 812
        },
        "endPoint": {
          "x": 1887,
          "y": 1012
        },
        "bendPoints": [
          {
            "x": 2018.25,
            "y": 812
          },
          {
            "x": 1900.125,
            "y": 912
          }
        ],
        "incomingShape": "I3q-CPQSBdclDZAfORtoS",
        "outgoingShape": "y5gOWdbgjwxQ6ftXSyiLg"
      }
    ],
    "container": "root"
  },
  {
    "id": "7jqNgUkcOrMEFyHw-I_oV",
    "source": "I3q-CPQSBdclDZAfORtoS",
    "target": "Ce5SqLytWelETkNBueZ0n",
    "sections": [
      {
        "id": "7jqNgUkcOrMEFyHw-I_oV_s0",
        "startPoint": {
          "x": 2055.75,
          "y": 812
        },
        "endPoint": {
          "x": 2187,
          "y": 1012
        },
        "bendPoints": [
          {
            "x": 2055.75,
            "y": 812
          },
          {
            "x": 2173.875,
            "y": 912
          }
        ],
        "incomingShape": "I3q-CPQSBdclDZAfORtoS",
        "outgoingShape": "Ce5SqLytWelETkNBueZ0n"
      }
    ],
    "container": "root"
  },
  {
    "id": "Nb7xPpLaTGsuN_rZrfLFk",
    "source": "I3q-CPQSBdclDZAfORtoS",
    "target": "h7fZmnnQZnyM8mqGY2m8F",
    "sections": [
      {
        "id": "Nb7xPpLaTGsuN_rZrfLFk_s0",
        "startPoint": {
          "x": 1980.75,
          "y": 812
        },
        "endPoint": {
          "x": 1587,
          "y": 1012
        },
        "bendPoints": [
          {
            "x": 1980.75,
            "y": 812
          },
          {
            "x": 1626.375,
            "y": 912
          }
        ],
        "incomingShape": "I3q-CPQSBdclDZAfORtoS",
        "outgoingShape": "h7fZmnnQZnyM8mqGY2m8F"
      }
    ],
    "container": "root"
  },
  {
    "id": "rGvPVdQAE7fwPd3PAcQxe",
    "source": "13J4-8NZIVGrfNovVTq26",
    "target": "J2XljHjG17fSTm7UsnKuk",
    "sections": [
      {
        "id": "rGvPVdQAE7fwPd3PAcQxe_s0",
        "startPoint": {
          "x": 492,
          "y": 312
        },
        "endPoint": {
          "x": 87,
          "y": 512
        },
        "bendPoints": [
          {
            "x": 492,
            "y": 312
          },
          {
            "x": 127.5,
            "y": 412
          }
        ],
        "incomingShape": "13J4-8NZIVGrfNovVTq26",
        "outgoingShape": "J2XljHjG17fSTm7UsnKuk"
      }
    ],
    "container": "root"
  },
  {
    "id": "-F_JQ7f_ugn1qj_VnWNB2",
    "source": "13J4-8NZIVGrfNovVTq26",
    "target": "nF_bOmOXonb2DmOzqxOXo",
    "sections": [
      {
        "id": "-F_JQ7f_ugn1qj_VnWNB2_s0",
        "startPoint": {
          "x": 522,
          "y": 312
        },
        "endPoint": {
          "x": 387,
          "y": 512
        },
        "bendPoints": [
          {
            "x": 522,
            "y": 312
          },
          {
            "x": 400.5,
            "y": 412
          }
        ],
        "incomingShape": "13J4-8NZIVGrfNovVTq26",
        "outgoingShape": "nF_bOmOXonb2DmOzqxOXo"
      }
    ],
    "container": "root"
  },
  {
    "id": "gKco3PaEShVw7zuK-geDr",
    "source": "13J4-8NZIVGrfNovVTq26",
    "target": "AfeEyOE2Qj0dNCYnQ_7oe",
    "sections": [
      {
        "id": "gKco3PaEShVw7zuK-geDr_s0",
        "startPoint": {
          "x": 552,
          "y": 312
        },
        "endPoint": {
          "x": 687,
          "y": 512
        },
        "bendPoints": [
          {
            "x": 552,
            "y": 312
          },
          {
            "x": 673.5,
            "y": 412
          }
        ],
        "incomingShape": "13J4-8NZIVGrfNovVTq26",
        "outgoingShape": "AfeEyOE2Qj0dNCYnQ_7oe"
      }
    ],
    "container": "root"
  },
  {
    "id": "P-71Ti__fSrJNVgEUpMml",
    "source": "AfeEyOE2Qj0dNCYnQ_7oe",
    "target": "Kc5TZRLv-CbMIC27TMF7y",
    "sections": [
      {
        "id": "P-71Ti__fSrJNVgEUpMml_s0",
        "startPoint": {
          "x": 687,
          "y": 562
        },
        "endPoint": {
          "x": 687,
          "y": 762
        },
        "bendPoints": [
          {
            "x": 687,
            "y": 562
          },
          {
            "x": 687,
            "y": 662
          }
        ],
        "incomingShape": "AfeEyOE2Qj0dNCYnQ_7oe",
        "outgoingShape": "Kc5TZRLv-CbMIC27TMF7y"
      }
    ],
    "container": "root"
  },
  {
    "id": "k7rlIwRmgfp_gDbh4RClR",
    "source": "Kc5TZRLv-CbMIC27TMF7y",
    "target": "Cl4g9tLNXqE2R8p5t1Dtn",
    "sections": [
      {
        "id": "k7rlIwRmgfp_gDbh4RClR_s0",
        "startPoint": {
          "x": 662,
          "y": 812
        },
        "endPoint": {
          "x": 387,
          "y": 1012
        },
        "bendPoints": [
          {
            "x": 662,
            "y": 812
          },
          {
            "x": 414.5,
            "y": 912
          }
        ],
        "incomingShape": "Kc5TZRLv-CbMIC27TMF7y",
        "outgoingShape": "Cl4g9tLNXqE2R8p5t1Dtn"
      }
    ],
    "container": "root"
  },
  {
    "id": "EyKpqDMPTJwfP8b1QV9ux",
    "source": "Kc5TZRLv-CbMIC27TMF7y",
    "target": "2owVKiOQlID9xNu9s-4Z9",
    "sections": [
      {
        "id": "EyKpqDMPTJwfP8b1QV9ux_s0",
        "startPoint": {
          "x": 712,
          "y": 812
        },
        "endPoint": {
          "x": 987,
          "y": 1012
        },
        "bendPoints": [
          {
            "x": 712,
            "y": 812
          },
          {
            "x": 959.5,
            "y": 912
          }
        ],
        "incomingShape": "Kc5TZRLv-CbMIC27TMF7y",
        "outgoingShape": "2owVKiOQlID9xNu9s-4Z9"
      }
    ],
    "container": "root"
  },
  {
    "id": "KA1hp2ISMcTmW_igoi99j",
    "source": "Kc5TZRLv-CbMIC27TMF7y",
    "target": "yg1h4Vc0F0eP48PwCufJb",
    "sections": [
      {
        "id": "KA1hp2ISMcTmW_igoi99j_s0",
        "startPoint": {
          "x": 737,
          "y": 812
        },
        "endPoint": {
          "x": 1287,
          "y": 1012
        },
        "bendPoints": [
          {
            "x": 737,
            "y": 812
          },
          {
            "x": 1232,
            "y": 912
          }
        ],
        "incomingShape": "Kc5TZRLv-CbMIC27TMF7y",
        "outgoingShape": "yg1h4Vc0F0eP48PwCufJb"
      }
    ],
    "container": "root"
  },
  {
    "id": "wJsjVeEBvXrB9cMXlTlLq",
    "source": "Kc5TZRLv-CbMIC27TMF7y",
    "target": "qcPxgcI1nTUzlI9m1Nwmi",
    "sections": [
      {
        "id": "wJsjVeEBvXrB9cMXlTlLq_s0",
        "startPoint": {
          "x": 687,
          "y": 812
        },
        "endPoint": {
          "x": 687,
          "y": 1012
        },
        "bendPoints": [
          {
            "x": 687,
            "y": 812
          },
          {
            "x": 687,
            "y": 912
          }
        ],
        "incomingShape": "Kc5TZRLv-CbMIC27TMF7y",
        "outgoingShape": "qcPxgcI1nTUzlI9m1Nwmi"
      }
    ],
    "container": "root"
  },
  {
    "id": "06o-euRXwtrlSUhijdW6j",
    "source": "Kc5TZRLv-CbMIC27TMF7y",
    "target": "0ulNSgD2n8rIJ9TfcTKkE",
    "sections": [
      {
        "id": "06o-euRXwtrlSUhijdW6j_s0",
        "startPoint": {
          "x": 637,
          "y": 812
        },
        "endPoint": {
          "x": 87,
          "y": 1012
        },
        "bendPoints": [
          {
            "x": 637,
            "y": 812
          },
          {
            "x": 142,
            "y": 912
          }
        ],
        "incomingShape": "Kc5TZRLv-CbMIC27TMF7y",
        "outgoingShape": "0ulNSgD2n8rIJ9TfcTKkE"
      }
    ],
    "container": "root"
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
