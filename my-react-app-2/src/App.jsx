import { useState, useRef, useEffect } from "react";
import reactLogo from "./assets/react.svg";
import viteLogo from "/vite.svg";
import "./App.css";

import ELK from "elkjs/lib/elk.bundled.js";
import "litegraph.js/css/litegraph.css";
import LiteGraph from "litegraph.js";

const json_data = {
  squadName: "Super hero squad",
  homeTown: "Metro City",
  formed: 2016,
  secretBase: "Super tower",
  active: true,
  powers: ["Million tonne punch", "Damage resistance", "Superhuman reflexes"],
  soo: {
    name: "Madame Uppercut",
    age: 39,
    secretIdentity: "Jane Wilson",
  },
  members: [
    {
      name: "Molecule Man",
      age: 29,
      secretIdentity: "Dan Jukes",
      powers: ["Radiation resistance", "Turning tiny", "Radiation blast"],
    },
    {
      name: "Madame Uppercut",
      age: 39,
      secretIdentity: "Jane Wilson",
    },
    {
      name: "Eternal Flame",
      age: 1000000,
      secretIdentity: "Unknown",
      powers: [
        "Immortality",
        "Heat Immunity",
        "Inferno",
        "Teleportation",
        "Interdimensional travel",
      ],
    },
  ],
};

const initialNodes = [];
const initialEdges = [];
const list = [];

function parseJsonData(json_data, parentKey, levelIndex = 0, elementIndex = 0) {
  const _arrayALL = {};
  const _arrayStringNumber = {};
  const _object = {};

  Object.entries(json_data).forEach(([key, value]) => {
    if (!value) return;

    if (Array.isArray(value)) {
      _arrayALL[key] = value;
    } else if (typeof value === "object") {
      parseJsonData(
        value,
        key,
        elementIndex === 0 ? levelIndex : levelIndex + 1,
        elementIndex
      );
    } else {
      _object[key] = value;
    }
  });

  Object.entries(_arrayALL).forEach(([key, value]) => {
    let i = 0;
    value.forEach((element) => {
      if (typeof element === "object") {
        // console.log("object -->", i, key);
        parseJsonData(element, key, levelIndex + 1, i);
      } else {
        _arrayStringNumber[key] = value;
      }
      i++;
    });
  });

  // create node for parentKey, label = parentKey
  let parentNodeId = null;
  if (parentKey) {
    parentNodeId = `${parentKey}--object-parent`;

    // check initialNodes already have parentKey node
    let isParentNodeExist = false;
    initialNodes.forEach((node) => {
      if (node.id === parentNodeId) {
        isParentNodeExist = true;
      }
    });

    if (!isParentNodeExist) {
      initialNodes.push({
        id: parentNodeId,
        position: { x: 200 * levelIndex, y: 200 * elementIndex },
        data: { label: parentKey },
      });
    }

    if (elementIndex === 0) {
      list.push(parentNodeId);
    }
  }

  // handle object
  const objectKeys = Object.keys(_object);
  if (objectKeys.length > 0) {
    const objectNodeId = `${parentKey}-${levelIndex}-${elementIndex}-object`;

    if (parentKey) {
      initialEdges.push({
        id: `${parentKey}-${levelIndex}-${elementIndex}-object-edge`,
        source: parentNodeId,
        target: objectNodeId,
      });
    } else if (elementIndex === 0) {
      // create links for sub nodes
      list.forEach((nodeId) => {
        initialEdges.push({
          id: `${nodeId}-${objectNodeId}-edge`,
          source: objectNodeId,
          target: nodeId,
        });
      });
    }

    initialNodes.push({
      id: objectNodeId,
      position: { x: 450 + levelIndex * 200, y: 300 + elementIndex * 200 },
      data: {
        label: {
          ..._object,
        },
      },
    });

    parentNodeId = objectNodeId;
  }

  // handle array
  const arrayKeys = Object.keys(_arrayStringNumber);
  if (arrayKeys.length > 0) {
    Object.entries(_arrayStringNumber).forEach(([key, value]) => {
      // node for key label = key
      let keyNodeId = `${key}-${levelIndex}-${elementIndex}-array-key`;
      initialNodes.push({
        id: keyNodeId,
        position: { x: 50 + levelIndex * 200, y: 500 + elementIndex * 200 },
        data: { label: key },
      });
      initialEdges.push({
        id: `${key}-${levelIndex}-${elementIndex}-array-key-edge`,
        source: parentNodeId,
        target: keyNodeId,
      });

      // node for value label = value
      value.forEach((element, index) => {
        let valueNodeId = `${key}-${levelIndex}-${elementIndex}-${index}-array-value`;
        initialNodes.push({
          id: valueNodeId,
          position: {
            x: 70 + levelIndex * 200 * index,
            y: 700 + elementIndex * 200 + index,
          },
          data: { label: element },
        });
        initialEdges.push({
          id: `${key}-${levelIndex}-${elementIndex}-${index}-array-value-edge`,
          source: keyNodeId,
          target: valueNodeId,
        });
      });
    });
  }
}

parseJsonData(json_data, null);

const elk = new ELK();

/**
 *
 * @param {*} nodes array of nodes from store
 * @param {*} edges array of edges from store
 * @param {*} options options from elkOptions. Used for layouting tree
 * @returns promises that contains array of nodes or edges that already get layouted or repositioned
 */
const getLayoutedElements = (nodes, edges, options = {}) => {
  const isHorizontal = options?.["elk.direction"] === "RIGHT";
  // console.log(isHorizontal);
  const graph = {
    id: "root",
    layoutOptions: options,
    //Passed array of nodes that contains target position and source position. The target position and source position change based on isHorizontal
    children: nodes.map((node) => ({
      ...node,
      targetPosition: isHorizontal ? "left" : "top",
      sourcePosition: isHorizontal ? "right" : "bottom",
      //Hardcode a width and height for node so that elk can use it when layouting.
      width: 150,
      height: 50,
    })),
    edges: edges,
  };

  // console.log(graph);

  //Return promises
  return elk
    .layout(graph)
    .then((layoutedGraph) => ({
      nodes: layoutedGraph.children.map((node) => ({
        ...node,
        // React Flow expects a position property on the node instead of `x` and `y` fields.
        position: { x: node.x, y: node.y },
      })),
      edges: layoutedGraph.edges,
    }))
    .catch(console.error);
};

//Elk options for layouting the tree
const elkOptions = {
  "elk.algorithm": "layered",
  "elk.layered.spacing.nodeNodeBetweenLayers": "200",
  "elk.spacing.nodeNode": "150",
  "elk.edgeRouting": "SPLINES",
};

function normalNode() {
  this.addOutput("value", "");
  this.addInput("value", "");
  this.addProperty("value", "");
  this.widget = this.addWidget("text", "", "", "value");
  this.widgets_up = true;
  console.log("this.widget-normal", this);
}

normalNode.title = "Value";

normalNode.prototype.setValue = function (value) {
  this.setProperty("value", value);
};

function normalNodeObject() {
  this.addOutput("value", "");
  this.addInput("value", "");
  this.addProperty("value", { property: "" });
  // this.name_widget = this.addWidget(
  //   "text",
  //   "prop.",
  //   this.properties.value,
  //   "value"
  // );

  const inputEl = document.createElement("textarea");
  inputEl.className = "comfy-multiline-input";
  inputEl.value = "kdkdkdskdmxkmkdodcdkdmmx";
  inputEl.placeholder = "skdskdkjfdkjdkjdkkdkd";

  // const widget = this.addCustomWidget("d-name", "customtext", inputEl, {
  //   getValue() {
  //     return inputEl.value;
  //   },
  //   setValue(v) {
  //     inputEl.value = v;
  //   },
  // });
  // widget.inputEl = inputEl;

  // inputEl.className = "comfy-multiline-input";
  // inputEl.value = "opts.defaultVal";
  // inputEl.placeholder = "opts.placeholder";

  // const element = inputEl;

  this.addCustomWidget({
    type: "d-name",
    name: "customtext",
    get value() {
      return inputEl.value;
    },
    set value(v) {
      inputEl.value = v;
    },
    element: inputEl,
    options: {},
    draw: function (ctx, node, widgetWidth, y, widgetHeight) {
      // console.log("draw", ctx, node, widgetWidth, y, widgetHeight);
      ctx.fillStyle = "#AAA";
      ctx.fillRect(0, y, widgetWidth, widgetHeight);
      ctx.fillStyle = "#000";
      ctx.fillText("customtext", 0, y + 10);

      // console.log("context", ctx, "node", node, "widgetWidth", widgetWidth);

      Object.assign(inputEl.style, {
        position: "absolute",
        left: "0px",
        top: y + "px",
        width: widgetWidth + "px",
        height: widgetHeight + "px",
        outline: "red",
        background: "green",
        transformOrigin: "0 0",
        transform: "scale(0.5)",
        zIndex: 1,
      });
    },
    value: "jdjdjdjdj\nsiwen\ncndej\ndjnendndn",
  });

  // this.widget.inputEl = inputEl;

  // inputEl.addEventListener("input", () => {
  //   this.widget.callback?.(widget.value);
  // });

  this.widget = this.addWidget("d-name", "customtext", "", "value");

  this.widgets_up = true;
  this.serialize_widgets = true;

  console.log("this.widget", this);
}

normalNodeObject.title = "Object";

normalNodeObject.prototype.setValue = function (value) {
  this.setProperty("value", value);
  // this.widget.element.value = value;
  // console.log("this", this);
};

function App() {
  const canvasRef = useRef(null);

  useEffect(() => {
    const graph = new LiteGraph.LGraph();
    const canvas = new LiteGraph.LGraphCanvas(canvasRef.current, graph);
    LiteGraph.LiteGraph.registerNodeType("basic/normal", normalNode);
    LiteGraph.LiteGraph.registerNodeType(
      "basic/normalObject",
      normalNodeObject
    );

    // const inputEl = document.createElement("textarea");
    // LiteGraph.LGraphNode.prototype.addCustomWidget({
    //   type: "d-name",
    //   name: "customtext",
    //   get value() {
    //     return inputEl.value;
    //   },
    //   set value(v) {
    //     inputEl.value = v;
    //   },
    //   element: inputEl,
    //   options: {},
    // });

    getLayoutedElements(initialNodes, initialEdges, {
      "elk.direction": "RIGHT",
      ...elkOptions,
    })
      .then(({ nodes: layoutedNodes, edges: layoutedEdges }) => {
        layoutedEdges.forEach((edge) => {
          const source = layoutedNodes.find((node) => node.id === edge.source);
          const target = layoutedNodes.find((node) => node.id === edge.target);

          if (!source || !target) {
            return;
          }

          var sourceConst;
          if (!source.node) {
            if (typeof source?.data?.label === "object") {
              sourceConst =
                LiteGraph.LiteGraph.createNode("basic/normalObject");
              // console.log("sourceConst_label", source?.data?.label);
              sourceConst.setValue(
                JSON.stringify(source?.data?.label, null, 2)
              );
            } else {
              sourceConst = LiteGraph.LiteGraph.createNode("basic/normal");
              sourceConst.setValue(source?.data?.label);
            }
            sourceConst.pos = [source.position.x, source.position.y];
            // sourceConst.setValue(source?.data?.label);
            graph.add(sourceConst);
            source.node = sourceConst;
          } else {
            sourceConst = source.node;
          }

          var targetConst;
          if (!target.node) {
            if (typeof target?.data?.label === "object") {
              targetConst =
                LiteGraph.LiteGraph.createNode("basic/normalObject");
              // console.log("sourceConst_label", target?.data?.label);
              targetConst.setValue(
                JSON.stringify(target?.data?.label, null, 2)
              );
            } else {
              targetConst = LiteGraph.LiteGraph.createNode("basic/normal");
              targetConst.setValue(target?.data?.label);
            }
            targetConst.pos = [target.position.x, target.position.y];
            // targetConst.setValue(target?.data?.label);
            graph.add(targetConst);
            target.node = targetConst;
          } else {
            targetConst = target.node;
          }

          sourceConst.connect(0, targetConst, 0);
        });
      })
      .catch(console.error),
      graph.start();
  }, []);

  return (
    <canvas
      ref={canvasRef}
      width={window.innerWidth}
      height={window.innerHeight}
      style={{ border: "1px solid" }}
    ></canvas>
  );
}

export default App;
