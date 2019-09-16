function formatDate(now) {
     now = new Date(now*1000);
    let   year=now.getFullYear();
    let   month=now.getMonth()+1;
    month = month<10?'0'+month:month
    let   date=now.getDate();
    date = date<10?'0'+date:date
    let   hour=now.getHours();
    hour = hour<10?'0'+hour:hour
    let   minute=now.getMinutes();
    minute = minute<10?'0'+minute:minute
    let   second=now.getSeconds();
    second = second<10?'0'+second:second
    return   year+"-"+month+"-"+date+"   "+hour+":"+minute+":"+second;
}