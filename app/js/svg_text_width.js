var fontFamily = "'Open Sans', Helvetica, Arial, sans-serif",
    fontSize = 16;

function textWidth(str) {
    var svg = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
        text = document.createElementNS('http://www.w3.org/2000/svg', 'text');

    svg.width = 500;
    svg.height = 500;
    svg.style.position = 'absolute';
    svg.style.left = '-1000px';

    text.appendChild(document.createTextNode(str))
    text.style.fontFamily = fontFamily;
    text.style.fontSize = fontSize +'px';

    svg.appendChild(text);
    document.body.appendChild(svg);
    var box = text.getBBox();
    document.body.removeChild(svg);

    return box.width;
}
