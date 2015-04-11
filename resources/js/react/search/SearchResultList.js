import React from 'react';
var SearchResultListItem = require('./SearchResultListItem');

//<SearchResultList items=[] />

class SearchResultList extends React.Component {

    constructor(props) {
        super(props);
    }

    render() {
        var itemNodes = this.props.items.map(function(item) {
            return (
                <SearchResultListItem key={item.id} item={item} />
            );
        });
        var className = "list-group " + (this.props.className);
        return (
            <div className={className}>
                {itemNodes}
            </div>
        );
    }
}

export default SearchResultList
