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
