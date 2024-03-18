# This Python file uses the following encoding: utf-8
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
        self.set_label('hello world')
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


if __name__ == '__main__':
    app = QtWidgets.QApplication([])

    node_graph = NodeGraph()
    node_graph.register_node(LabelNode)
    node_graph.widget.show()

    # here we create a couple nodes in the node graph.
    node_c = node_graph.create_node('io.github.jchanvfx.LabelNode', name='label')
    node_c.set_text('<font color="#FF0000">hello </font><br><font color="#00FF00">world</font>')

    app.exec_()

