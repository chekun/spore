import React from 'react';

class SearchLoadMore extends React.Component {

    constructor(props) {
        super(props)
    }

    render() {
        var className = 'btn btn-default ' + this.props.className;
        return (
            <button className={className} type="button" onClick={this.props.onClick}>加载更多...</button>
        );
    }
}

export default SearchLoadMore
