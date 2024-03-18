# This Python file uses the following encoding: utf-8
import uuid

from Qt import QtWidgets

from NodeGraphQt import NodeGraph, BaseNode, NodeBaseWidget


class LabelWidget(QtWidgets.QWidget):
    def __init__(self, parent=None):
        super(LabelWidget, self).__init__(parent)
        self.label = QtWidgets.QLabel(text='hello world')

        layout = QtWidgets.QHBoxLayout(self)
        layout.setContentsMargins(0, 0, 0, 0)
        layout.addWidget(self.label)


class LabelWidgetWrapper(NodeBaseWidget):
    def __init__(self, parent=None):
        super(LabelWidgetWrapper, self).__init__(parent)
        self.set_name('label')
        self.set_custom_widget(LabelWidget())

    def get_value(self):
        widget = self.get_custom_widget()
        return widget.label.text()

    def set_value(self, value):
        widget = self.get_custom_widget()
        widget.label.setText(value)


class LabelNode(BaseNode):
    __identifier__ = 'io.github.jchanvfx'

    NODE_NAME = 'label'

    def __init__(self):
        super(LabelNode, self).__init__()
        self.add_input('in')
        self.add_output('out')

        label_widget = LabelWidgetWrapper(self.view)
        self.add_custom_widget(label_widget, tab='Custom')

    def set_text(self, text):
        widget = self.get_widget('label')
        widget.set_value(text)


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


def create_label_node(graph, name, level, index, text):
    node_item = graph.create_node('io.github.jchanvfx.LabelNode', name=name, pos=[level * 200, index * 100])
    node_item.set_text(text)
    return node_item


def create_nodes(node, parent_node=None):
    if node.key:
        key_node = create_label_node(node_graph, node.key, node.level, node.index,
                                     f'<font color="#FF0000">{node.key}</font>')
    elif node.values:
        text = '<br>'.join(
            [f'<font color="#FF0000">{k}</font>:<font color="#00FF00">{v}</font>' for k, v in node.values.items()])
        key_node = create_label_node(node_graph, f'{node.key}_value', node.level, node.index, text)
    else:
        key_node = create_label_node(node_graph, 'root', node.level, node.index, '<font color="#FF0000">root</font>')

    if parent_node:
        parent_node.set_output(0, key_node.input(0))

    for key, value in node.lists.items():
        list_node = create_label_node(node_graph, key, node.level, node.index, f'<font color="#FF0000">{key}</font>')
        key_node.set_output(0, list_node.input(0))

        for i, item in enumerate(value):
            item_node = create_label_node(node_graph, f'{key}_{i}', node.level, node.index,
                                          f'<font color="#FF0000">{item}</font>')
            list_node.set_output(0, item_node.input(0))

    if node.values and node.key:
        text = '<br>'.join(
            [f'<font color="#FF0000">{k}</font>:<font color="#00FF00">{v}</font>' for k, v in node.values.items()])
        key_value_node = create_label_node(node_graph, f'{node.key}_value', node.level, node.index, text)
        key_node.set_output(0, key_value_node.input(0))

    for child in node.children:
        create_nodes(child, parent_node=key_node)


if __name__ == '__main__':
    app = QtWidgets.QApplication([])

    node_graph = NodeGraph()
    node_graph.register_node(LabelNode)
    node_graph.widget.show()

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

    # Parse the JSON data into a tree of JsonNode objects.
    root_node = parse_json(json_data)
    print(root_node)
    create_nodes(root_node)

    app.exec_()
