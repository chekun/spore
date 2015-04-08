import React from 'react';
var progress = require('nprogress');
var $ = require('jquery');

class Loading extends React.Component {

    componentDidMount() {
        $(document).ajaxStart(function() {
            progress.start();
        });
        $(document).ajaxComplete(function() {
            progress.done();
        });
    }

    render() {
        return (
            <div className="none"></div>
        );
    }
}

export default Loading
