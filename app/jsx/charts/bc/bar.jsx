var Bar = React.createClass({
    mixins: [ReactRouter.Navigation, ChartAnimationMixin],

    height: 28,
    labelPaddingH: 5,     // Label horizontal padding
    labelPaddingV: 2,     // Label vertical padding
    labelMarginV: 4,      // Tune text vertical alignment
    labelHeight: 16,      // Text size
    labelOuterHeight: 20, // labelHeight + 2*labelPaddingV,

    getInitialState: function() {
        return {
            labelX: 0,
            lastLabelX: this.labelPaddingH
        };
    },

    componentDidMount: function() {
        this.calculateLabelPosition();
    },

    componentWillReceiveProps: function(newProps) {
        if (_.isEqual(this.props, newProps)) {
            return;
        }

        this.setState({
            lastBarX: (this.props.x !== undefined ? this.props.x : newProps.x),
            lastBarWidth: (this.props.width !== undefined ? this.props.width : newProps.width),
            lastLabelX: this.state.labelX
        }, this.calculateLabelPosition);
    },

    calculateLabelPosition: function() {
        var val = this.props.value,
            offset = this.props.offset,
            label = this.props.item + ': ' + formatNumber(val),
            labelWidth = textWidth(label),
            labelOffsetWidth = labelWidth + 2*this.labelPaddingH,
            labelX;

        if (offset === 0) {
            labelX = this.labelPaddingH;
        } else {
            if (val < 0) {
                if (offset >= labelOffsetWidth) {
                    labelX = offset - labelOffsetWidth + this.labelPaddingH;
                } else {
                    labelX = offset + this.labelPaddingH;
                }
            } else {
                if (offset + labelOffsetWidth <= this.props.max) {
                    labelX = offset + this.labelPaddingH;
                } else {
                    labelX = offset - labelOffsetWidth + this.labelPaddingH;
                }
            }
        }

        this.setState({
            labelX: labelX,
            barWidth: (this.props.item && this.props.width < 5 ? 5 : this.props.width)
        }, this.animateAll);
    },

    animateAll: function() {
        this.clearAnimations(this.refs.bar);
        this.clearAnimations(this.refs.label);
        this.animate(this.refs.bar, 'width', this.state.lastBarWidth, this.state.barWidth);
        this.animate(this.refs.bar, 'x', this.state.lastBarX, this.props.x);
        this.animate(this.refs.label, 'x', this.state.lastLabelX, this.state.labelX);
    },

    render: function() {
        var label = this.props.item ? (this.props.item + ': ' + formatNumber(this.props.value)) : '',
            labelWidth = textWidth(label),
            labelOuterWidth = (labelWidth == 0 ? 0 : labelWidth + 2*this.labelPaddingH),
            barX = (this.state.lastBarX && this.state.lastBarX !== this.props.x
                ? this.state.lastBarX
                : this.props.x),
            barWidth = (this.state.lastBarWidth && this.state.lastBarWidth !== this.state.barWidth
                ? this.state.lastBarWidth
                : this.state.barWidth);

        return (
            <g onClick={this.props.onClick}>
                <rect ref="bar" className="bar"
                    fill={this.props.color}
                    width={barWidth}
                    height={this.props.height}
                    x={barX}
                    y={this.props.y}
                    rx="2"
                    ry="2" />
                <text ref="label" className="label"
                    x={this.state.labelX}
                    y={this.props.y + this.labelMarginV + this.labelHeight}>
                    {label}
                </text>
            </g>
        );
    }
});
