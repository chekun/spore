import React from 'react';
var Router = require('react-router');

var RouteHandler = Router.RouteHandler;
var Link = Router.Link;

class SporedApp extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            location: location.pathname
        }
    }

    componentWillReceiveProps(props) {
        this.setState({location: location.pathname});
    }

    render() {
        return (
            <div>
                <nav className="navbar navbar-default navbar-spored navbar-fixed-top">
                  <div className="container">
                    <div className="navbar-header">
                      <Link className="navbar-brand" to="/">Project Spore</Link>
                    </div>
                    <div id="navbar" className="navbar-collapse collapse">
                      <ul className="nav navbar-nav navbar-right">
                        <li className={this.state.location === '/' ? "active" : ''} id="route-index"><Link to="/">首页</Link></li>
                        <li className={this.state.location === '/search' ? "active" : ''} id="route-search"><Link to="/search">搜索</Link></li>
                        <li className={this.state.location === '/rank' ? "active" : ''} id="route-rank"><Link to="/rank">排行榜</Link></li>
                      </ul>
                    </div>
                  </div>
                </nav>

                <div className="container container-spored">
                    <RouteHandler {...this.props} />
                </div>

                <footer className="footer"></footer>
            </div>
        );
    }
}

export default SporedApp
