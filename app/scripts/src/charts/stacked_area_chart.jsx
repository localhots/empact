var StackedAreaChart = React.createClass({
    mixins: [ReactRouter.Navigation, ReactRouter.State, SVGChartMixin, ChartDataMixin],

    numElements: 10,
    maxWeeks: 30,
    height: 350,
    xAxisHeight: 20,

    words: {
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

    getInitialState: function() {
        return {
            item: this.props.items[0],
            rawData: [],
            top: [],
            weeks: [],
            max: 1
        };
    },

    componentDidMount: function() {
        this.calculateViewBoxWidth();
        window.addEventListener('resize', this.calculateViewBoxWidth);
    },

    componentWillReceiveProps: function(newProps) {
        this.setState({
            item: (_.isEqual(newProps.items, this.props.items)
                ? this.state.item
                : newProps.items[0]),
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
        if (this.props.items[i] !== this.state.item) {
            this.setState({
                item: this.props.items[i],
                state: 'newProps'
            }, this.fetchData);
        }
    },

    handleClick: function(item) {
        var params = {org: this.getParams().org};
        params[this.state.item] = item;
        this.transitionTo(this.state.item, params, this.getQuery());
    },

    handleFocusIn: function(i) {
        var node = this.refs.container.getDOMNode();
        node.className = 'sachart-container focused item-'+ i;
    },

    handleFocusOut: function() {
        var node = this.refs.container.getDOMNode();
        node.className = 'sachart-container';
    },

    handleNewData: function() {
        // Group commits by items
        var weeksList = _(this.state.rawData).pluck('week').uniq().sort().reverse().take(this.maxWeeks).value();
        if (weeksList.length < 2) {
            this.setState({
                weeks: [],
                state: 'newPoints'
            });
            return;
        }

        var counts = _.reduce(this.state.rawData, function(res, el) {
            if (weeksList.indexOf(el.week) === -1) {
                return res;
            }
            if (res[el.item] === undefined) {
                res[el.item] = el.commits;
            } else {
                res[el.item] += el.commits;
            }
            return res;
        }, {});

        // Extract top items from
        var top = _(_.pairs(counts)) // Take [item, count] pairs from counts object
            .sortBy(1).reverse() // sort them by count (descending)
            .take(this.numElements) // take first N pairs
            .pluck(0) // keep only items, omit the counts
            .value();
        for (var i = top.length; i < this.numElements; i++) {
            top[i] = null;
        };

        var weeks = _.reduce(this.state.rawData, function(res, el) {
            if (weeksList.indexOf(el.week) === -1) {
                return res;
            }
            if (res[el.week] === undefined) {
                res[el.week] = {};
            }
            if (top.indexOf(el.item) > -1) {
                res[el.week][el.item] = el.commits;
            }
            return res;
        }, {});

        var max = _.max(_.map(weeksList, function(week){ return _.sum(_.values(weeks[week])); }));

        this.setState({
            top: top,
            weeks: weeks,
            max: max,
            state: 'newPoints'
        });
    },

    buildPathD: function(points) {
        var maxWidth = this.state.canvasWidth,
            maxHeight = this.height;

        var dots = this.buildDots(points);
        var first = dots.shift();
        var d = _.map(dots, function(dot){ return 'L'+ dot.x +','+ dot.y; });
        d.unshift('M'+ first.x +','+ first.y);
        d.push('L'+ maxWidth +','+ maxHeight);
        d.push('L0,'+ maxHeight +' Z');

        return d.join(' ');
    },

    buildDots: function(points) {
        var maxWidth = this.state.canvasWidth,
            maxHeight = this.height,
            maxValue = this.state.max,
            len = points.length;

        return _.map(points, function(point, i) {
            point.x = i/(len-1)*maxWidth;
            point.y = maxHeight - point.point;
            return point;
        });
    },

    render: function() {
        var maxWidth = this.state.canvasWidth,
            maxHeight = this.height,
            top = this.state.top,
            max = this.state.max;

        // [week, [{val, point}, ...]]
        var points = _(this.state.weeks)
            .map(function(items, week) {
                var values = _.map(top, function(item) {
                    return items[item] || 0;
                });

                var sum = 0;
                var points = _.map(values, function(val) {
                    sum += val/max*maxHeight*0.96;
                    return {
                        val: val,
                        point: sum
                    };
                });

                return [parseInt(week, 10), points];
            })
            .sort(0)
            .reverse()
            .take(this.maxWeeks)
            .reverse()
            .value();

        // [item, [{val, point}, ...]]
        var paths = _.map(top, function(item, i) {
            var itemPoints = _.map(points, function(pair) {
                return pair[1][i];
            });
            return[item, itemPoints];
        });

        var areas = _.map(paths, function(pair, i) {
            var item = pair[0],
                path = pair[1];

            return (
                <StackedArea key={'area-'+ i}
                    item={item} i={i}
                    d={this.buildPathD(path)}
                    color={Colors[i]}
                    onMouseOver={this.handleFocusIn.bind(this, i)} />
            );
        }.bind(this));

        var words = this.words,
            who = this.getParams().repo ||
                  this.getParams().team ||
                  this.getParams().user ||
                  this.getParams().org;

        var params = Object.keys(this.getParams());
        params.splice(params.indexOf('org'), 1);
        var subject = params[0];

        var renderDot = function(item, i, dot, j) {
            if (dot.val === 0) {
                return null;
            }

            var maxWidth = this.state.canvasWidth,
                maxHeight = this.height,
                radius = 10,
                x = dot.x,
                y = dot.y;

            if (x < radius) {
                x = radius
            } else if (x > maxWidth - radius) {
                x = maxWidth - radius;
            }
            if (y < radius) {
                y = radius;
            } else if (y > maxHeight - radius) {
                y = maxHeight - radius;
            }

            return (
                <Dot key={'dot-'+ i +'-'+ j}
                    item={item} i={i}
                    value={dot.val}
                    x={x}
                    y={y}
                    onMouseOver={this.handleFocusIn.bind(this, i)} />
            );
        }.bind(this);

        var renderedDots = _.map(paths, function(pair, i) {
            var item = pair[0], path = pair[1];
            var dots = this.buildDots(path);
            return dots.map(renderDot.bind(this, item, i));
        }.bind(this));

        var renderLegend = function(item, i){
            return (
                <li key={'legend-'+ item}
                    className={'label label-'+ i}
                    onMouseOver={this.handleFocusIn.bind(this, i)}
                    onMouseOut={this.handleFocusOut.bind(this, i)}
                    onClick={this.handleClick.bind(this, item)}
                    >
                    <div className="color-dot" style={{backgroundColor: Colors[i]}}></div>
                    {item}
                </li>
            );
        }.bind(this);

        var legend = _(paths).pluck(0).filter(function(el){ return el !== null; }).value();

        return (
            <div ref="container" className="sachart-container">
                <div className="whatsgoingon">
                    This stacked area chart represents <em>{words.items[this.state.item]}</em> {words.actions[this.state.item]} <em>{who}</em> {words.item[subject]} <WeekIntervalSelector />
                </div>
                <div className="filters">
                    <Selector thing="sort"
                        title="Show"
                        items={['commits']}
                        value={'commits'} />
                    <Selector thing="item"
                        title="Grouped by"
                        items={this.props.items}
                        value={this.state.item}
                        onChange={this.handleFilter.bind(this, 'item')} />
                </div>
                <svg ref="svg" className="sachart" key="sachart-svg"
                    width="100%"
                    height={this.height + this.xAxisHeight}
                    viewBox={"0 0 "+ (this.state.canvasWidth || 0) + " "+ (this.height + this.xAxisHeight)}
                    onMouseOut={this.handleFocusOut}
                    >
                    <g ref="areas">{areas.reverse()}</g>
                    <g ref="dots">{renderedDots}</g>
                    <Axis
                        weeks={_.pluck(points, 0)}
                        y={this.height + 3}
                        width={this.state.canvasWidth} />
                </svg>
                <ul className="legend">
                    {legend.map(renderLegend)}
                </ul>
            </div>
        );
    }
});

var StackedArea = React.createClass({
    mixins: [ChartAnimationMixin],

    getInitialState: function() {
        return {};
    },

    componentWillReceiveProps: function(newProps) {
        this.setState({
            lastd: this.props.d || newProps.d,
        }, this.animateAll);
    },

    animateAll: function() {
        this.clearAnimations(this.refs.path);
        this.animate(this.refs.path, 'd', this.state.lastd, this.props.d);
    },

    render: function() {
        return (
            <path ref="path"
                className={'path path-'+ this.props.i}
                d={this.state.lastd || this.props.d}
                fill={this.props.color}
                onMouseOver={this.props.onMouseOver}
                shapeRendering="optimizeQuality" />
        );
    }
});

var Dot = React.createClass({
    mixins: [ChartAnimationMixin],

    radius: 10,

    getInitialState: function() {
        return {};
    },

    componentWillReceiveProps: function(newProps) {
        this.setState({
            lastY: this.props.y || newProps.y
        }, this.animateAll);
    },

    animateAll: function() {
        this.clearAnimations(this.refs.dot);
        this.animate(this.refs.dot, 'cy', this.state.lastY, this.props.y);
    },

    render: function() {
        return (
            <g className={'dot dot-'+ this.props.i} onMouseOver={this.props.onMouseOver}>
                <circle ref="dot"
                    cx={this.props.x}
                    cy={this.state.lastY || this.props.y}
                    r={this.radius}
                    fill="#fff"
                    stroke="#f0f0f0"
                    strokeWidth="2" />
                <text ref="value"
                    x={this.props.x}
                    y={this.props.y+4}
                    textAnchor="middle">
                    {this.props.value}
                </text>
            </g>
        );
    }
});

var Axis = React.createClass({
    render: function() {
        if (this.props.weeks.length === 0) {
            return null;
        }
        var renderMark = function(week, i) {
            var len = this.props.weeks.length,
                x = i/(len - 1)*this.props.width,
                showLabel,
                ta = (i === 0
                    ? 'start'
                    : (i === len - 1
                        ? 'end'
                        : 'middle'));

            // Thin out labels
            if (len > 20) {
                showLabel = (i % 2 === 0);
            } else if (len > 10) {
                showLabel = (i % 2 === 0);
            } else {
                showLabel = true;
            }

            return (
                <g key={'mark-'+ i}>
                    <line
                        x1={x}
                        y1={this.props.y}
                        x2={x}
                        y2={this.props.y + 4}
                        stroke="#666" strokeWidth="1" />
                    {showLabel ? <text className="axis-mark"
                        x={x}
                        y={this.props.y + 15}
                        textAnchor={ta}
                        fill="#666"
                        >
                        {formatDate(week)}
                    </text> : null}
                </g>
            );
        }.bind(this);

        return (
            <g ref="axis">
                <line
                    x1="0"
                    y1={this.props.y}
                    x2={this.props.width}
                    y2={this.props.y}
                    stroke="#666" strokeWidth="1" />
                {this.props.weeks.map(renderMark)}
                <line
                    x1={this.props.width - 1}
                    y1={this.props.y}
                    x2={this.props.width - 1}
                    y2={this.props.y + 4}
                    stroke="#666" strokeWidth="1" />
            </g>
        )
    }
});
