const mongoose = require("mongoose");

let Schema = mongoose.Schema;
let articleSchema = new Schema({
  type : {type: String,required: true},
  title : {type: String,required: true},
  content : {type:String,required: true},
  tag : {type: String,required: true},
  updateDate : {type: Date,default: Date.now},
  date : {type: Date, default: Date.now},
  surface : {type:String,default: 'http://47.96.127.142:80/img/defaultSurface.jpg'},
  pv : {type: Number,default: 0},
  comment : [{type:Schema.Types.ObjectID,ref:"comment"}]
})

articleSchema.pre("update",function(next){
  this.update = new Date;
  next();
});

let article = mongoose.model("article",articleSchema);


//假数据
/* for (let i=0;i<50;i++){
  article.create({
    type : ["原创","转载"][(Math.random()*2)|0],
    title : `第${i+1}篇文章`,
    content : (""+i+i+i+i+i+"哈哈哈哈哈我是文章的内容").repeat(10),
    pv : Math.random()*99999|0,
    tag : ["HTML&Css","JavaScript","Node","Vue&React","Other"][(Math.random()*5)|0]
  });
} */


module.exports = article;