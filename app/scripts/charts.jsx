var SVGNS = "http://www.w3.org/2000/svg",
    Router = ReactRouter;

var BarChart = React.createClass({
    mixins: [Router.Navigation, Router.State],
    barHeight: 40,
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

    onFilter: function(thing, i) {
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
        console.log(s);
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
        console.log("State:", this.state)
        return (
            <div className="barchart-container">
                <div className="filters">
                    <Selector thing="item"
                        items={this.props.items}
                        value={this.state.item}
                        onChange={this.onFilter.bind(this, 'item')} />
                    <Selector thing="sort"
                        items={['commits', 'delta']}
                        value={this.state.sort}
                        onChange={this.onFilter.bind(this, 'sort')} />
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
            x = (min >= 0 ? 0 : -min/max2*maxWidth - (val >= 0 ? 0 : width)),
            y = this.y(i);
            console.log(point.item, {val: val, max: max, x: x, y: y, width: width})

        return (
            <Bar key={point.item} point={point} i={i}
                metric={this.state.sort}
                x={x} y={y} width={width} height={height}
                link={this.props.link} />
        );
    }
});

var Bar = React.createClass({
    mixins: [Router.Navigation],
    handleClick: function(e) {
        this.transitionTo(this.props.link + this.props.point.item);
    },

    render: function() {
        var p = this.props.point,
            val = p[this.props.metric],
            w = this.props.width,
            label = p.item + ': ' + val,
            labelm = 10, // Margin
            labelw = label.length*9.3 + 2*labelm, // Width
            textx = labelm;
        if (labelw + 2*labelm > w) {
            textx = w + textx;
        }
        return (
            <g onClick={this.handleClick}>
                <rect className="bar" fill={Colors2[this.props.i]}
                    width={w} height={this.props.height}
                    x={this.props.x} y={this.props.y} rx="2" ry="2" />
                <rect className="label_underlay"
                    x={textx-6} y={this.props.y+10}
                    height={20} width={labelw}
                    rx="3" ry="3" />
                <text className="label" x={textx} y={this.props.y + 26}>{label}</text>
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
