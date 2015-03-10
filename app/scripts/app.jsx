var Router = ReactRouter,
    Link = Router.Link;

var App = React.createClass({
    mixins: [Router.Navigation],
    render: function(){
        return (
            <section className="app">
                <Menu/>
                <Router.RouteHandler/>
            </section>
        );
    }
});

var Menu = React.createClass({
    mixins: [Router.Navigation, Router.State],
    api_url: "/api/teams?org=",

    getInitialState: function() {
        return {teams: []};
    },

    componentDidMount: function() {
        this.loadTeams();
    },

    loadTeams: function() {
        $.get(this.api_url + this.getParams().org, function(res){
            this.setState({teams: res})
        }.bind(this));
    },

    render: function() {
        var renderTeam = function(team) {
            return (
                <li key={team.name} className="nav team">
                    <Link to="team" params={{org: team.owner, team: team.name}}>{team.name}</Link>
                </li>
            )
        };
        return (
            <section className="menu">
                <ul>
                    <li key="empact" className="nav empact">
                        <Link to="org" params={{org: this.getParams().org}}>Empact</Link>
                    </li>
                    <li key="dash" className="nav dash">
                        <Link to="org" params={{org: this.getParams().org}}>{this.getParams().org}</Link>
                    </li>
                    <li key="teams-header" className="nav header">Teams:</li>
                    {this.state.teams.map(renderTeam)}
                </ul>
            </section>
        );
    }
});

var Org = React.createClass({
    render: function(){
        return (
            <Router.RouteHandler/>
        );
    }
});

var OrgStats = React.createClass({
    mixins: [Router.Navigation, Router.State],
    render: function(){
        return (
            <section className="content">
                <BarChart key={this.getParams().team} api="/api/stat/orgs/top"
                    params={this.getParams()} items={["repo", "team", "user"]} />
            </section>
        );
    }
});

var TeamStats = React.createClass({
    mixins: [Router.Navigation, Router.State],
    render: function(){
        return (
            <section className="content">
                <BarChart key={this.getParams().team} api="/api/stat/teams/top"
                    params={this.getParams()} items={["repo", "user"]} />
            </section>
        );
    }
});

var UserStats = React.createClass({
    mixins: [Router.Navigation, Router.State],
    render: function(){
        return (
            <section className="content">
                <BarChart key={this.getParams().team} api="/api/stat/users/top"
                    params={this.getParams()} items={["repo"]} />
            </section>
        );
    }
});

var RepoStats = React.createClass({
    mixins: [Router.Navigation, Router.State],
    render: function(){
        return (
            <section className="content">
                <BarChart key={this.getParams().team} api="/api/stat/repos/top"
                    params={this.getParams()} items={["user", "team"]} />
            </section>
        );
    }
});

var NotFound = React.createClass({
    render: function(){
        return (
            <section className="content">NOT FOUND :(</section>
        );
    }
});

var routes = [
        <Router.Route name="root" path="/app/" handler={App}>
            <Router.Route name="org" path=":org" handler={Org}>
                <Router.DefaultRoute handler={OrgStats} />
                <Router.Route name="team" path="teams/:team" handler={TeamStats} />
                <Router.Route name="user" path="users/:user" handler={UserStats} />
                <Router.Route name="repo" path="repos/:repo" handler={RepoStats} />
            </Router.Route>
            <Router.NotFoundRoute handler={NotFound}/>
        </Router.Route>
    ];
Router.run(routes, Router.HistoryLocation, function(Handler) {
    React.render(<Handler/>, document.body);
});
