var Area = React.createClass({
    mixins: [ChartAnimationMixin],

    getInitialState: function() {
        return {};
    },

    componentWillReceiveProps: function(newProps) {
        this.setState({
            lastd: this.props.d || newProps.d,
        }, this.animateAll);
    },

    animateAll: function() {
        this.clearAnimations(this.refs.path);
        this.animate(this.refs.path, 'd', this.state.lastd, this.props.d);
    },

    render: function() {
        return (
            <path ref="path"
                className={'path path-'+ this.props.i}
                d={this.state.lastd || this.props.d}
                fill={this.props.color}
                onMouseOver={this.props.onMouseOver}
                shapeRendering="optimizeQuality" />
        );
    }
});
