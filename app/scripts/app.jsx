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
                <li key={team.slug} className="nav team">
                    <Link to="team" params={{org: team.owner, team: team.slug}}>{team.name}</Link>
                </li>
            )
        };
        return (
            <section className="menu">
                <ul>
                    <li key="empact" className="nav empact">
                        <Link to="dashboard" params={{org: this.getParams().org}}>Empact</Link>
                    </li>
                    <li key="dash" className="nav dash">
                        <Link to="dashboard" params={{org: this.getParams().org}}>Dashboard</Link>
                    </li>
                    <li key="teams-header" className="nav header">Teams:</li>
                    {this.state.teams.map(renderTeam)}
                </ul>
            </section>
        );
    }
});

var Dashboard = React.createClass({
    render: function(){
        return (
            <Router.RouteHandler/>
        );
    }
});

var OrgStats = React.createClass({
    mixins: [Router.Navigation, Router.State],
    render: function(){
        var topTeams = "/api/stat/teams/top?org="+ this.getParams().org;
        return (
            <section className="content">
                <BarChart url={topTeams} />
            </section>
        );
    }
});

var TeamStats = React.createClass({
    mixins: [Router.Navigation, Router.State],
    render: function(){
        var url = "/api/stat/teams/top?org="+ this.getParams().org;
        return (
            <section className="content">
                <div className="left">
                    <BarChart url={url} />
                </div>
            </section>
        );
    }
});

var UserStats = React.createClass({
    render: function(){
        return (
            <section className="content">User stats!</section>
        );
    }
});

var RepoStats = React.createClass({
    render: function(){
        return (
            <section className="content">Repo Stats!</section>
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
            <Router.Route name="dashboard" path=":org" handler={Dashboard}>
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
