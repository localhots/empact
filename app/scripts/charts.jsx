var Router = ReactRouter;

var BarChart = React.createClass({
    barHeight: 40,
    barMargin: 5,

    getInitialState: function() {
        return {points: [], max: 1};
    },

    componentDidMount: function() {
        $.get(this.props.api, function(res){
            var max = 1;
            res.map(function(el) {
                if (el.value > max) {
                    max = el.value
                }
            });
            this.setState({points: res, max: max});
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
        return (
            <svg className="barchart" width="100%" height={this.height()}>
                {this.state.points.map(this.renderBar)}
            </svg>
        );
    },

    renderBar: function(point, i) {
        return (
            <Bar key={point.item} point={point} i={i} link={this.props.link}
                y={this.y(i)}
                width={point.value/this.state.max}
                height={this.barHeight} />
        );
    }
});

var Bar = React.createClass({
    mixins: [Router.Navigation],
    handleClick: function(e) {
        this.transitionTo(this.props.link + this.props.point.item);
    },

    render: function() {
        var p = this.props.point
            w = this.props.width*500,
            label = p.item + ": " + p.value,
            tx = 10;
        if (label.length*15 > w) {
            tx = w + tx;
        }
        return (
            <g onClick={this.handleClick}>
                <rect className="bar" fill={Colors2[this.props.i]}
                    width={this.props.width*500}
                    height={this.props.height}
                    x="0" y={this.props.y} rx="2" ry="2" />
                <rect className="label_underlay" x={tx-6} y={this.props.y+10} height={20} width={label.length*10+5} rx="3" ry="3" />
                <text className="label" x={tx} y={this.props.y + 26}>{label}</text>
            </g>
        );
    }
});
