import React from 'react';
var RankGroup = require('./RankGroup');

class RankApp extends React.Component {

  constructor(props) {
      super(props)
      this.state = {
          user_ranks: [],
          group_ranks: []
      }
  }

  render() {
      return (
          <div className="search">
            <h1>排行榜</h1>
            <RankGroup items={this.state.user_ranks} className="col-xs-6 col-md-3" />
            <RankGroup items={this.state.group_ranks} className="col-xs-6 col-md-4" />
          </div>
      );
  }
}

export default RankApp
