var Axis = React.createClass({
    topMargin: 2,
    markHeight: 5,

    render: function() {
        if (this.props.weeks.length === 0) {
            return null;
        }
        var renderMark = function(week, i) {
            var len = this.props.weeks.length,
                x = i/(len - 1)*this.props.width + 1,
                showLabel,
                ta = (i === 0 // Text anchor for the leftmost label
                    ? 'start'
                    : (i === len - 1 // Text anchor for the rightmost label
                        ? 'end'
                        : 'middle')); // Text anchor for other labels

            // Thin out labels
            if (len > 20) {
                showLabel = (i % 3 === 0);
            } else if (len > 10) {
                showLabel = (i % 2 === 0);
            } else {
                showLabel = true;
            }

            return (
                <g key={'mark-'+ i}>
                    <line className="axis"
                        x1={x}
                        y1={this.props.y + this.topMargin}
                        x2={x}
                        y2={this.props.y + this.topMargin + this.markHeight} />
                    {!showLabel ? null : <text className="axis-mark"
                        x={x}
                        y={this.props.y + this.topMargin + 14}
                        textAnchor={ta}
                        >
                        {formatDate(week)}
                    </text>}
                </g>
            );
        }.bind(this);

        return (
            <g ref="axis">
                <rect // This rect hides area bouncing glitches
                    x="0"
                    y={this.props.y}
                    width={this.props.width}
                    height={this.props.height}
                    fill="#fff" />
                <line className="axis"
                    x1="0"
                    y1={this.props.y + this.topMargin}
                    x2={this.props.width}
                    y2={this.props.y + this.topMargin} />
                {this.props.weeks.map(renderMark)}
                <line className="axis"
                    x1={this.props.width - 1}
                    y1={this.props.y + this.topMargin}
                    x2={this.props.width - 1}
                    y2={this.props.y + this.topMargin + this.markHeight} />
            </g>
        )
    }
});
