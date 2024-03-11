import React, { useCallback } from "react";
import {
  ReactFlow,
  Controls,
  Background,
  BackgroundVariant,
  useNodesState,
  useEdgesState,
} from "reactflow";
import ELK from "elkjs/lib/elk.bundled.js";
import JsonVisNode from "./Custom-Nodes/JsonVisNode";
import JsonVisEdge from "./Custom Edges/JsonVisEdge";

import "./App.css";

import "reactflow/dist/style.css";

const nodeTypes = { jsonVis: JsonVisNode };
const edgeTypes = {
  jsonVis: JsonVisEdge,
};

const defaultEdgeOpt = { type: "jsonVis" };

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
        label: (
          <div>
            {objectKeys.map((key) => (
              <div>
                <span style={{ color: "red" }}>{key}:</span>
                <span style={{ color: "blue" }}> {_object[key]}</span>
              </div>
            ))}
          </div>
        ),
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

      console.log("keu", key, "value", value);
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

export default function App() {
  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges);

  // const onConnect = useCallback(
  //   (params) => setEdges((eds) => addEdge(params, eds)),
  //   [setEdges]
  // );

  const onConnect = useCallback(
    getLayoutedElements(initialNodes, initialEdges, {
      "elk.direction": "RIGHT",
      ...elkOptions,
    })
      .then(({ nodes: layoutedNodes, edges: layoutedEdges }) => {
        //add layouted or repositioned nodes and edges to store, so that react flow will render the layouted or repositioned nodes and edges
        setNodes(layoutedNodes);
        setEdges(layoutedEdges);
      })
      .catch(console.error),
    [setEdges]
  );

  return (
    <>
      <div className="app-cont">
        <div className="react-flow-cont">
          <ReactFlow
            nodes={nodes}
            edges={edges}
            onNodesChange={onNodesChange}
            onEdgesChange={onEdgesChange}
            onConnect={onConnect}
            nodeTypes={nodeTypes}
            edgeTypes={edgeTypes}
            defaultEdgeOptions={defaultEdgeOpt}
          >
            <Background
              gap={30}
              color={"#373737"}
              variant={BackgroundVariant.Lines}
            ></Background>
            <Controls showInteractive={false}></Controls>
          </ReactFlow>
        </div>
      </div>
    </>
  );
}
