import React from 'react';

//<SearchTabControl badges=[] />

class SearchTabControl extends React.Component {

    constructor(props) {
        super(props);
    }

    render() {
        return (
            <div className={this.props.className} onClick={this.props.onClick}>
                <ul className="nav nav-tabs nav-justified">
                  <li role="presentation" className={this.props.selected == 0 ? "active" : ""}>
                    <a href="javascript:void(0);" data-tab="0">
                        <span data-tab="0">全部  </span>
                        <span className="badge" data-tab="0">{this.props.badges.total}</span>
                    </a>
                  </li>
                  <li role="presentation" className={this.props.selected == 1 ? "active" : ""}>
                    <a href="javascript:void(0);" data-tab="1">
                        <span data-tab="1">用户  </span>
                        <span data-tab="1" className="badge">{this.props.badges.users}</span>
                    </a>
                  </li>
                  <li role="presentation" className={this.props.selected == 2 ? "active" : ""}>
                    <a href="javascript:void(0);" data-tab="2">
                        <span data-tab="2">群组  </span>
                        <span data-tab="2" className="badge">{this.props.badges.groups}</span>
                    </a>
                  </li>
                  <li role="presentation" className={this.props.selected == 3 ? "active" : ""}>
                    <a href="javascript:void(0);" data-tab="3">
                        <span data-tab="3">帖子  </span>
                        <span data-tab="3" className="badge">{this.props.badges.threads}</span>
                    </a>
                  </li>
                </ul>
            </div>
        );
    }
}

export default SearchTabControl
