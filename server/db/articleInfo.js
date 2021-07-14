const mongoose = require("mongoose");


let articleInfo = mongoose.model("articleInfo",new mongoose.Schema({
  tags : [String],
  num : Number
}));

articleInfo.create({
  tags : ["HTML&Css","JavaScript","Node","Vue&React","Other"],
  num : 100
});

module.exports = articleInfo;