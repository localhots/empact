var BarChart = React.createClass({
    mixins: [Router.Navigation, Router.State],

    numElements: 15,
    barHeight: 30,
    barMargin: 5,

    getInitialState: function() {
        return {
            item: this.props.items[0],
            sort: 'commits',
            rawData: [],
            points: [],
            min: 0,
            max: 1,
            canvasWidth: 500
        };
    },

    calculateViewBoxWidth: function() {
        this.setState({
           canvasWidth: this.refs.svg.getDOMNode().offsetWidth
        });
    },

    componentDidMount: function() {
        this.fetchData();
        this.calculateViewBoxWidth();
        window.addEventListener('resize', this.calculateViewBoxWidth);
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
        var points = this.state.rawData.sort(sortFun).slice(0, this.numElements);

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

        this.setState({
            points: points,
            min: min,
            max: max
        });
    },

    apiParams: function() {
        var params = _.clone(this.props.params);
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
                <svg ref="svg" className="barchart"
                    width="100%" height={this.height()}
                    viewBox={"0 0 "+ this.state.canvasWidth + " "+ this.height()}>
                    {this.state.points.map(this.renderBar)}
                </svg>
            </div>
        );
    },

    renderBar: function(point, i) {
        var maxWidth = this.state.canvasWidth,
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
            labelPaddingH = 5, // Horizontal
            labelPaddingV = 2, // Vertical
            labelWidth = textWidth(label),
            labelHeight = 16,
            labelOuterWidth = labelWidth + 2*labelPaddingH,
            labelOffsetWidth = labelOuterWidth + 2*labelPaddingH,
            labelOuterHeight = labelHeight + 2*labelPaddingV,
            labelMarginV = (this.props.height - labelOuterHeight)/2,
            labelX = 0,
            labelY = this.props.y + labelOuterHeight + 1, // 1 is magic
            barX = this.props.x;

        if (labelOffsetWidth <= width) {
            if (offset > 0) {
                if (barX === offset) {
                    labelX = barX + 2*labelPaddingH;
                } else {
                    labelX = barX + width - labelOffsetWidth + 2*labelPaddingH;
                }
            } else {
                labelX = barX + 2*labelPaddingH;
            }
        } else {
            if (barX === offset) {
                labelX = barX + width + 2*labelPaddingH;
            } else if (labelOffsetWidth <= barX) {
                labelX = barX - labelOffsetWidth + 2*labelPaddingH;
            } else {
                labelX = barX + width + labelPaddingH;
            }
        }

        return (
            <g onClick={this.props.onClick}>
                <rect className="bar" fill={this.props.color}
                    width={width} height={this.props.height}
                    x={this.props.x} y={this.props.y} rx="2" ry="2" />
                <rect className="label_underlay"
                    x={labelX - labelPaddingH} y={this.props.y + labelMarginV}
                    height={labelOuterHeight} width={labelOuterWidth}
                    rx="3" ry="3" />
                <text className="label" x={labelX} y={labelY}>{label}</text>
            </g>
        );
    }
});
