import json
import uuid

import dearpygui.dearpygui as dpg
from collections import defaultdict, deque


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
        return (f"JsonNode(id={self.id}, key={self.key}, values={self.values}, lists={self.lists}, "
                f"children=[{', '.join(repr(child) for child in self.children)}])")


def parse_json(data, key=None, parent_id=None, level=0, index=0):
    node = JsonNode(key, parent_id, level, index)

    if isinstance(data, dict):
        for k, v in data.items():
            if isinstance(v, dict):
                child = parse_json(v, k, node.id, level + 1, 0)
                node.add_child(child)
            elif isinstance(v, list):
                process_list(v, node, k, level + 1)
            else:
                node.add_value(k, v)
    return node


def process_list(lst, parent_node, key, level):
    for i, item in enumerate(lst):
        if isinstance(item, dict):
            child = parse_json(item, key, parent_node.id, level, i)
            parent_node.add_child(child)
        else:
            parent_node.add_list(key, item)


def json_node_to_nodes_and_links(node, edges):
    if node.parent_id:
        edges.append((node.parent_id, node.id))

    for key, value in node.lists.items():
        edges.append((node.id, node.id + "_" + key))

        for i, item in enumerate(value):
            edges.append((node.id + "_" + key, node.id + f"_{key}_{i}"))

    if node.values and node.key:
        edges.append((node.id, node.id + "_values"))

    for child in node.children:
        json_node_to_nodes_and_links(child, edges)


# callback runs when user attempts to connect attributes
def link_callback(sender, app_data):
    # app_data -> (link_id1, link_id2)
    dpg.add_node_link(app_data[0], app_data[1], parent=sender)


# callback runs when user attempts to disconnect attributes
def delink_callback(sender, app_data):
    # app_data -> link_id
    dpg.delete_item(app_data)


def file_dialog_callback(sender, app_data):
    print(f"File Dialog Callback: {app_data}")
    json_file = app_data["file_path_name"]
    current_filter = app_data["current_filter"]
    if current_filter != ".json":
        print("Invalid file type")
        return
    with open(json_file) as f:
        data = json.load(f)
    node_graph, max_x, max_y, pos = get_graph(data)
    create_nodes(node_graph, max_x, max_y, pos)


def create_nodes2(node: JsonNode, parent=None, pos: dict = None, node_editor=None):
    if node.key:
        with dpg.node(parent=node_editor, pos=(pos[node.id][0], pos[node.id][1])):
            with dpg.node_attribute(attribute_type=dpg.mvNode_Attr_Output) as key_node:
                dpg.add_text(node.key)
            with dpg.node_attribute(attribute_type=dpg.mvNode_Attr_Input) as key_node_:
                pass
    elif node.values:
        with dpg.node(parent=node_editor, pos=(pos[node.id][0], pos[node.id][1])):
            with dpg.node_attribute(attribute_type=dpg.mvNode_Attr_Output) as key_node:
                txt = "\n".join(f"{k}: {v}" for k, v in node.values.items())
                dpg.add_text(txt)
            with dpg.node_attribute(attribute_type=dpg.mvNode_Attr_Input) as key_node_:
                pass
    else:
        with dpg.node(parent=node_editor, pos=(pos[node.id][0], pos[node.id][1])):
            with dpg.node_attribute(attribute_type=dpg.mvNode_Attr_Output) as key_node:
                pass
            with dpg.node_attribute(attribute_type=dpg.mvNode_Attr_Input) as key_node_:
                pass

    if parent:
        dpg.add_node_link(parent, key_node_, parent=node_editor)

    for key, value in node.lists.items():
        with dpg.node(parent=node_editor, pos=(pos[node.id + "_" + key][0], pos[node.id + "_" + key][1])):
            with dpg.node_attribute(attribute_type=dpg.mvNode_Attr_Input) as key_input:
                dpg.add_node_link(key_node, key_input, parent=node_editor)

            with dpg.node_attribute(attribute_type=dpg.mvNode_Attr_Output) as key_output:
                dpg.add_text(key)

        for i, item in enumerate(value):
            with dpg.node(parent=node_editor, pos=(pos[node.id + f"_{key}_{i}"][0], pos[node.id + f"_{key}_{i}"][1])):
                with dpg.node_attribute(attribute_type=dpg.mvNode_Attr_Input) as key_input_:
                    dpg.add_text(str(item))
                    dpg.add_node_link(key_output, key_input_, parent=node_editor)

    if node.values and node.key:
        with dpg.node(parent=node_editor, pos=(pos[node.id + "_values"][0], pos[node.id + "_values"][1])):
            with dpg.node_attribute(attribute_type=dpg.mvNode_Attr_Input) as key_value_node:
                txt = "\n".join(f"{k}: {v}" for k, v in node.values.items())
                dpg.add_text(txt)

        dpg.add_node_link(key_node, key_value_node, parent=node_editor)

    for child in node.children:
        create_nodes2(child, key_node, pos, node_editor)

def layout_tree(edges, direction, nodeSpacing, levelSpacing):
    # Build the tree structure
    children = defaultdict(list)
    in_degree = defaultdict(int)
    for u, v in edges:
        children[u].append(v)
        in_degree[v] += 1

    # Identify the root node(s)
    roots = [node for node in children if in_degree[node] == 0]

    # Initialize positions dict
    positions = {}

    # Function to calculate positions
    def bfs_layout():
        queue = deque([(root, 0, 0) for root in roots]) # node, level, position within level
        level_width = defaultdict(int)
        while queue:
            node, level, pos = queue.popleft()
            # Calculate x and y based on direction
            if direction == 'top-to-bottom':
                x = pos * nodeSpacing
                y = level * levelSpacing
            else:  # 'left-to-right'
                x = level * levelSpacing
                y = pos * nodeSpacing
            positions[node] = (x, y)
            for child in children[node]:
                level_width[level + 1] += 1
                queue.append((child, level + 1, level_width[level + 1]))
    
    bfs_layout()
    return positions

""""
with dpg.window(label="Tutorial", width=400, height=400) as window:

    with dpg.node_editor(callback=link_callback, delink_callback=delink_callback) as node_editor:
        # with dpg.node(label="Node 1", parent=node_editor) as node1:
        with dpg.node(label="Node 1") as node1:
            with dpg.node_attribute(label="Node A1"):
                dpg.add_input_float(label="F1", width=150)
                dpg.add_text("Hello Node A1")

            with dpg.node_attribute(label="Node A2", attribute_type=dpg.mvNode_Attr_Output) as node_a2_out:
                dpg.add_input_float(label="F2", width=150)
                dpg.add_text("Hello Node A2 Output")

                # with dpg.node(label="Node 1.1", parent=node1) as node1_1:
        with dpg.node(label="Node 1.1") as node1_1:
            with dpg.node_attribute(label="Node A1.1"):
                dpg.add_input_float(label="F1.1", width=150)
                dpg.add_text("Hello Node A1.1")

            with dpg.node_attribute(label="Node A2.1", attribute_type=dpg.mvNode_Attr_Input) as node_a2_1_in:
                dpg.add_input_float(label="F2.1", width=150)
                dpg.add_text("Hello Node A2.1 Output")

        print(node_a2_out, node_a2_1_in)
        dpg.add_node_link(node_a2_out, node_a2_1_in, parent=node_editor)

        with dpg.node(label="Node 2"):
            with dpg.node_attribute(label="Node A3"):
                dpg.add_input_float(label="F3", width=200)

            with dpg.node_attribute(label="Node A4", attribute_type=dpg.mvNode_Attr_Output):
                dpg.add_input_float(label="F4", width=200)
"""

json_data = {
    "squadName": "Super hero squad",
    "homeTown": "Metro City",
    "formed": 2016,
    "secretBase": "Super tower",
    "active": True,
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
        {
            "name": "Madame Uppercut",
            "age": 39,
            "secretIdentity": "Jane Wilson",
            "powers": [
                "Million tonne punch",
                "Damage resistance",
                "Superhuman reflexes"
            ]
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


def get_graph(json_in: dict) -> tuple:
    node_graph = parse_json(json_in)
    edges = []
    json_node_to_nodes_and_links(node_graph, edges)
    pos = layout_tree(edges, direction='top-to-bottom', nodeSpacing=200, levelSpacing=200)
    max_x = max(x for x, y in pos.values())
    max_y = max(y for x, y in pos.values())
    return node_graph, max_x, max_y, pos


def create_nodes(node_graph, max_x, max_y, pos):
    dpg.delete_item("node_editor")
    node_editor = dpg.add_node_editor(parent=window, callback=link_callback, delink_callback=delink_callback,
                                      # minimap=True,
                                      width=max_x + 100, height=max_y + 100, tag="node_editor")
    create_nodes2(node=node_graph, pos=pos, node_editor=node_editor)


dpg.create_context()

dpg.create_viewport(title='Json X')

# file dialog for opening json file
with dpg.file_dialog(label="Open File", callback=file_dialog_callback, directory_selector=False,
                     show=False, width=700, height=400) as file_dialog:
    dpg.add_file_extension(".*")
    dpg.add_file_extension("", color=(150, 255, 150, 255))
    dpg.add_file_extension(".json", color=(150, 255, 150, 255), custom_text="[JSON]")

window = dpg.add_window(label="Json_x", width=dpg.get_viewport_width(), height=dpg.get_viewport_height())
dpg.add_button(label="Open File", callback=lambda: dpg.show_item(file_dialog), parent=window)

node_graph, max_x, max_y, pos = get_graph(json_data)
create_nodes(node_graph, max_x, max_y, pos)

dpg.setup_dearpygui()
dpg.show_viewport()
dpg.start_dearpygui()
dpg.destroy_context()
