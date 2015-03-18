var App = React.createClass({
    mixins: [ReactRouter.Navigation, ReactRouter.State],

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
                    <ReactRouter.RouteHandler />
                </div>
            </div>
        );
    }
});
