<template>
  <div class="Links">
    <Nav></Nav>
    <div class="canvas" ref="can">
      <div class="txt">
        <h2>友情链接</h2>
        <p>真常应物,真常得性;常应常静,常清静矣</p>
      </div>
      <Yinshi></Yinshi>
    </div>
    <div class="link-main">
      <div class="rule">
        <h3>链接申请说明</h3>
        <p>
          <i class="fa fa-close"></i>经常宕机&nbsp;
          <i class="fa fa-close"></i>不合法规&nbsp;
          <i class="fa fa-close"></i>插边球站&nbsp;
          <i class="fa fa-close"></i>红标报毒&nbsp;
          <i class="fa fa-check"></i>原创优先&nbsp;
          <i class="fa fa-check"></i>技术优先
        </p>
        <p>交换友链可在留言板留言.本站链接如下：<br>
        名称：楚岚<br>
        网址：https://www.chulan.fun<br>
        图标：https://www.chulan.fun<br>
        描述：楚岚·一个认真的普通人</p>

      </div>
      <div class="container">
        <ul>
          <li v-for="item in linksList">
            <a :href="item.href" target="_blank">
              <img src="../../assets/img/wy.jpg">
              <h3>{{item.name}}</h3>
              <p>{{item.des}}</p>
            </a>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script>
  import Nav from "../../components/Nav";
  import Yinshi from "../../components/Yinshi";
  export default {
    name: "Links",
    data(){
      return {
        linksList:[
          {name:"欢迎加入",href:"https://www.baidu.com",icon:"",des:"热爱并坚持使用.欢迎分享"},
          {name:"欢迎加入",href:"https://www.baidu.com",icon:"",des:"热爱并坚持使用.欢迎分享"},
          {name:"欢迎加入",href:"https://www.baidu.com",icon:"",des:"热爱并坚持使用.欢迎分享"},
          {name:"欢迎加入",href:"https://www.baidu.com",icon:"",des:"热爱并坚持使用.欢迎分享"},
          {name:"欢迎加入",href:"https://www.baidu.com",icon:"",des:"热爱并坚持使用.欢迎分享"},
          {name:"欢迎加入",href:"https://www.baidu.com",icon:"",des:"热爱并坚持使用.欢迎分享"},
          {name:"欢迎加入",href:"https://www.baidu.com",icon:"",des:"热爱并坚持使用.欢迎分享"},
          {name:"欢迎加入",href:"https://www.baidu.com",icon:"",des:"热爱并坚持使用.欢迎分享"},
          {name:"欢迎加入",href:"https://www.baidu.com",icon:"",des:"热爱并坚持使用.欢迎分享"},
          {name:"欢迎加入",href:"https://www.baidu.com",icon:"",des:"热爱并坚持使用.欢迎分享"},
          {name:"欢迎加入",href:"https://www.baidu.com",icon:"",des:"热爱并坚持使用.欢迎分享"},
          {name:"欢迎加入",href:"https://www.baidu.com",icon:"",des:"热爱并坚持使用.欢迎分享"},
          {name:"欢迎加入",href:"https://www.baidu.com",icon:"",des:"热爱并坚持使用.欢迎分享"},
          {name:"欢迎加入",href:"https://www.baidu.com",icon:"",des:"热爱并坚持使用.欢迎分享"},
          {name:"欢迎加入",href:"https://www.baidu.com",icon:"",des:"热爱并坚持使用.欢迎分享"},
          {name:"欢迎加入",href:"https://www.baidu.com",icon:"",des:"热爱并坚持使用.欢迎分享"}
        ]
      }
    },
    components:{
      Nav,Yinshi
    }
  }
</script>

<style scoped lang="less">
  .Links{
    width: 100%;
    padding-top: 61px;
    .canvas{
      position: relative;
      width: 100%;
      height: 260px;
      .txt{
        position: absolute;
        left: 50%;
        top: 50%;
        text-align: center;
        transform: translate(-50%,-50%);
        color: #fff;
        h2{
          font-size: 25px;
          font-weight: normal;
        }
      }
    }
    .link-main{
      box-sizing: border-box;
      width: 100%;
      max-width: 1360px;
      padding: 30px 50px;
      /*height: 500px;*/
      margin: 0 auto;
      .rule{
        box-sizing: border-box;
        padding: 20px;
        width: 100%;
        background-color: #fff;
        h3{
          margin-bottom: 1rem;
          padding-left: 1rem;
          border-left: .25rem solid #6bc30d;
          font-weight: 700;
          font-size: 1.5rem;
          line-height: 1.5rem;
        }
        p{
          margin: 1.5rem 0;
          padding-left: 18px;
          font-family: Raleway,Segoe UI,Lucida Grande,Helvetica,Arial,Microsoft YaHei;
          line-height: 1.75rem;
          color: #515250;
        }
      }
      .container{
        width: 100%;
        ul{
          display: flex;
          width: 100%;
          flex-wrap: wrap;
          justify-content: space-between;
          li{
            box-sizing: border-box;
            width: 23.5%;
            height: 150px;
            background-color: #fff;
            margin-top: 2%;
            transition: .5s;
            cursor: pointer;
            &:hover{
              transform: translateY(-2px);
              box-shadow: 0 0 5px #fff;
            }
            a{
              display: block;
              padding: 20px;
              background: #fff;
              color: #333;
              img{
                float: left;
                width: 40px;
                height: 40px;
                border-radius: 100%;
              }
              h3{
                overflow: hidden;
                padding: 10px 15px 9px;
                height: 21px;
                font-weight: normal;
              }
              p{
                clear: both;
                display: -webkit-box;
                overflow: hidden;
                margin-top: 1rem;
                height: 60px;
                text-overflow: ellipsis;
                font-size: 13px;
                -webkit-box-orient: vertical;
                -webkit-line-clamp: 3;
              }
            }
          }
        }

      }
    }
  }
</style>