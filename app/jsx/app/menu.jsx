var Menu = React.createClass({
    mixins: [ReactRouter.State],

    render: function() {
        var renderOrg = function(org) {
            return (
                <li key={'org-'+ org.login} className="nav org">
                    <ReactRouter.Link to="org"
                        params={{org: org.login}}
                        query={this.getQuery()}>
                        {org.login}
                    </ReactRouter.Link>
                </li>
            )
        }.bind(this);
        var renderTeam = function(team) {
            return (
                <li key={'team-'+ team.name} className="nav team">
                    <ReactRouter.Link to="team"
                        params={{org: this.getParams().org, team: team.name}}
                        query={this.getQuery()}>
                        {team.name}
                    </ReactRouter.Link>
                </li>
            )
        }.bind(this);
        return (
            <div className="menu">
                <ul>
                    <li className="empact">
                        <ReactRouter.Link to="org" params={this.getParams()} className="logo-button">
                            <div className="logo e">e</div>
                            <div className="logo m">m</div>
                            <div className="logo p">p</div>
                            <div className="logo a">a</div>
                            <div className="logo c">c</div>
                            <div className="logo t">t</div>
                        </ReactRouter.Link>
                    </li>
                    <li className="nav header">Organizations:</li>
                    {this.props.orgs.map(renderOrg)}
                    <li className="nav header">Teams:</li>
                    {this.props.teams.map(renderTeam)}
                </ul>
            </div>
        );
    }
});
