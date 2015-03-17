function getXHR(){
    var xhr;
    try {
        xhr = new ActiveXObject("Msxml2.XMLHTTP");
    } catch (e) {
        try {
            xhr = new ActiveXObject("Microsoft.XMLHTTP");
        } catch (e) {
            xhr = false;
        }
    }
    if (!xhr && typeof XMLHttpRequest !== 'undefined') {
        xhr = new XMLHttpRequest();
    }
    return xhr;
}

function getURL(url, params, callback) {
    if (Object.keys(params).length > 0) {
        var pairs = [];
        for (key in params) {
            pairs.push(key +'='+ encodeURIComponent(params[key]))
        }
        url += '?'+ pairs.join('&');
    }

    var xhr = getXHR();
    xhr.open('GET', url, true);
    xhr.onreadystatechange = function() {
        if (xhr.readyState === 4) {
            if(xhr.status === 200) {
                callback(JSON.parse(xhr.responseText));
            }
        }
    };
    xhr.send(null);
}
