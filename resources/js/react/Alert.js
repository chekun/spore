import React from 'react';

class Alert extends React.Component {

    constructor(props) {
        super(props)
    }

    render() {
        var className = 'alert alert-'+this.props.type + ' ' + this.props.className;
        return (
            <div className={className}  role="alert">{this.props.text}</div>
        );
    }
}

export default Alert
