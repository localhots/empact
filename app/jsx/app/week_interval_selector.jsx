var WeekIntervalSelector = React.createClass({
    mixins: [ReactRouter.Navigation, ReactRouter.State],

    getInitialState: function() {
        var ms = 1000,
            daySeconds = 86400,
            weekSeconds = daySeconds*7,
            today = new Date(),
            sunday = new Date(today - daySeconds*ms*today.getDay()),
            perfectSunday = new Date(Date.UTC(sunday.getFullYear(), sunday.getMonth(), sunday.getDate())),
            lastWeek = perfectSunday.setHours(0)/ms,
            firstWeek = lastWeek - 51*weekSeconds;

        var weeks = [];
        for (var i = lastWeek; i >= firstWeek; i -= weekSeconds) {
            weeks.push(i);
        };

        return {
            weeks: weeks.sort()
        };
    },

    handleChange: function(thing, e) {
        var params = this.getQuery();
        params[thing.slice(0, 1)] = e.target.value/100;
        this.transitionTo(document.location.pathname, null, params);
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
