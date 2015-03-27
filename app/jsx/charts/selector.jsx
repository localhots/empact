var Selector = React.createClass({
    names: {
        "repo":    "Repositories",
        "team":    "Teams",
        "user":    "Users",
        "commits": "Commits",
        "delta":   "Delta"
    },

    renderItem: function(item, i) {
        var itemClass = (item === this.props.value ? 'active state-button' : 'state-button'),
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
            <div className="selector">
                <div className="title">{this.props.title}</div>
                <ul className="items">
                    {this.props.items.map(this.renderItem)}
                </ul>
            </div>
        );
    }
});
