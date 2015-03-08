var Router = ReactRouter,
    Route = Router.Route,
    Link = Router.Link,
    RouteHandler = Router.RouteHandler,
    DefaultRoute = Router.DefaultRoute,
    NotFoundRoute = Router.NotFoundRoute;

var App = React.createClass({
    mixins: [Router.Navigation],
    render: function(){
        return (
            <section className="app">
                <Menu/>
                <RouteHandler/>
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
            <RouteHandler/>
        );
    }
});

var OrgStats = React.createClass({
    mixins: [Router.Navigation, Router.State],
    render: function(){
        return (
            <section className="content">Org stats for {this.getParams().org}</section>
        );
    }
});

var TeamStats = React.createClass({
    mixins: [Router.Navigation, Router.State],
    render: function(){
        var url = "/api/stat/teams/top?org="+ this.getParams().org +"&from=1417878086&to=1425654067";
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
        <Route name="root" path="/app/" handler={App}>
            <Route name="dashboard" path=":org" handler={Dashboard}>
                <DefaultRoute handler={OrgStats} />
                <Route name="team" path="teams/:team" handler={TeamStats} />
                <Route name="user" path="users/:user" handler={UserStats} />
                <Route name="repo" path="repos/:repo" handler={RepoStats} />
            </Route>
            <NotFoundRoute handler={NotFound}/>
        </Route>
    ];
Router.run(routes, Router.HistoryLocation, function(Handler) {
    React.render(<Handler/>, document.body);
});
