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
        var topRepos = "/api/stat/orgs/top?org="+ this.getParams().org +"&item=repo",
            repoURL = "/app/"+ this.getParams().org +"/repos/";
        return (
            <section className="content">
                <BarChart api={topRepos} link={repoURL}/>
            </section>
        );
    }
});

var TeamStats = React.createClass({
    mixins: [Router.Navigation, Router.State],
    render: function(){
        var topRepos = "/api/stat/teams/top"+
                "?org="+ this.getParams().org +
                "&team="+ this.getParams().team +
                "&item=repo",
            repoURL = "/app/"+ this.getParams().org +"/repos/";
        return (
            <section className="content">
                <BarChart api={topRepos} link={repoURL}/>
            </section>
        );
    }
});

var UserStats = React.createClass({
    mixins: [Router.Navigation, Router.State],
    render: function(){
        var topRepos = "/api/stat/users/top"+
                "?org="+ this.getParams().org +
                "&author="+ this.getParams().user +
                "&item=repo",
            repoURL = "/app/"+ this.getParams().org +"/repos/";
        return (
            <section className="content">
                <BarChart api={topRepos} link={repoURL}/>
            </section>
        );
    }
});

var RepoStats = React.createClass({
    mixins: [Router.Navigation, Router.State],
    render: function(){
        var topAuthors = "/api/stat/repos/top"+
                "?org="+ this.getParams().org +
                "&repo="+ this.getParams().repo +
                "&item=author",
            userURL = "/app/"+ this.getParams().org +"/users/";
        return (
            <section className="content">
                <BarChart api={topAuthors} link={userURL}/>
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
