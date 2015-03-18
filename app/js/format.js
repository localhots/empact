var monthNames = [
    'Jan', 'Feb', 'Mar',
    'Apr', 'May', 'Jun',
    'Jul', 'Aug', 'Sep',
    'Oct', 'Nov', 'Dec'
];

function formatDate(ts, showYear) {
    var d = new Date(ts*1000),
        day = d.getDate(),
        month = monthNames[d.getMonth()],
        year = (''+ d.getFullYear()).slice(2);

    if (showYear) {
        return month +' '+ day +" '"+ year;
    } else {
        return month +' '+ day;
    }
}

function formatNumber(num) {
    x = (''+ num).split('.');
    x1 = x[0];
    x2 = x.length > 1 ? '.' + x[1] : '';
    while (/(\d+)(\d{3})/.test(x1)) {
        x1 = x1.replace(/(\d+)(\d{3})/, '$1,$2');
    }
    return (x1 + x2).replace('-', '−');
}
