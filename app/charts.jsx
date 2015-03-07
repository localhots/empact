var BarChart = React.createClass({
    barHeight: 40,
    barMargin: 5,

    getInitialState: function() {
        return {points: []};
    },

    componentDidMount: function() {
        $.get(this.props.url, function(res){
            this.setState({points: res});
        }.bind(this))
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
        var renderPoint = function(point, i) {
            return (
                <g key={point.item}>
                    <rect
                        fill={colorFor(point.item)}
                        width={point.value}
                        height={this.barHeight}
                        x="0"
                        y={this.y(i)} />
                    <text x="20" y={this.y(i) + 25}>{point.item + ": " + point.value}</text>
                </g>
            );
        }.bind(this)
        return (
            <svg className="bar_chart" width="100%" height={this.height()}>
                {this.state.points.map(renderPoint)}
            </svg>
        );
    }
});
