var Router = ReactRouter,
    Route = Router.Route,
    Link = Router.Link,
    RouteHandler = Router.RouteHandler,
    DefaultRoute = Router.DefaultRoute,
    NotFoundRoute = Router.NotFoundRoute;

var App = React.createClass({displayName: "App",
    mixins: [Router.Navigation],
    render: function(){
        return (
            React.createElement("section", {className: "app"}, 
                React.createElement(Menu, null), 
                React.createElement(RouteHandler, null)
            )
        );
    }
});

var Menu = React.createClass({displayName: "Menu",
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
                React.createElement("li", {key: team.slug, className: "nav team"}, 
                    React.createElement(Link, {to: "team", params: {org: team.owner, team: team.slug}}, team.name)
                )
            )
        };
        return (
            React.createElement("section", {className: "menu"}, 
                React.createElement("ul", null, 
                    React.createElement("li", {key: "empact", className: "nav empact"}, 
                        React.createElement(Link, {to: "dashboard", params: {org: this.getParams().org}}, "Empact")
                    ), 
                    React.createElement("li", {key: "dash", className: "nav dash"}, 
                        React.createElement(Link, {to: "dashboard", params: {org: this.getParams().org}}, "Dashboard")
                    ), 
                    React.createElement("li", {key: "teams-header", className: "nav header"}, "Teams:"), 
                    this.state.teams.map(renderTeam)
                )
            )
        );
    }
});

var Dashboard = React.createClass({displayName: "Dashboard",
    render: function(){
        return (
            React.createElement(RouteHandler, null)
        );
    }
});

var OrgStats = React.createClass({displayName: "OrgStats",
    mixins: [Router.Navigation, Router.State],
    render: function(){
        return (
            React.createElement("section", {className: "content"}, "Org stats for ", this.getParams().org)
        );
    }
});

var TeamStats = React.createClass({displayName: "TeamStats",
    mixins: [Router.Navigation, Router.State],
    render: function(){
        var url = "/api/stat/teams/top?org="+ this.getParams().org +"&from=1417878086&to=1425654067";
        return (
            React.createElement("section", {className: "content"}, 
                React.createElement("div", {className: "left"}, 
                    React.createElement(BarChart, {url: url})
                )
            )
        );
    }
});

var UserStats = React.createClass({displayName: "UserStats",
    render: function(){
        return (
            React.createElement("section", {className: "content"}, "User stats!")
        );
    }
});

var RepoStats = React.createClass({displayName: "RepoStats",
    render: function(){
        return (
            React.createElement("section", {className: "content"}, "Repo Stats!")
        );
    }
});

var NotFound = React.createClass({displayName: "NotFound",
    render: function(){
        return (
            React.createElement("section", {className: "content"}, "NOT FOUND :(")
        );
    }
});

var routes = [
        React.createElement(Route, {name: "root", path: "/app/", handler: App}, 
            React.createElement(Route, {name: "dashboard", path: ":org", handler: Dashboard}, 
                React.createElement(DefaultRoute, {handler: OrgStats}), 
                React.createElement(Route, {name: "team", path: "teams/:team", handler: TeamStats}), 
                React.createElement(Route, {name: "user", path: "users/:user", handler: UserStats}), 
                React.createElement(Route, {name: "repo", path: "repos/:repo", handler: RepoStats})
            ), 
            React.createElement(NotFoundRoute, {handler: NotFound})
        )
    ];
Router.run(routes, Router.HistoryLocation, function(Handler) {
    React.render(React.createElement(Handler, null), document.body);
});
