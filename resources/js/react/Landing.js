import React from 'react';
var Router = require('react-router');
var Link = Router.Link;

class Landing extends React.Component {

    constructor(props) {
        super(props)
    }

    componentDidMount() {
        document.getElementById('slogan').className="magictime puffIn";
    }

    render() {
        return (
            <div className="container landing-container">
                <img id="slogan" src="/public/assets/slogan.png" />
                <div>
                    <Link to="/search" className="btn btn-info .btn-lg">孢子大搜索</Link>
                    <Link to="/rank" className="btn btn-info .btn-lg">孢子风云榜</Link>
                </div>
            </div>
        );
    }
}

export default Landing
