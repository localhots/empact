var SVGNS = 'http://www.w3.org/2000/svg',
    fontFamily = 'Helvetica Neue',
    fontSize = '16px',
    Router = ReactRouter;

var Chart = {
    calculateViewBoxWidth: function() {
        this.setState({
           canvasWidth: this.refs.svg.getDOMNode().offsetWidth
        });
    },

    animate: function(ref, attr, from, to) {
        var node = ref.getDOMNode(),
            anim = anim = document.createElementNS(SVGNS, 'animate');

        _.map(node.childNodes, function(el) {
            node.removeChild(el);
        });

        anim.setAttributeNS(null, 'attributeType', 'XML');
        anim.setAttributeNS(null, 'attributeName', attr);
        anim.setAttributeNS(null, 'from', from);
        anim.setAttributeNS(null, 'to', to);
        anim.setAttributeNS(null, 'dur', '350ms');
        // anim.setAttributeNS(null, 'keySplines', [this.easing, this.easing, this.easing].join(';'));
        anim.setAttributeNS(null, 'repeatCount', '1');
        node.appendChild(anim);
        anim.beginElement();
    }
};

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
