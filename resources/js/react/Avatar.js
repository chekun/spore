import React from 'react';

//<Avatar icon=[] />

class Avatar extends React.Component {

    constructor(props) {
        super(props)
    }

    render() {

        var avatarPrefix = 'http://#region.as.baoz.cn/f/#file';
        var url = 'http://d.as.baoz.cn/f/default.t30x30.png';
        if (this.props.icon) {
            url = avatarPrefix.replace("#region", this.props.icon.crop.substr(0, 1)).replace("#file", this.props.icon.crop)+'.t30x30.png';
        }

        return (
            <img width="30" height="30" src={url} alt={this.props.alt} className="img-circle" />
        );
    }
}

export default Avatar
