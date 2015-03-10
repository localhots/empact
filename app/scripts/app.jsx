var Router = ReactRouter,
    Link = Router.Link;

var Storage = {
    set: function(category, key, value) {
        window.localStorage.setItem(category +'-'+ key, JSON.stringify(value));
    },

    get: function(category, key) {
        var val = window.localStorage.getItem(category +'-'+ key);
        return val === null ? null : JSON.parse(val);
    }
};

var App = React.createClass({
    mixins: [Router.Navigation, Router.State],

    orgsURL: "/api/orgs",
    teamsURL: "/api/teams?org=",

    getInitialState: function() {
        return {
            orgs: [],
            org: null,
            teams: [],
            team: null
        };
    },

    componentDidMount: function() {
        this.loadOrgs();
        this.loadTeams();
    },

    loadOrgs: function() {
        $.get(this.orgsURL, function(res){
            this.setState({orgs: res});
            if (res !== null) {
                for (var i = 0; i < res.length; i++) {
                    var org = res[i];
                    Storage.set('org', org.login, org);
                }
            }
        }.bind(this));
    },

    loadTeams: function() {
        $.get(this.teamsURL + this.getParams().org, function(res){
            this.setState({teams: res})
            if (res !== null) {
                for (var i = 0; i < res.length; i++) {
                    var team = res[i];
                    Storage.set('team', team.name, team);
                }
            }
        }.bind(this));
    },

    render: function(){
        return (
            <section className="app">
                <Menu orgs={this.state.orgs} teams={this.state.teams} />
                <Router.RouteHandler />
            </section>
        );
    }
});

var Menu = React.createClass({
    mixins: [Router.State],

    render: function() {
        var renderOrg = function(org) {
            return (
                <li key={'org-'+ org.login} className="nav org">
                    <Link to="org" params={{org: org.login}}>{org.login}</Link>
                </li>
            )
        };
        var renderTeam = function(team) {
            return (
                <li key={'team-'+ team.name} className="nav team">
                    <Link to="team" params={{org: team.owner, team: team.name}}>{team.name}</Link>
                </li>
            )
        };
        return (
            <section className="menu">
                <ul>
                    <li key="empact" className="nav empact">
                        <Link to="org" params={this.getParams()}>Empact</Link>
                    </li>
                    <li key="orgs-header" className="nav header">Organizations:</li>
                    {this.props.orgs.map(renderOrg)}
                    <li key="teams-header" className="nav header">Teams:</li>
                    {this.props.teams.map(renderTeam)}
                </ul>
            </section>
        );
    }
});

var Org = React.createClass({
    render: function(){
        return (
            <Router.RouteHandler />
        );
    }
});

var OrgStats = React.createClass({
    mixins: [Router.State],

    render: function(){
        var org = Storage.get('org', this.getParams().org);
        return (
            <section className="content">
                <InfoBlock image={org.avatar_url} title={org.login} text={org.descr} />
                <BarChart key={this.getParams().team} api="/api/stat/orgs/top"
                    params={this.getParams()} items={["repo", "team", "user"]} />
            </section>
        );
    }
});

var TeamStats = React.createClass({
    mixins: [Router.State],

    render: function(){
        return (
            <section className="content">
                <InfoBlock key={"info-block-"+ this.getParams().team}
                    image="https://media.licdn.com/mpr/mpr/p/8/005/058/14b/0088c48.jpg"
                    title={this.getParams().team}
                    text={"The most awesome team in "+ this.getParams().org} />
                <BarChart key={'bar-chart-'+ this.getParams().team} api="/api/stat/teams/top"
                    params={this.getParams()} items={["repo", "user"]} />
            </section>
        );
    }
});

var UserStats = React.createClass({
    mixins: [Router.State],
    render: function(){
        return (
            <section className="content">
                <InfoBlock title={this.getParams().user} />
                <BarChart key={'bar-chart-'+ this.getParams().user} api="/api/stat/users/top"
                    params={this.getParams()} items={["repo"]} />
            </section>
        );
    }
});

var RepoStats = React.createClass({
    mixins: [Router.State],
    render: function(){
        return (
            <section className="content">
                <InfoBlock title={this.getParams().repo} />
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

var SelectOrg = React.createClass({
    render: function(){
        return (
            <section className="content">Please select organization from the menu!</section>
        );
    }
});

var InfoBlock = React.createClass({
    render: function() {
        var img = <div className="img" style={{backgroundImage: "url("+ this.props.image +")"}} />;
        return (
            <div className="info-block">
                { this.props.image ? img : null }
                <h1>{this.props.title}</h1>
                <p>{this.props.text}</p>
            </div>
        )
    }
});

var routes = [
    <Router.Route name="root" path="/app/" handler={App}>
        <Router.DefaultRoute handler={SelectOrg} />
        <Router.NotFoundRoute handler={NotFound} />
        <Router.Route name="org" path=":org" handler={Org}>
            <Router.DefaultRoute handler={OrgStats} />
            <Router.Route name="team" path="teams/:team" handler={TeamStats} />
            <Router.Route name="user" path="users/:user" handler={UserStats} />
            <Router.Route name="repo" path="repos/:repo" handler={RepoStats} />
        </Router.Route>
    </Router.Route>
];
Router.run(routes, Router.HistoryLocation, function(Handler) {
    React.render(<Handler />, document.body);
});
