import React from 'react';

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
        var avatarPrefix = 'http://#region.as.baoz.cn/f/#file';
        if (item.icon) {
            item.icon.url = avatarPrefix.replace("#region", item.icon.crop.substr(0, 1)).replace("#file", item.icon.crop)+'.t30x30.png';
        }
        switch (item.type) {
            case 1:
                return (
                    <a href="javascript:void(0);" className="list-group-item">
                        <img width="30" height="30" src={item.icon.url} alt={item.name} className="img-circle" />
                        <span>用户</span>: {item.name}
                    </a>
                );
                break;
            case 2:
                var link = "http://baoz.cn/" + item.id;
                return (
                    <a href={link} target="_blank" className="list-group-item">
                        <img width="30" height="30" src={item.icon.url} alt={item.name} className="img-circle" />
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
