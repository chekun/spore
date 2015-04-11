import React from 'react';
var Rank = require('./Rank');

class RankGroup extends React.Component {

  constructor(props) {
      super(props)
  }

  render() {
      var $this = this;
      var nodes = this.props.items.map(function(item) {
          return (
              <div className={$this.props.className} key={item.key}>
                  <Rank title={item.title} items={item.items} />
              </div>
          );
      });

      return (
          <div className="row">
            {nodes}
          </div>
      );
  }
}

export default RankGroup
