import React from "react";
import ReactFlow from "reactflow";

import styles from "./App.module.css";

import "reactflow/dist/style.css";

const initialNodes = {{ .InitialNodes }}

const initialEdges = {{ .InitialEdges }}

export default function App() {
  return (
    <div className={styles.container}>
      <header className={styles.header}>React Flow - Vite Example</header>
      <div className={styles.flow}>
        <ReactFlow nodes={initialNodes} edges={initialEdges} />
      </div>
    </div>
  );
}
