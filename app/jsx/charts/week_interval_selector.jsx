var WeekIntervalSelector = React.createClass({
    mixins: [ReactRouter.Navigation, ReactRouter.State],

    getInitialState: function() {
        return {
            weeks: []
        };
    },

    componentDidMount: function() {
        this.loadWeekRange(this.props.org);
    },

    componentWillReceiveProps: function(newProps) {
        this.loadWeekRange(newProps.org);
    },

    handleChange: function(thing, e) {
        var params = this.getQuery();
        params[thing.slice(0, 1)] = e.target.value/100;
        this.transitionTo(document.location.pathname, null, params);
    },

    loadWeekRange: function(org) {
        getURL("/api/weeks", {org: org}, function(res){
            var weeks = [],
                min = res[0],
                max = res[1];
            for (var i = min; i <= max; i += 86400*7) {
                weeks.push(i);
            };
            this.setState({
                weeks: weeks
            });
        }.bind(this));
    },

    render: function() {
        var daySeconds = 86400,
            weekSeconds = daySeconds*7,
            lastWeek = this.state.weeks[this.state.weeks.length-1],
            from = (this.getQuery().f ? parseInt(this.getQuery().f, 10)*100 : lastWeek - 29*weekSeconds),
            to = (this.getQuery().t ? parseInt(this.getQuery().t, 10)*100 : lastWeek);

        var weeksBefore = _(this.state.weeks)
            .filter(function(week) {
                return week < to;
            })
            .reverse()
            .value();
        var weeksAfter = _(this.state.weeks)
            .filter(function(week) {
                return week > from;
            })
            .reverse()
            .value();

        var renderOption = function(ts) {
            return (
                <option key={ts} value={ts}>{formatDate(ts, true)}</option>
            );
        };

        return (
            <div className="week-selector">
                <span>from</span>
                <div ref="from" className="selector">
                    <em ref="label">{formatDate(from)}</em>
                    <select ref="select" value={from} onChange={this.handleChange.bind(this, 'from')}>
                        {weeksBefore.map(renderOption)}
                    </select>
                </div>
                <span>to</span>
                <div ref="to" className="selector">
                    <em ref="label">{formatDate(to)}</em>
                    <select ref="select" value={to} onChange={this.handleChange.bind(this, 'to')}>
                        {weeksAfter.map(renderOption)}
                    </select>
                </div>
            </div>
        );
    }
});
