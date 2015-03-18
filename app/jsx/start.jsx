var Org = React.createClass({
    render: function(){
        return (
            <ReactRouter.RouteHandler />
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

var routes = [
    <ReactRouter.Route name="root" path="/app/" handler={App}>
        <ReactRouter.DefaultRoute handler={SelectOrg} />
        <ReactRouter.NotFoundRoute handler={NotFound} />
        <ReactRouter.Route name="org" path=":org" handler={Org}>
            <ReactRouter.DefaultRoute handler={Dashboard} />
            <ReactRouter.Route name="team" path="teams/:team" handler={Dashboard} />
            <ReactRouter.Route name="user" path="users/:user" handler={Dashboard} />
            <ReactRouter.Route name="repo" path="repos/:repo" handler={Dashboard} />
        </ReactRouter.Route>
    </ReactRouter.Route>
];

ReactRouter.run(routes, ReactRouter.HistoryLocation, function(Handler) {
    React.render(<Handler />, document.body);
});
