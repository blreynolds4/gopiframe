import http from "../http-common";

class UploadFilesService {
  upload(file, onUploadProgress) {
    // uploads one file in a post via form data
    // returns a promise and calls onUploadProgress function
    let formData = new FormData();

    formData.append("uploadImages", file);

    let p = http.post("/photos", formData, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
      onUploadProgress,
    });

    console.log("Upload promise: ", p);
    return p;
  }
}

export default new UploadFilesService();
