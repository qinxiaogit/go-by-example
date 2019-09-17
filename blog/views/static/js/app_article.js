(function($) {
    'use strict';
    $(function() {
        $('#admin-fullscreen').on('click', function() {
            $.AMUI.fullscreen.toggle();
        });
        axios.defaults.baseURL = 'https://api.owletmiao.cn'
        axios.get('/api/v1/articles/'+getQueryString("article_id")).then(function (response) {
            let code = response.data.code
            console.log(response.data)
            if(!!!code){
                let content = response.data.data.content
                let title = response.data.data.title
                $("#article_detail").html(marked(content))
                $(".am-article-title").html(title)
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
