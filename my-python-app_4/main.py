import uuid

import dearpygui.dearpygui as dpg
import matplotlib.pyplot as plt
import networkx as nx


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


def json_node_to_nodes_and_links(node, G):
    if node.key:
        G.add_node(node.id, label=node.key)
    elif node.values:
        G.add_node(node.id, label="\n".join(f"{k}: {v}" for k, v in node.values.items()))
    else:
        G.add_node(node.id)

    if node.parent_id:
        G.add_edge(node.parent_id, node.id)

    for key, value in node.lists.items():
        G.add_node(node.id + "_" + key, label=key)
        G.add_edge(node.id, node.id + "_" + key)

        for i, item in enumerate(value):
            G.add_node(node.id + f"_{key}_{i}", label=str(item))
            G.add_edge(node.id + "_" + key, node.id + f"_{key}_{i}")

    if node.values and node.key:
        G.add_node(node.id + "_values", label="\n".join(f"{k}: {v}" for k, v in node.values.items()))
        G.add_edge(node.id, node.id + "_values")

    for child in node.children:
        json_node_to_nodes_and_links(child, G)


# callback runs when user attempts to connect attributes
def link_callback(sender, app_data):
    # app_data -> (link_id1, link_id2)
    dpg.add_node_link(app_data[0], app_data[1], parent=sender)


# callback runs when user attempts to disconnect attributes
def delink_callback(sender, app_data):
    # app_data -> link_id
    dpg.delete_item(app_data)


def create_nodes2(node, parent=None):
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
        create_nodes2(child, parent=key_node)


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
    "name": "Example",
    "active": True,
    "count": 42,
    "nested": {
        "key": "value",
        "list": [1, 2, 3],
        "nested2": {
            "key": "value"
        }
    },
    "list": [
        {"key": "value"},
        {"key": "value"}
    ]
}

root_node = parse_json(json_data)
print(root_node)

G = nx.DiGraph()
json_node_to_nodes_and_links(root_node, G)
pos = nx.nx_agraph.graphviz_layout(G, prog="dot")
# print(pos)
nx.draw(G, pos, with_labels=True)
plt.show()

dpg.create_context()

window = dpg.add_window(label="Json_x", width=800, height=600)
node_editor = dpg.add_node_editor(parent=window, callback=link_callback, delink_callback=delink_callback)
create_nodes2(root_node)

dpg.create_viewport(title='Json X', width=800, height=600)
dpg.setup_dearpygui()
dpg.show_viewport()
dpg.start_dearpygui()
dpg.destroy_context()
