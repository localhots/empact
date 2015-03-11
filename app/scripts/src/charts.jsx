var SVGNS = 'http://www.w3.org/2000/svg',
    fontFamily = 'Helvetica Neue',
    fontSize = '16px',
    Router = ReactRouter;

var Selector = React.createClass({
    names: {
        "repo": "Repositories",
        "team": "Teams",
        "user": "Users",
        "commits": "Commits",
        "delta": "Delta"
    },

    itemWithName: function(name) {
        for (item in this.names) {
            if (this.names[item] === name) {
                return item;
            }
        }
    },

    renderItem: function(item, i) {
        var itemClass = (item === this.props.value ? 'active' : ''),
            clickEvent = null;
        if (this.props.onChange) {
            clickEvent = this.props.onChange.bind(this, i);
        }
        return (
            <li key={item} onClick={clickEvent} className={itemClass}>{this.names[item]}</li>
        );
    },

    render: function() {
        return (
            <ul className={this.props.thing}>
                {this.props.items.map(this.renderItem)}
            </ul>
        );
    }
});

function textWidth(str) {
    var svg = document.createElementNS(SVGNS, "svg");
        text = document.createElementNS(SVGNS, "text");

    svg.width = 500;
    svg.height = 500;
    svg.style.position = 'absolute';
    svg.style.left = '-1000px';

    text.appendChild(document.createTextNode(str))
    text.style.fontFamily = fontFamily;
    text.style.fontSize = fontSize;

    svg.appendChild(text);
    document.body.appendChild(svg);
    var box = text.getBBox();
    document.body.removeChild(svg);

    return box.width;
}
