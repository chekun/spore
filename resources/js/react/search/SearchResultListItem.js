import React from 'react';
var Avatar = require('../Avatar');

//<SearchResultListItem item=[] />

class SearchResultListItem extends React.Component {
    render() {
        var item = this.props.item;
        if (item["title"]) {
            item.type = 3;
        } else if (item["intro"]) {
            item.type = 2;
        } else {
            item.type = 1;
        }
        switch (item.type) {
            case 1:
                return (
                    <a href="javascript:void(0);" className="list-group-item">
                        <Avatar width="30" height="30" icon={item.icon} alt={item.name} />
                        <span>用户</span>: {item.name}
                    </a>
                );
                break;
            case 2:
                var link = "http://baoz.cn/" + item.id;
                return (
                    <a href={link} target="_blank" className="list-group-item">
                        <Avatar width="30" height="30" icon={item.icon} alt={item.name} />
                        <span>群组</span>: {item.name}
                    </a>
                );
                break;
            case 3:
                var link = "http://baoz.cn/" + item.id;
                return (
                    <a href={link} target="_blank" className="list-group-item">
                        <span>帖子</span>: {item.title}
                    </a>
                );
        }
    }
}

export default SearchResultListItem
