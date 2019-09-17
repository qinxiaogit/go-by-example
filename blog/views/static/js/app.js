(function($) {
  'use strict';
  $(function() {
    $('#admin-fullscreen').on('click', function() {
      $.AMUI.fullscreen.toggle();
    });
    axios.defaults.baseURL = 'https://api.owletmiao.cn'
    axios.get('/api/v1/articles').then(function (response) {
          // console.log(response.data);
        let code = response.data.code
        if(!!!code){
          let list = response.data.data.lists
            for (let i=0;i<list.length;i++){
              let item = list[i]
              $(".am-list").append(
                  "<li><a href='article/index.html?article_id="+item.id+"'>"+item.title+"-<small>"+formatDate(item.created_on)+"</small></a></li>"
              )
            }
        }
    }).catch(function (error) {
          console.log(error);
    });
  });
})(window.Zepto);

function getQueryString(name) {
  let reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)", "i");
  let r = window.location.search.substr(1).match(reg);
  if (r != null) return unescape(r[2]); return null;
}
