try {
    new ActiveXObject("Msxml2.XMLHTTP");
    getXHR = function() {
        return new ActiveXObject("Msxml2.XMLHTTP");
    }
} catch (e) {
    try {
        new ActiveXObject("Microsoft.XMLHTTP");
        getXHR = function() {
            return new ActiveXObject("Microsoft.XMLHTTP");
        }
    } catch (e) {
        if (typeof XMLHttpRequest !== 'undefined') {
            getXHR = function() {
                return new XMLHttpRequest();
            }
        } else {
            alert("Something went really wrong!");
            console.log("XHR is not available");
        }
    }
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
