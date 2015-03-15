var BarChart = React.createClass({
    mixins: [ReactRouter.Navigation, ReactRouter.State, SVGChartMixin, ChartDataMixin],

    sorts: ['commits', 'delta'],
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
            max: 1
        };
    },

    componentDidMount: function() {
        this.calculateViewBoxWidth();
        window.addEventListener('resize', this.calculateViewBoxWidth);
        this.componentWillReceiveProps(this.props);
    },

    componentWillReceiveProps: function(newProps) {
        this.setState({
            item: newProps.items[0],
            state: 'newProps'
        }, this.fetchData);
    },

    shouldComponentUpdate: function(newProps, newState) {
        if (!newState.canvasWidth) {
            return false;
        }
        if (newState.state !== 'newPoints')  {
            return false;
        }
        return true;
    },

    handleFilter: function(thing, i) {
        if (thing === 'item' && this.props.items[i] !== this.state.item) {
            this.setState({
                item: this.props.items[i],
                state: 'newProps'
            }, this.fetchData);
        } else if (thing === 'sort' && this.sorts[i] !== this.state.sort) {
            this.setState({
                sort: this.sorts[i],
                state: 'newProps'
            }, this.handleNewData);
        }
    },

    handleClick: function(point) {
        var params = {org: this.getParams().org};
        params[this.state.item] = point.item;
        this.transitionTo(this.state.item, params, this.getQuery());
    },

    handleNewData: function() {
        var min = 0, max = 1;
        var points = _.chain(this.state.rawData)
            .sort(function(a, b) {
                return Math.abs(b[this.state.sort]) - Math.abs(a[this.state.sort]);
            }.bind(this))
            .take(this.numElements)
            .value();
        for (var i = points.length; i < this.numElements; i++) {
            var point = {};
            point[this.state.sort] = 0;
            points.push(point);
        }

        this.setState({
            points: points,
            min: _.min(points, this.state.sort)[this.state.sort],
            max: _.max(points, this.state.sort)[this.state.sort],
            state: 'newPoints'
        });
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
        var words = {
                items: {
                    repo: 'repositories',
                    team: 'teams',
                    user: 'contributors'
                },
                item: {
                    repo: 'repository',
                    team: 'team'
                },
                actions: {
                    repo: 'which were the most attended by',
                    team: 'which were the most active working on',
                    user: 'which were the most active working on'
                }
            },
            who = this.getParams().repo || this.getParams().team || this.getParams().user || this.getParams().org;

        var params = Object.keys(this.getParams());
        params.splice(params.indexOf('org'), 1);
        var subject = params[0];

        return (
            <div className="barchart-container">
                <div className="whatsgoingon">
                    This bar chart represents <em>{words.items[this.state.item]}</em> {words.actions[this.state.item]} <em>{who}</em> {words.item[subject]} from <em>W11, Mar 9</em> to <em>W18, Apr 27</em>
                </div>
                <div className="filters">
                    <Selector thing="item"
                        items={this.props.items}
                        value={this.state.item}
                        onChange={this.handleFilter.bind(this, 'item')} />
                    <Selector thing="sort"
                        items={this.sorts}
                        value={this.state.sort}
                        onChange={this.handleFilter.bind(this, 'sort')} />
                </div>
                <svg ref="svg" className="barchart" key="barchart-svg"
                    width="100%" height={this.height()}
                    viewBox={"0 0 "+ (this.state.canvasWidth || 0) + " "+ this.height()}>
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
            width = Math.floor(Math.abs(val)/max2*maxWidth),
            height = this.barHeight,
            offset = (min < 0 ? -min : 0)/max2*maxWidth,
            x = (min >= 0 ? 0 : offset - (val >= 0 ? 0 : width)),
            y = this.y(i);

        return (
            <Bar key={'bar-'+ i}
                item={point.item}
                value={val}
                color={Colors2[i]}
                x={x}
                y={y}
                offset={offset}
                width={width}
                height={height}
                max={maxWidth}
                onClick={this.handleClick.bind(this, point)} />
        );
    }
});

var Bar = React.createClass({
    mixins: [ReactRouter.Navigation, ChartAnimationMixin],

    height: 30,
    labelPaddingH: 5,     // Label horizontal padding
    labelPaddingV: 2,     // Label vertical padding
    labelMarginV: 5,      // Same as padding
    labelHeight: 16,      // Text size
    labelOuterHeight: 20, // labelHeight + 2*labelPaddingV,

    getInitialState: function() {
        return {
            labelX: 0,
            lastLabelX: 2*this.labelPaddingH
        };
    },

    componentDidMount: function() {
        this.calculateLabelPosition();
    },

    componentWillReceiveProps: function(newProps) {
        if (_.isEqual(this.props, newProps)) {
            return;
        }

        this.setState({
            lastBarX: (this.props.x !== undefined ? this.props.x : newProps.x),
            lastBarWidth: (this.props.width !== undefined ? this.props.width : newProps.width),
            lastLabelX: this.state.labelX
        }, this.calculateLabelPosition);
    },

    calculateLabelPosition: function() {
        var val = this.props.value,
            offset = this.props.offset,
            label = this.props.item + ': ' + numberFormat(val),
            labelWidth = textWidth(label),
            labelOuterWidth = labelWidth + 2*this.labelPaddingH,
            labelOffsetWidth = labelOuterWidth + 2*this.labelPaddingH,
            labelX;

        if (offset === 0) {
            labelX = 2*this.labelPaddingH;
        } else {
            if (val < 0) {
                if (offset >= labelOffsetWidth) {
                    labelX = offset - labelOffsetWidth + 2*this.labelPaddingH;
                } else {
                    labelX = offset + 2*this.labelPaddingH;
                }
            } else {
                if (offset + labelOffsetWidth <= this.props.max) {
                    labelX = offset + 2*this.labelPaddingH;
                } else {
                    labelX = offset - labelOffsetWidth + 2*this.labelPaddingH;
                }
            }
        }

        this.setState({
            labelX: labelX,
            barWidth: (this.props.item && this.props.width < 5 ? 5 : this.props.width)
        }, this.animateAll);
    },

    animateAll: function() {
        var bar = this.refs.bar,
            underlay = this.refs.label.refs.underlay,
            text = this.refs.label.refs.text,
            padH = this.labelPaddingH;

        this.clearAnimations(bar);
        this.clearAnimations(underlay);
        this.clearAnimations(text);
        this.animate(bar, 'width', this.state.lastBarWidth, this.state.barWidth);
        this.animate(bar, 'x', this.state.lastBarX, this.props.x);
        this.animate(underlay, 'x', this.state.lastLabelX - padH, this.state.labelX - padH);
        this.animate(text, 'x', this.state.lastLabelX, this.state.labelX);
    },

    render: function() {
        var label = this.props.item ? (this.props.item + ': ' + numberFormat(this.props.value)) : '',
            labelWidth = textWidth(label),
            labelOuterWidth = (labelWidth == 0 ? 0 : labelWidth + 2*this.labelPaddingH),
            barX = (this.state.lastBarX && this.state.lastBarX !== this.props.x
                ? this.state.lastBarX
                : this.props.x),
            barWidth = (this.state.lastBarWidth && this.state.lastBarWidth !== this.state.barWidth
                ? this.state.lastBarWidth
                : this.state.barWidth);

        return (
            <g onClick={this.props.onClick}>
                <rect ref="bar"
                    className="bar"
                    fill={this.props.color}
                    width={barWidth}
                    height={this.props.height}
                    x={barX}
                    y={this.props.y}
                    rx="2"
                    ry="2" />
                <BarLabel ref="label"
                    item={this.props.item}
                    value={this.props.value}
                    width={labelOuterWidth}
                    height={this.labelOuterHeight}
                    x={this.state.labelX}
                    y={this.props.y + this.labelMarginV} />
            </g>
        );
    }
});

var BarLabel = React.createClass({
    render: function() {
        var text = (this.props.item ? this.props.item +': '+ numberFormat(this.props.value) : '');
        return (
            <g>
                <rect ref="underlay" key="underlay"
                    className="label_underlay"
                    width={this.props.width}
                    height={this.props.height}
                    x={this.props.x}
                    y={this.props.y}
                    rx="3"
                    ry="3" />
                <text ref="text" key="text"
                    className="label"
                    x={this.props.x}
                    y={this.props.y + 16}
                    height={20}>
                    {text}
                </text>
            </g>
        );
    }
})
