var Selector = React.createClass({
    names: {
        "repo":    "Repositories",
        "team":    "Teams",
        "user":    "Users",
        "commits": "Commits",
        "delta":   "Delta"
    },

    renderItem: function(item, i) {
        var itemClass = (item === this.props.value ? 'active' : ''),
            clickEvent = null;
        if (this.props.onChange) {
            clickEvent = this.props.onChange.bind(this, i);
        }
        return (
            <li key={item} onClick={clickEvent} className={itemClass}>{this.names[item]}</li>
        );
    },

    render: function() {
        return (
            <ul className={this.props.thing}>
                {this.props.items.map(this.renderItem)}
            </ul>
        );
    }
});
