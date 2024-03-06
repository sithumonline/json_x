import React from "react";
import ReactFlow from "reactflow";

import styles from "./App.module.css";

import "reactflow/dist/style.css";

const initialNodes = [
  {
    id: "1",
    position: { x: 0, y: 0 },
    data: {
      label: (
        <div>
          <span style={{ color: "red" }}>hello:</span>
          <span style={{ color: "blue" }}> world</span>
          <br />
          <span style={{ color: "red" }}>some:</span>
          <span style={{ color: "blue" }}> text</span>
        </div>
      ),
    },
  },
  { id: "2", position: { x: 0, y: 100 }, data: { label: "2" } },
  { id: "3", position: { x: 200, y: 0 }, data: { label: "3" } },
];
const initialEdges = [
  { id: "e1-2", source: "1", target: "2" },
  { id: "e2-3", source: "2", target: "3" },
  { id: "e3-1", source: "3", target: "1" },
];

export default function App() {
  return (
    <div className={styles.container}>
      <header className={styles.header}>React Flow - Vite Example</header>
      <div style={{ width: "100vw", height: "100vh" }}>
        <ReactFlow nodes={initialNodes} edges={initialEdges} />
      </div>
    </div>
  );
}
