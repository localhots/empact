var SVGNS = "http://www.w3.org/2000/svg",
    Router = ReactRouter;

var BarChart = React.createClass({
    barHeight: 40,
    barMargin: 5,

    getInitialState: function() {
        return {points: [], max: 1};
    },

    componentDidMount: function() {
        $.get(this.props.api, function(res){
            res = res.slice(0, 15);
            var max = 1;
            res.map(function(el) {
                if (el.commits > max) {
                    max = el.commits
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
                width={point.commits/this.state.max}
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
            label = p.item + ': ' + p.commits,
            labelm = 10, // Margin
            labelw = label.length*9 + 2*labelm, // Width
            textx = labelm;
        if (labelw + 2*labelm > w) {
            textx = w + textx;
        }
        return (
            <g onClick={this.handleClick}>
                <rect className="bar" fill={Colors2[this.props.i]}
                    width={this.props.width*500}
                    height={this.props.height}
                    x="0" y={this.props.y} rx="2" ry="2" />
                <rect className="label_underlay"
                    x={textx-6} y={this.props.y+10}
                    height={20} width={labelw}
                    rx="3" ry="3" />
                <text className="label" x={textx} y={this.props.y + 26}>{label}</text>
            </g>
        );
    }
});
