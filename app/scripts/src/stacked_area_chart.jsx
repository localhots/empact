var StackedAreaChart = React.createClass({
    mixins: [Router.Navigation, Router.State],

    numElements: 10,
    height: 250,

    getInitialState: function() {
        return {
            item: this.props.items[0],
            rawData: [],
            top: [],
            max: 1,
            weeks: [],
            canvasWidth: 500
        };
    },

    componentDidMount: function() {
        this.fetchData();
        this.calculateViewBoxWidth();
        window.addEventListener('resize', this.calculateViewBoxWidth);
    },

    componentWillReceiveProps: function(newProps) {
        this.setState({
            'item': newProps.items[0],
            'sort': 'commits'
        }, this.fetchData);
    },

    calculateViewBoxWidth: function() {
        this.setState({
           canvasWidth: this.refs.svg.getDOMNode().offsetWidth
        });
    },

    handleFilter: function(thing, i) {
        if (this.props.items[i] !== this.state.item) {
            this.setState({
                item: this.props.items[i]
            }, this.fetchData);
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
            }, this.buildPoints);
        }.bind(this));
    },

    apiParams: function() {
        var params = _.clone(this.props.params);
        params['item'] = this.state.item;
        return params;
    },

    buildPoints: function() {
        // Group commits by items
        var counts = _.reduce(this.state.rawData, function(res, el) {
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

        var weeks = _.reduce(this.state.rawData, function(res, el) {
            if (res[el.week] === undefined) {
                res[el.week] = {};
            }
            if (top.indexOf(el.item) > -1) {
                res[el.week][el.item] = el.commits;
            }
            return res;
        }, {});
        var max = _.chain(weeks).keys().sort().reverse().take(15).map(function(week) {
            return _.sum(_.values(weeks[week]));
        })
        .max()
        .value();

        // var max = _.max(_.map(weeks, function(items, week) {
        //         return _.sum(_.values(items));
        //     }));

        this.setState({
            top: top,
            max: max,
            weeks: weeks
        });
    },

    buildPathD: function(points) {
        var maxWidth = this.state.canvasWidth,
            maxHeight = this.height,
            maxValue = this.state.max,
            len = points.length;
        var d = _.map(points, function(point, i) {
            return 'L'+ Math.floor(i/len*maxWidth) +','+ Math.floor(maxHeight - point);
        });
        d.unshift('M0,'+ maxHeight);
        d.push('L'+ maxWidth +','+ maxHeight +'Z');

        // for (var i = 0; i < missing; i++) {
        //     d.push('L'+ i +','+ this.props.height/2);
        // }
        // for (var i = 0; i < points.length; i++) {
        //     d.push('L'+ missing+i +','+ points[i]);
        // }
        // d.push('L'+ this.props.width +','+ this.props.height/2, 'Z');

        return d.join(' ');
    },

    render: function() {
        var maxWidth = this.state.canvasWidth,
            maxHeight = this.height,
            rtop = this.state.top.reverse(),
            max = this.state.max;

        var points = _.chain(this.state.weeks)
            .map(function(items, week) {
                var values = _.map(rtop, function(item) {
                    return items[item] || 0;
                });

                var sum = 0;
                // console.log('----------');
                var points = _.map(values, function(val) {
                    sum += Math.floor(val/max*maxHeight);
                    // console.log(val, max, maxHeight, sum);
                    return sum;
                });

                return [week, points];
            })
            .sort(0)
            .value();

        var paths = _.reduce(rtop, function(res, item, i) {
            res[item] = _.map(points, function(pair) {
                return pair[1][i];
            }).slice(-15);
            return res;
        }, {});

        var i = -1;
        var colors = {}
        var areas = _.map(paths, function(path, item) {
            i++;
            colors[item] = Colors2[i];
            return (
                <StackedArea key={'sa-item-'+ item}
                    item={item}
                    path={roundPathCorners(this.buildPathD(path), 5)}
                    color={Colors2[i]} />
            );
        }.bind(this));

        return (
            <div className="sachart-container">
                <div className="filters">
                    <Selector thing="item"
                        items={this.props.items}
                        value={this.state.item}
                        onChange={this.handleFilter.bind(this, 'item')} />
                    <Selector thing="sort"
                        items={['commits']}
                        value={'commits'} />
                </div>
                <svg ref="svg" className="sachart"
                    width="100%" height={maxHeight}
                    viewBox={"0 0 "+ this.state.canvasWidth + " "+ maxHeight}>
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
    render: function() {
        return (
            <path key={'sac-area-'+ this.props.item}
                d={this.props.path}
                fill={this.props.color}
                shapeRendering="optimizeQuality" />
        );
    }
});
