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

    orgsURL: '/api/orgs',
    teamsURL: '/api/teams?org=',
    usersURL: '/api/users?org=',

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
        $.get(this.usersURL + this.getParams().org, function(res){
            this.setState({users: res});
            if (res !== null) {
                for (var i = 0; i < res.length; i++) {
                    var user = res[i];
                    Storage.set('user', user.login, user);
                }
            }
        }.bind(this));
    },

    render: function(){
        return (
            <div className="master">
                <div className="app" id="app">
                    <Menu orgs={this.state.orgs} teams={this.state.teams} />
                    <Router.RouteHandler />
                </div>
                <WeekIntervalSelector />
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
            <div className="menu">
                <ul>
                    <li key="empact">
                        <Link to="org" params={this.getParams()} className="logo-button">
                            <div className="logo e">e</div>
                            <div className="logo m">m</div>
                            <div className="logo p">p</div>
                            <div className="logo a">a</div>
                            <div className="logo c">c</div>
                            <div className="logo t">t</div>
                        </Link>
                    </li>
                    <li key="orgs-header" className="nav header">Organizations:</li>
                    {this.props.orgs.map(renderOrg)}
                    <li key="teams-header" className="nav header">Teams:</li>
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
            infoImage, infoTitle,
            bcApi, bcItems,
            sacApi, sacItems;

        if (p.team) {
            infoTitle = p.team;
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
                <InfoBlock image={infoImage} title={infoTitle} />
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
                <div className={'img'+ (this.props.image ? '' : ' empty')}
                    style={{backgroundImage: "url("+ (this.props.image || '') +")"}} />
                <h1>{this.props.title}</h1>
            </div>
        )
    }
});

var WeekIntervalSelector = React.createClass({
    monthNames: [
        'Jan', 'Feb', 'Mar',
        'Apr', 'May', 'Jun',
        'Jul', 'Aug', 'Sep',
        'Oct', 'Nov', 'Dec'
    ],

    getInitialState: function() {
        return {
            from: 0,
            to: 0
        };
    },

    componentDidMount: function() {
        // window.addEventListener('resize', this.resize);
    },

    showSelector: function() {

    },

    render: function() {
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

        var renderWeek = function(week, i) {
            var d = new Date(week*ms),
                month = this.monthNames[d.getMonth()],
                season = this.seasons[d.getMonth()],
                day = d.getDate();

            return (
                <div key={'week-'+ i} ref="blocks" className="week">
                    <div key={''+ i +'-week'} className="square">{day}</div>
                    <div key={''+ i +'-month'} className={'square '+ season}>{month}</div>
                    <div key={''+ i +'-day'} className="square">{day}</div>
                </div>
            )
        }.bind(this);

        return (
            <div className="week-selector">
                <div className="week">
                    Mar 9
                </div>
                &mdash;
                <div className="week">
                    Apr 27
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
