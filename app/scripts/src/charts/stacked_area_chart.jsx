var StackedAreaChart = React.createClass({
    mixins: [ReactRouter.Navigation, ReactRouter.State, SVGChartMixin, ChartDataMixin],

    numElements: 10,
    maxWeeks: 30,
    height: 350,

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
        if (this.props.items[i] !== this.state.item) {
            this.setState({
                item: this.props.items[i],
                state: 'newProps'
            }, this.fetchData);
        }
    },

    handleClick: function(point) {
        var params = {org: this.getParams().org};
        params[this.state.item] = point.item;
        this.transitionTo(this.state.item, params);
    },

    handleNewData: function() {
        // Group commits by items
        var weeksList = _.chain(this.state.rawData)
            .pluck('week')
            .uniq()
            .sort()
            .reverse()
            .take(this.maxWeeks)
            .value();

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
        var top = _.chain(_.pairs(counts)) // Take [item, count] pairs from counts object
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

        var max = _.max(_.map(weeksList, function(week) {
                return _.sum(_.values(weeks[week]));
            }));

        this.setState({
            top: top,
            weeks: weeks,
            max: max,
            state: 'newPoints'
        });
    },

    buildPathD: function(points) {
        var maxWidth = this.state.canvasWidth,
            maxHeight = this.height,
            maxValue = this.state.max,
            len = points.length;

        var d = _.map(points, function(point, i) {
                return 'L'+ Math.floor(i/(len-1)*maxWidth) +','+ Math.floor(maxHeight - point);
            });
        d.unshift('M0,'+ maxHeight);
        d.push('L'+ maxWidth +','+ maxHeight +'Z');

        return d.join(' ');
    },

    render: function() {
        var maxWidth = this.state.canvasWidth,
            maxHeight = this.height,
            top = this.state.top,
            max = this.state.max;

        var points = _.chain(this.state.weeks)
            .map(function(items, week) {
                var values = _.map(top, function(item) {
                    return items[item] || 0;
                });

                var sum = 0;
                var points = _.map(values, function(val) {
                    sum += Math.floor(val/max*maxHeight);
                    return sum;
                });

                return [week, points];
            })
            .sort(0)
            .reverse()
            .take(this.maxWeeks)
            .reverse()
            .value();

        var paths = _.map(top, function(item, i) {
            var itemPoints = _.map(points, function(pair) {
                return pair[1][i];
            });
            return[item, itemPoints];
        });

        var colors = {};
        var areas = _.map(paths, function(pair, i) {
            var item = pair[0], path = pair[1];
            if (item !== null) {
                colors[item] = Colors[i];
            }
            return (
                <StackedArea key={'area-'+ i}
                    item={item}
                    d={roundPathCorners(this.buildPathD(path), 3)}
                    color={Colors[i]} />
            );
        }.bind(this));

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
            <div className="sachart-container">
                <div className="whatsgoingon">
                    This stacked area chart represents <em>{words.items[this.state.item]}</em> {words.actions[this.state.item]} <em>{who}</em> {words.item[subject]} from <em>W11, Mar 9</em> to <em>W18, Apr 27</em>
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
                    width="100%" height={maxHeight}
                    viewBox={"0 0 "+ (this.state.canvasWidth || 0) + " "+ maxHeight}>
                    {areas.reverse()}
                </svg>
                <ul className="legend">
                    {_.pairs(colors).map(function(pair){
                        return (
                            <li key={'legend-'+ pair[0]}>
                                <div className="color-dot" style={{backgroundColor: pair[1]}}></div>
                                {pair[0]}
                            </li>
                        );
                    })}
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
                d={this.state.lastd || this.props.d}
                fill={this.props.color}
                shapeRendering="optimizeQuality" />
        );
    }
});
