const path = require("path");

module.exports = {
  entry: "./physarum_pb.js",
  output: {
    path: path.resolve(__dirname),
    filename: "bundle.js",
  },
};
