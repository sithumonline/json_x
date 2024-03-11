import React, { useCallback } from "react";
import {
  ReactFlow,
  Controls,
  Background,
  BackgroundVariant,
  useNodesState,
  useEdgesState,
  addEdge,
} from "reactflow";
import JsonVisNode from "./Custom-Nodes/JsonVisNode";
import JsonVisEdge from "./Custom Edges/JsonVisEdge";

import "./App.css";

import "reactflow/dist/style.css";

const nodeTypes = { jsonVis: JsonVisNode };
const edgeTypes = {
  jsonVis: JsonVisEdge,
};

const defaultEdgeOpt = { type: "jsonVis" };

const initialNodes = {{ .InitialNodes }}

const initialEdges = {{ .InitialEdges }}

export default function App() {
  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges);

  const onConnect = useCallback(
    (params) => setEdges((eds) => addEdge(params, eds)),
    [setEdges]
  );

  return (
    <>
      <div className="app-cont">
        <div className="react-flow-cont"
        >
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
