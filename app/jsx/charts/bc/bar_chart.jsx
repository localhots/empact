var BarChart = React.createClass({
    mixins: [ReactRouter.Navigation, ReactRouter.State, SVGChartMixin, ChartDataMixin],

    sorts: ['commits', 'delta'],
    numElements: 15,
    barHeight: 30,
    barMargin: 5,

    words: {
        values: { // Sort
            commits: "the number of commits",
            delta: "the delta between lines added and lines removed"
        },
        actions: { // Item-Sort
            "repo-commits": "made to",
            "repo-delta": "for",
            "team-commits": "made by the most active",
            "team-delta": "for the most active",
            "user-commits": "made by the most active",
            "user-delta": "for the most active"
        },
        items: { // Item
            repo: "repositories",
            team: "teams",
            user: "users"
        },
        whatHappened: { // Item-Target
            "user-repo": "working on",
            "team-repo": "working on",
            "team-org": "working on repositories of",
            "user-org": "working on repositories of",
            "repo-org": "that were most actively modified by the members of",
            "user-team": "working on repositories of the",
            "repo-team": "that were most actively modified by the members of the",
            "repo-user": "that were most actively modified by"
        },
        targetSuffix: { // Subject of current context
            repo: "repository",
            team: "team"
        },
    },

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
        var words = this.words,
            target = (this.getParams().repo ? 'repo'
                : this.getParams().team ? 'team'
                : this.getParams().user ? 'user'
                : 'org');
            subject = this.getParams()[target];

        var points = _.clone(this.state.points);
        if (points.length > 0) {
            for (var i = points.length; i < this.numElements; i++) {
                points.push({commits: 0, delta: 0});
            }
        }

        return (
            <div className="barchart-container">
                <div className="whatsgoingon">
                    This bar chart shows <em>{words.values[this.state.sort]}</em> {words.actions[this.state.item +'-'+ this.state.sort]} <em>{words.items[this.state.item]}</em> {words.whatHappened[this.state.item +'-'+ target]} <em>{subject}</em> {words.targetSuffix[target]} <WeekIntervalSelector />
                </div>
                <div className="filters">
                    <Selector thing="sort"
                        title="Show"
                        items={this.sorts}
                        value={this.state.sort}
                        onChange={this.handleFilter.bind(this, 'sort')} />
                    <Selector thing="item"
                        title="Grouped by"
                        items={this.props.items}
                        value={this.state.item}
                        onChange={this.handleFilter.bind(this, 'item')} />
                </div>
                <svg ref="svg" className="barchart" key="barchart-svg"
                    width="100%" height={this.height()}
                    viewBox={"0 0 "+ (this.state.canvasWidth || 0) + " "+ this.height()}>
                    {points.map(this.renderBar)}
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
                color={Colors[i]}
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
