package data

var Nodes = `
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
