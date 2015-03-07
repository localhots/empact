var Menu = React.createClass({
    getInitialState: function() {
        return {teams: []};
    },

    componentDidMount: function() {
        this.loadTeams();
    },

    loadTeams: function() {
        $.get(this.props.api, function(res){
            this.setState({teams: res})
        }.bind(this));
    },

    render: function() {
        var renderTeam = function(team) {
            return (
                <li className="nav team">{team.name}</li>
            )
        };
        return (
            <ul>
                <li className="nav empact">Empact</li>
                <li className="nav dash">Dashboard</li>
                <li className="nav repos">Repos</li>
                <li className="nav header">Teams:</li>
                {this.state.teams.map(renderTeam)}
            </ul>
        );
    }
});
