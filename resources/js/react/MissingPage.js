import React from 'react';
var Router = require('react-router');
var Link = Router.Link;

class MissingPage extends React.Component {

    constructor(props) {
        super(props)
    }

    render() {
        return (
            <div className="container landing-container">
                <h2>404 Not Found!</h2>
                <div>
                    <Link to="/search" className="btn btn-info .btn-lg">孢子大搜索</Link>
                    <Link to="/rank" className="btn btn-info .btn-lg">孢子风云榜</Link>
                </div>
            </div>
        );
    }
}

export default MissingPage
