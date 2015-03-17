var Router = ReactRouter,
    Link = Router.Link;

var Storage = {
    set: function(category, key, value) {
        window.localStorage.setItem(category +'-'+ key, JSON.stringify(value));
    },

    get: function(category, key) {
        var val = window.localStorage.getItem(category +'-'+ key);
        return val === null ? {} : JSON.parse(val);
    }
};

var App = React.createClass({
    mixins: [Router.Navigation, Router.State],

    orgsURL: '/api/orgs',
    teamsURL: '/api/teams',
    usersURL: '/api/users',

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
        this.loadUsers();
    },

    loadOrgs: function() {
        getURL(this.orgsURL, {}, function(res) {
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
        getURL(this.teamsURL, {org: this.getParams().org}, function(res) {
            this.setState({teams: res});
            if (res !== null) {
                for (var i = 0; i < res.length; i++) {
                    var team = res[i];
                    Storage.set('team', team.name, team);
                }
            }
        }.bind(this));
    },

    loadUsers: function() {
        getURL(this.usersURL, {org: this.getParams().org}, function(res) {
            this.setState({users: res});
            if (res !== null) {
                for (var i = 0; i < res.length; i++) {
                    var user = res[i];
                    Storage.set('user', user.login, user);
                }
            }
        }.bind(this));
    },

    render: function() {
        return (
            <div className="master">
                <div className="app" id="app">
                    <Menu orgs={this.state.orgs} teams={this.state.teams} />
                    <Router.RouteHandler />
                </div>
            </div>
        );
    }
});

var Menu = React.createClass({
    mixins: [Router.State],

    render: function() {
        var renderOrg = function(org) {
            return (
                <li key={'org-'+ org.login} className="nav org">
                    <Link to="org"
                        params={{org: org.login}}
                        query={this.getQuery()}>
                        {org.login}
                    </Link>
                </li>
            )
        }.bind(this);
        var renderTeam = function(team) {
            return (
                <li key={'team-'+ team.name} className="nav team">
                    <Link to="team"
                        params={{org: team.owner, team: team.name}}
                        query={this.getQuery()}>
                        {team.name}
                    </Link>
                </li>
            )
        }.bind(this);
        return (
            <div className="menu">
                <ul>
                    <li className="empact">
                        <Link to="org" params={this.getParams()} className="logo-button">
                            <div className="logo e">e</div>
                            <div className="logo m">m</div>
                            <div className="logo p">p</div>
                            <div className="logo a">a</div>
                            <div className="logo c">c</div>
                            <div className="logo t">t</div>
                        </Link>
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

var Org = React.createClass({
    render: function(){
        return (
            <Router.RouteHandler />
        );
    }
});

var Dashboard = React.createClass({
    mixins: [Router.State],

    render: function(){
        var p = this.getParams(),
            infoImage, infoImageClass, infoTitle,
            bcApi, bcItems,
            sacApi, sacItems;

        if (p.team) {
            infoTitle = p.team +' Team';
            infoImageClass = 'team';
            bcApi = '/api/stat/teams/top';
            bcItems = ['repo', 'user'],
            sacApi = '/api/stat/teams/activity';
            sacItems = ['user', 'repo'];
        } else if (p.user) {
            var info = Storage.get('user', p.user);
            infoImage = info ? info.avatar_url : null;
            infoTitle = info && info.name ? info.name : p.user;
            bcApi = '/api/stat/users/top';
            bcItems = ['repo'],
            sacApi = '/api/stat/users/activity';
            sacItems = ['repo'];
        } else if (p.repo) {
            infoTitle = p.repo;
            infoImageClass = 'repo';
            bcApi = '/api/stat/repos/top';
            bcItems = ['user', 'team'],
            sacApi = '/api/stat/repos/activity';
            sacItems = ['user', 'team'];
        } else {
            var info = Storage.get('org', p.org);
            infoImage = info.avatar_url;
            infoTitle = info.login;
            bcApi = '/api/stat/orgs/top';
            bcItems = ['repo', 'team', 'user'],
            sacApi = '/api/stat/orgs/activity';
            sacItems = ['team', 'user', 'repo'];
        }

        return (
            <div className="content">
                <InfoBlock image={infoImage} className={infoImageClass} title={infoTitle} />
                <BarChart api={bcApi} params={this.getParams()} items={bcItems} />
                <StackedAreaChart api={sacApi} params={this.getParams()} items={sacItems} />
            </div>
        );
    }
});

var NotFound = React.createClass({
    render: function(){
        return (
            <div className="content">NOT FOUND :(</div>
        );
    }
});

var SelectOrg = React.createClass({
    render: function(){
        return (
            <div className="content">Please select organization from the menu!</div>
        );
    }
});

var InfoBlock = React.createClass({
    render: function() {
        return (
            <div className="info-block">
                <div className={'img'+ (this.props.image ? '' : ' empty') +' '+ this.props.className}
                    style={this.props.image ? {backgroundImage: "url("+ this.props.image +")"} : null} />
                <h1>{this.props.title}</h1>
            </div>
        )
    }
});

var WeekIntervalSelector = React.createClass({
    mixins: [ReactRouter.Navigation, ReactRouter.State],

    getInitialState: function() {
        var ms = 1000,
            daySeconds = 86400,
            weekSeconds = daySeconds*7,
            today = new Date(),
            sunday = new Date(today - daySeconds*ms*today.getDay()),
            perfectSunday = new Date(Date.UTC(sunday.getFullYear(), sunday.getMonth(), sunday.getDate())),
            lastWeek = perfectSunday.setHours(0)/ms,
            firstWeek = lastWeek - 51*weekSeconds;

        var weeks = [];
        for (var i = lastWeek; i >= firstWeek; i -= weekSeconds) {
            weeks.push(i);
        };

        return {
            weeks: weeks.sort()
        };
    },

    handleChange: function(thing, e) {
        var params = this.getQuery();
        params[thing.slice(0, 1)] = e.target.value/100;
        this.transitionTo(document.location.pathname, null, params);
    },

    render: function() {
        var daySeconds = 86400,
            weekSeconds = daySeconds*7,
            lastWeek = this.state.weeks[this.state.weeks.length-1],
            from = (this.getQuery().f ? parseInt(this.getQuery().f, 10)*100 : lastWeek - 7*weekSeconds),
            to = (this.getQuery().t ? parseInt(this.getQuery().t, 10)*100 : lastWeek);

        var weeksBefore = _(this.state.weeks)
            .filter(function(week) {
                return week < to;
            })
            .reverse()
            .value();
        var weeksAfter = _(this.state.weeks)
            .filter(function(week) {
                return week > from;
            })
            .reverse()
            .value();

        var renderOption = function(ts) {
            return (
                <option key={ts} value={ts}>{formatDate(ts, true)}</option>
            );
        };

        return (
            <div className="week-selector">
                <span>from</span>
                <div ref="from" className="selector">
                    <em ref="label">{formatDate(from)}</em>
                    <select ref="select" value={from} onChange={this.handleChange.bind(this, 'from')}>
                        {weeksBefore.map(renderOption)}
                    </select>
                </div>
                <span>to</span>
                <div ref="to" className="selector">
                    <em ref="label">{formatDate(to)}</em>
                    <select ref="select" value={to} onChange={this.handleChange.bind(this, 'to')}>
                        {weeksAfter.map(renderOption)}
                    </select>
                </div>
            </div>
        );
    }
});

var routes = [
    <Router.Route name="root" path="/app/" handler={App}>
        <Router.DefaultRoute handler={SelectOrg} />
        <Router.NotFoundRoute handler={NotFound} />
        <Router.Route name="org" path=":org" handler={Org}>
            <Router.DefaultRoute handler={Dashboard} />
            <Router.Route name="team" path="teams/:team" handler={Dashboard} />
            <Router.Route name="user" path="users/:user" handler={Dashboard} />
            <Router.Route name="repo" path="repos/:repo" handler={Dashboard} />
        </Router.Route>
    </Router.Route>
];
Router.run(routes, Router.HistoryLocation, function(Handler) {
    React.render(<Handler />, document.body);
});
