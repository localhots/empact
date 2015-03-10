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
    mixins: [Router.State],
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
    mixins: [Router.State],
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
    mixins: [Router.State],
    render: function(){
        return (
            <section className="content">
                <InfoBlock key={"info-block-"+ this.getParams().team}
                    image="https://media.licdn.com/mpr/mpr/p/8/005/058/14b/0088c48.jpg"
                    title={this.getParams().team}
                    text={"The most awesome team in "+ this.getParams().org} />
                <BarChart key={"bar-chart"+ this.getParams().team} api="/api/stat/teams/top"
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
                <BarChart key={this.getParams().team} api="/api/stat/users/top"
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
