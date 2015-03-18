var Dot = React.createClass({
    mixins: [ChartAnimationMixin],

    radius: 10,

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
        return (
            <g className={'dot dot-'+ this.props.i} onMouseOver={this.props.onMouseOver}>
                <circle ref="dot"
                    cx={this.props.x}
                    cy={this.state.lastY || this.props.y}
                    r={this.radius} />
                <text ref="value"
                    x={this.props.x}
                    y={this.props.y+4}>
                    {this.props.value}
                </text>
            </g>
        );
    }
});
