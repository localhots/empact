var SVGNS = "http://www.w3.org/2000/svg",
    Router = ReactRouter;

var BarChart = React.createClass({
    mixins: [Router.Navigation, Router.State],
    barHeight: 30,
    barMargin: 5,

    getInitialState: function() {
        return {
            item: this.props.items[0],
            sort: 'commits',
            rawData: [],
            points: [],
            min: 0,
            max: 1
        };
    },

    componentDidMount: function() {
        this.fetchData();
    },

    handleFilter: function(thing, i) {
        if (thing === 'item' && this.props.items[i] !== this.state.item) {
            this.setState({
                item: this.props.items[i]
            }, this.fetchData);
        } else if (thing === 'sort' && ['commits', 'delta'][i] !== this.state.sort) {
            this.setState({
                sort: ['commits', 'delta'][i]
            }, this.sort);
        }
    },

    handleClick: function(point) {
        var params = {org: this.getParams().org};
        params[this.state.item] = point.item;
        this.transitionTo(this.state.item, params);
    },

    fetchData: function() {
        $.get(this.props.api, this.apiParams(), function(res){
            this.setState({
                rawData: res
            }, this.sort);
        }.bind(this));
    },

    sort: function() {
        var sortFun = function(a, b) {
            return Math.abs(b[this.state.sort]) - Math.abs(a[this.state.sort]);
        }.bind(this);
        var points = this.state.rawData.sort(sortFun).slice(0, 15);

        var min = 0, max = 1;
        points.map(function(el) {
            var val = el[this.state.sort];
            if (val > max) {
                max = val;
            }
            if (val < min) {
                min = val;
            }
        }.bind(this));

        s = {
            points: points,
            min: min,
            max: max
        };
        // console.log(s);
        this.setState(s);
    },

    apiParams: function() {
        // Deep copy
        // Don't use jQuery.extend
        var params = JSON.parse(JSON.stringify(this.props.params));
        params['item'] = this.state.item;
        return params;
    },

    height: function() {
        if (this.state.points.length === 0) {
            return 0;
        } else {
            return this.y(this.state.points.length) - this.barMargin;
        }
    },

    y: function(i) {
        return i*(this.barHeight + this.barMargin);
    },

    render: function() {
        // console.log("State:", this.state)
        return (
            <div className="barchart-container">
                <div className="filters">
                    <Selector thing="item"
                        items={this.props.items}
                        value={this.state.item}
                        onChange={this.handleFilter.bind(this, 'item')} />
                    <Selector thing="sort"
                        items={['commits', 'delta']}
                        value={this.state.sort}
                        onChange={this.handleFilter.bind(this, 'sort')} />
                </div>
                <svg className="barchart" width="100%" height={this.height()}>
                    {this.state.points.map(this.renderBar)}
                </svg>
            </div>
        );
    },

    renderBar: function(point, i) {
        var maxWidth = 400,
            val = point[this.state.sort],
            min = this.state.min,
            max = this.state.max,
            max2 = (min < 0 ? max - min : max),
            width = Math.abs(val)/max2*maxWidth,
            height = this.barHeight,
            offset = -min/max2*maxWidth,
            x = (min >= 0 ? 0 : offset - (val >= 0 ? 0 : width)),
            y = this.y(i);

        return (
            <Bar key={point.item} item={point.item} value={val}
                color={Colors2[i]}
                x={x} y={y} offset={offset} width={width} height={height}
                onClick={this.handleClick.bind(this, point)} />
        );
    }
});

var Bar = React.createClass({
    mixins: [Router.Navigation],

    render: function() {
        var val = this.props.value,
            item = this.props.item,
            offset = this.props.offset,
            width = this.props.width,
            label = item + ': ' + val,
            labelPadding = 5,
            labelWidth = textWidth(label, 'Helvetica Neue', '16px') + 2*labelPadding,
            labelOuterWidth = labelWidth + 2*labelPadding,
            labelX = 0,
            barX = this.props.x;

        if (labelOuterWidth <= width) {
            if (offset > 0) {
                if (barX === offset) {
                    labelX = barX + 2*labelPadding;
                } else {
                    labelX = barX + width - labelOuterWidth + 2*labelPadding;
                }
            } else {
                labelX = barX + 2*labelPadding;
            }
        } else {
            if (labelOuterWidth <= barX) {
                labelX = barX - labelOuterWidth + 2*labelPadding;
            } else {
                labelX = barX + width + labelPadding;
            }
        }

        return (
            <g onClick={this.props.onClick}>
                <rect className="bar" fill={this.props.color}
                    width={width} height={this.props.height}
                    x={this.props.x} y={this.props.y} rx="2" ry="2" />
                <rect className="label_underlay"
                    x={labelX - labelPadding} y={this.props.y + 5}
                    height={20} width={labelWidth}
                    rx="3" ry="3" />
                <text className="label" x={labelX} y={this.props.y + 21}>{label}</text>
            </g>
        );
    }
});

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
            clickEvent = this.props.onChange.bind(this, i);
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

function textWidth(str, font, size) {
    var svg = document.createElementNS(SVGNS, "svg");
        text = document.createElementNS(SVGNS, "text");

    svg.width = 500;
    svg.height = 500;
    svg.style.position = 'absolute';
    svg.style.left = '-1000px';

    text.appendChild(document.createTextNode(str))
    text.style.fontFamily = font;
    text.style.fontSize = size;

    svg.appendChild(text);
    document.body.appendChild(svg);
    var box = text.getBBox();
    document.body.removeChild(svg);

    return box.width;
}
