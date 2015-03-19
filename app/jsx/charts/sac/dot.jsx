var Dot = React.createClass({
    mixins: [ChartAnimationMixin],

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
        var radius,
            val = ''+ this.props.value;
        if (val.length === 1) {
            radius = 8;
        } else if (val.length === 2) {
            radius = 9;
        } else {
            radius = 11;
        }
        return (
            <g className={'dot dot-'+ this.props.i} onMouseOver={this.props.onMouseOver}>
                <circle ref="dot"
                    cx={this.props.x}
                    cy={this.state.lastY || this.props.y}
                    r={radius} />
                <text ref="value"
                    x={this.props.x}
                    y={this.props.y+4}>
                    {this.props.value}
                </text>
            </g>
        );
    }
});
