(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-bb3efcee"],{"04d2":function(t,e,n){},"0737":function(t,e,n){"use strict";n.r(e);var i=function(){var t=this,e=t.$createElement,n=t._self._c||e;return n("div",{staticClass:"About"},[n("Nav"),n("div",{ref:"can",staticClass:"canvas"},[t._m(0),n("Yinshi")],1),t._m(1)],1)},r=[function(){var t=this,e=t.$createElement,n=t._self._c||e;return n("div",{staticClass:"txt"},[n("h2",[t._v("关于")]),n("p",[t._v("真常应物,真常得性;常应常静,常清静矣")])])},function(){var t=this,e=t.$createElement,i=t._self._c||e;return i("div",{staticClass:"about-main"},[i("div",{staticClass:"a-m-content"},[i("article",[i("section",[i("h1",[t._v("关于我")]),i("p",[t._v(" 毕业后在长沙做在线教育课程顾问，很迷茫。2020年决定转行做开发，然沉迷撸码，日渐消瘦。 ")]),i("span",[t._v(" 可以通过以下方式联系到我： ")]),i("ul",[i("li",[t._v("邮 箱 ："),i("a",{attrs:{href:"390226630@qq.com"}},[t._v("390226630@qq.com")])])])]),i("section",[i("h1",[t._v("关于本站")]),i("p",[t._v(" 本站建于2020年1月，主要是个人爱好写着玩。 ")]),i("span",[t._v(" 本站结构： ")]),i("ul",[i("li",[t._v("前 端 ："),i("code",[t._v("Layui + ElementUI")])]),i("li",[t._v("后 端 ："),i("code",[t._v("nodeJS + MongoDB")])])]),i("p",[t._v(" 本站采用阿里云提供的服务器ESC和存储对象OSS。 ")])]),i("section",[i("h1",[t._v("关于版权")]),i("p",[t._v(" 本站采用「 "),i("a",{attrs:{href:"https://creativecommons.org/licenses/by-nc/4.0/deed.zh",target:"_blank"}},[t._v("署名-非商业性使用 4.0 国际 (CC BY-NC 4.0)")]),t._v("」创作共享协议。 只要在使用时注明出处，那么您可以可以对本站所有原创内容进行转载、节选、二次创作，但是您不得对其用于商业目的。 ")])]),i("section",[i("h1",[t._v("特别说明")]),i("ul",[i("li",[t._v("本站文章仅代表个人观点，和任何组织或个人无关。")]),i("li",[t._v("本站前端开发代码没有考虑对IE浏览器的兼容。")])]),i("br"),i("br"),i("div",[i("img",{staticStyle:{width:"100%",height:"320px"},attrs:{src:n("d1b4")}})])])])])])}],a=n("216c"),o=n("6217"),s={name:"About",components:{Nav:a["a"],Yinshi:o["a"]}},c=s,u=(n("0c46"),n("2877")),h=Object(u["a"])(c,i,r,!1,null,"4430836f",null);e["default"]=h.exports},"0c46":function(t,e,n){"use strict";n("04d2")},1486:function(t,e,n){},"159b":function(t,e,n){var i=n("da84"),r=n("fdbc"),a=n("17c2"),o=n("9112");for(var s in r){var c=i[s],u=c&&c.prototype;if(u&&u.forEach!==a)try{o(u,"forEach",a)}catch(h){u.forEach=a}}},"17c2":function(t,e,n){"use strict";var i=n("b727").forEach,r=n("a640"),a=n("ae40"),o=r("forEach"),s=a("forEach");t.exports=o&&s?[].forEach:function(t){return i(this,t,arguments.length>1?arguments[1]:void 0)}},"185b":function(t,e,n){"use strict";n("1486")},"1dde":function(t,e,n){var i=n("d039"),r=n("b622"),a=n("2d00"),o=r("species");t.exports=function(t){return a>=51||!i((function(){var e=[],n=e.constructor={};return n[o]=function(){return{foo:1}},1!==e[t](Boolean).foo}))}},"25f0":function(t,e,n){"use strict";var i=n("6eeb"),r=n("825a"),a=n("d039"),o=n("ad6d"),s="toString",c=RegExp.prototype,u=c[s],h=a((function(){return"/a/b"!=u.call({source:"a",flags:"b"})})),f=u.name!=s;(h||f)&&i(RegExp.prototype,s,(function(){var t=r(this),e=String(t.source),n=t.flags,i=String(void 0===n&&t instanceof RegExp&&!("flags"in c)?o.call(t):n);return"/"+e+"/"+i}),{unsafe:!0})},4160:function(t,e,n){"use strict";var i=n("23e7"),r=n("17c2");i({target:"Array",proto:!0,forced:[].forEach!=r},{forEach:r})},"4de4":function(t,e,n){"use strict";var i=n("23e7"),r=n("b727").filter,a=n("1dde"),o=n("ae40"),s=a("filter"),c=o("filter");i({target:"Array",proto:!0,forced:!s||!c},{filter:function(t){return r(this,t,arguments.length>1?arguments[1]:void 0)}})},6217:function(t,e,n){"use strict";var i=function(){var t=this,e=t.$createElement,n=t._self._c||e;return n("canvas",{attrs:{id:"canvas",height:"260"}})},r=[];n("cb29"),n("4de4"),n("4160"),n("fb6a"),n("d3b7"),n("25f0"),n("159b");function a(t,e){if(!(t instanceof e))throw new TypeError("Cannot call a class as a function")}function o(t,e){for(var n=0;n<e.length;n++){var i=e[n];i.enumerable=i.enumerable||!1,i.configurable=!0,"value"in i&&(i.writable=!0),Object.defineProperty(t,i.key,i)}}function s(t,e,n){return e&&o(t.prototype,e),n&&o(t,n),t}var c=function(){function t(e,n){a(this,t),this.x=e,this.y=n,this.vx=3*Math.random()-1.5,this.vy=3*Math.random()-1.5,this.r=3*Math.random()+3,this.dead=!1,this.color="#"+Math.random().toString(16).slice(2,8)}return s(t,[{key:"render",value:function(t){t.beginPath(),t.fillStyle=this.color,t.globalCompositeOperation="lighter",t.arc(this.x,this.y,this.r,0,2*Math.PI),t.fill()}},{key:"update",value:function(){this.r*=.96,this.x+=this.vx,this.y+=this.vy,this.r<.01&&(this.dead=!0)}}]),t}(),u={name:"HelloWorld",data:function(){return{canvas:null,bubbles:[],timer:null,animateState:!0}},mounted:function(){this.init(),this.register(),this.render()},beforeDestroy:function(){clearInterval(this.timer),window.removeEventListener("resize",this.handleResize),this.canvas.removeEventListener("mousemove",this.handleMousemove),this.animateState=!1},watch:{width:function(){this.canvas.width=this.width||500},height:function(){this.canvas.height=this.height||260}},methods:{init:function(){this.canvas=document.querySelector("#canvas"),this.ctx=this.canvas.getContext("2d"),this.canvas.width=document.documentElement.offsetWidth,this.canvas.height=260,this.auto()},auto:function(){var t=this;clearInterval(this.timer),this.timer=setInterval((function(){t.bubbles.push(new c(t.canvas.width*Math.random(),t.canvas.height*Math.random()))}),30)},register:function(){window.addEventListener("resize",this.handleResize),this.canvas.addEventListener("mousemove",this.handleMousemove)},render:function(){var t=this;this.ctx.clearRect(0,0,this.canvas.width,this.canvas.height),this.bubbles.forEach((function(e){e.render(t.ctx),e.update()})),this.bubbles=this.bubbles.filter((function(t){return!t.dead})),this.animateState&&requestAnimationFrame((function(){return t.render()}))},handleResize:function(){this.canvas.width=document.documentElement.offsetWidth},handleMousemove:function(t){this.bubbles.push(new c(t.offsetX,t.offsetY))}}},h=u,f=(n("185b"),n("2877")),l=Object(f["a"])(h,i,r,!1,null,"08d3144c",null);e["a"]=l.exports},"65f0":function(t,e,n){var i=n("861d"),r=n("e8b5"),a=n("b622"),o=a("species");t.exports=function(t,e){var n;return r(t)&&(n=t.constructor,"function"!=typeof n||n!==Array&&!r(n.prototype)?i(n)&&(n=n[o],null===n&&(n=void 0)):n=void 0),new(void 0===n?Array:n)(0===e?0:e)}},"81d5":function(t,e,n){"use strict";var i=n("7b0b"),r=n("23cb"),a=n("50c4");t.exports=function(t){var e=i(this),n=a(e.length),o=arguments.length,s=r(o>1?arguments[1]:void 0,n),c=o>2?arguments[2]:void 0,u=void 0===c?n:r(c,n);while(u>s)e[s++]=t;return e}},ad6d:function(t,e,n){"use strict";var i=n("825a");t.exports=function(){var t=i(this),e="";return t.global&&(e+="g"),t.ignoreCase&&(e+="i"),t.multiline&&(e+="m"),t.dotAll&&(e+="s"),t.unicode&&(e+="u"),t.sticky&&(e+="y"),e}},b727:function(t,e,n){var i=n("0366"),r=n("44ad"),a=n("7b0b"),o=n("50c4"),s=n("65f0"),c=[].push,u=function(t){var e=1==t,n=2==t,u=3==t,h=4==t,f=6==t,l=7==t,d=5==t||f;return function(v,b,p,m){for(var g,y,_=a(v),w=r(_),E=i(b,p,3),x=o(w.length),A=0,S=m||s,C=e?S(v,x):n||l?S(v,0):void 0;x>A;A++)if((d||A in w)&&(g=w[A],y=E(g,A,_),t))if(e)C[A]=y;else if(y)switch(t){case 3:return!0;case 5:return g;case 6:return A;case 2:c.call(C,g)}else switch(t){case 4:return!1;case 7:c.call(C,g)}return f?-1:u||h?h:C}};t.exports={forEach:u(0),map:u(1),filter:u(2),some:u(3),every:u(4),find:u(5),findIndex:u(6),filterOut:u(7)}},cb29:function(t,e,n){var i=n("23e7"),r=n("81d5"),a=n("44d2");i({target:"Array",proto:!0},{fill:r}),a("fill")},d1b4:function(t,e,n){t.exports=n.p+"img/bg.2386f34d.jpg"},e8b5:function(t,e,n){var i=n("c6b6");t.exports=Array.isArray||function(t){return"Array"==i(t)}},fb6a:function(t,e,n){"use strict";var i=n("23e7"),r=n("861d"),a=n("e8b5"),o=n("23cb"),s=n("50c4"),c=n("fc6a"),u=n("8418"),h=n("b622"),f=n("1dde"),l=n("ae40"),d=f("slice"),v=l("slice",{ACCESSORS:!0,0:0,1:2}),b=h("species"),p=[].slice,m=Math.max;i({target:"Array",proto:!0,forced:!d||!v},{slice:function(t,e){var n,i,h,f=c(this),l=s(f.length),d=o(t,l),v=o(void 0===e?l:e,l);if(a(f)&&(n=f.constructor,"function"!=typeof n||n!==Array&&!a(n.prototype)?r(n)&&(n=n[b],null===n&&(n=void 0)):n=void 0,n===Array||void 0===n))return p.call(f,d,v);for(i=new(void 0===n?Array:n)(m(v-d,0)),h=0;d<v;d++,h++)d in f&&u(i,h,f[d]);return i.length=h,i}})}}]);
//# sourceMappingURL=chunk-bb3efcee.9db63c1b.js.map