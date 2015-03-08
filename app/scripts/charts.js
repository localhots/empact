var BarChart = React.createClass({displayName: "BarChart",
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
                React.createElement("g", {key: point.item}, 
                    React.createElement("rect", {
                        fill: colorFor(point.item), 
                        width: point.value, 
                        height: this.barHeight, 
                        x: "0", 
                        y: this.y(i)}), 
                    React.createElement("text", {x: "20", y: this.y(i) + 25}, point.item + ": " + point.value)
                )
            );
        }.bind(this)
        return (
            React.createElement("svg", {className: "bar_chart", width: "100%", height: this.height()}, 
                this.state.points.map(renderPoint)
            )
        );
    }
});
