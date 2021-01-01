import React from "react";
import "./App.css";
import "bootstrap/dist/css/bootstrap.min.css";

import UploadFiles from "./components/upload-files.component";

function App() {
  return (
    <div className="container" style={{ width: "600px" }}>
      <div className="my-2">
        <h3>Pi Photos</h3>
        <h4>Upload one or more files from your device</h4>
      </div>

      <UploadFiles />
    </div>
  );
}

export default App;
