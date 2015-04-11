var React = require('react');
var Router = require('react-router');
var DefaultRoute = Router.DefaultRoute;
var NotFoundRoute = Router.NotFoundRoute;
var Route = Router.Route;

var SporedApp = require('./react/SporedApp');
var SearchApp = require('./react/search/SearchApp');
var Landing = require('./react/Landing');
var RankApp = require('./react/rank/RankApp');


var routes = (
  <Route name="spored" handler={SporedApp} path="/">
    <DefaultRoute handler={Landing} />
    <Route name="search" handler={SearchApp} path="/search" />
    <Route name="rank" handler={RankApp} path="/rank" />
    <NotFoundRoute handler={SearchApp} />
  </Route>
);

Router.run(routes, Router.HistoryLocation, function(Handler, state) {
    var params = state.params;
    React.render(<Handler params={params} />, document.body);
});
