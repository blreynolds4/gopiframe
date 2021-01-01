import axios from "axios";

export default axios.create({
  baseURL: process.env.REACT_APP_UPLOAD_ENDPOINT,
  headers: {
    "Content-type": "application/json",
  },
});
