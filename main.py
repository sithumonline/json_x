import uuid
import json

class JsonNode:
    def __init__(self, key=None, parent_id=None, level=0, index=0):
        self.id = str(uuid.uuid4())[:8]  # Shortened for readability
        self.key = key
        self.parent_id = parent_id
        self.level = level
        self.index = index
        self.values = {}  # Simple key-value pairs
        self.lists = {}  # Lists of simple values
        self.children = []  # Nested JsonNode objects

    def add_value(self, key, value):
        self.values[key] = value

    def add_list(self, key, value):
        if key not in self.lists:
            self.lists[key] = []
        self.lists[key].append(value)

    def add_child(self, child):
        self.children.append(child)

    def __repr__(self):
        return f"JsonNode(id={self.id}, key={self.key}, values={self.values}, lists={self.lists}, children=[{', '.join(repr(child) for child in self.children)}])"

def parse_json(data, key=None, parent_id=None, level=0, index=0):
    # print(f"Parsing level: {level}, key: {key}")  # Debugging print
    node = JsonNode(key, parent_id, level, index)

    if isinstance(data, dict):
        # print(f"Processing dict: {data}")  # Debugging print
        for k, v in data.items():
            if isinstance(v, dict):
                child = parse_json(v, k, node.id, level+1, 0)
                node.add_child(child)
            elif isinstance(v, list):
                process_list(v, node, k, level+1)
            else:
                node.add_value(k, v)
    # elif isinstance(data, list):
    #     process_list(data, node, key, level)
    # else:
    #     node.add_value(key, data)
    # print(f"Returning node: {node}")  # Debugging print
    return node

def process_list(lst, parent_node, key, level):
    # print(f"Processing list: {lst}")  # Debugging print
    for i, item in enumerate(lst):
        if isinstance(item, dict):
            child = parse_json(item, key, parent_node.id, level, i)
            parent_node.add_child(child)
        else:
            parent_node.add_list(key, item)

# Example usage
if __name__ == "__main__":
    json_data = {
        "name": "Example",
        "active": True,
        "tags": ["demo", "json"],
        "details": {
            "description": "A simple JSON example.",
            "type": "example"
        },
        "items": [
            {"name": "Item 1", "value": 100},
            {"name": "Item 2", "value": 150}
        ]
    }

    json_data = {
        "squadName": "Super hero squad",
        "homeTown": "Metro City",
        "formed": 2016,
        "secretBase": "Super tower",
        "active": True,
        "powers": ["Million tonne punch", "Damage resistance", "Superhuman reflexes"],
        "soo": {
            "name": "Madame Uppercut",
            "age": 39,
            "secretIdentity": "Jane Wilson",
        },
        "members": [
            {
                "name": "Molecule Man",
                "age": 29,
                "secretIdentity": "Dan Jukes",
                "powers": ["Radiation resistance", "Turning tiny", "Radiation blast"],
            },
            {
                "name": "Madame Uppercut",
                "age": 39,
                "secretIdentity": "Jane Wilson",
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
                    "Interdimensional travel",
                ],
            },
        ],
    }

    # read json from file
    filePath = "/home/sithum/Downloads/n.json"
    filePath = "/home/sithum/Downloads/n.json"
    with open(filePath, "r") as file:
        json_data = json.load(file)

    # print(json_data)
    parsed_json = parse_json(json_data)
    print(parsed_json)


