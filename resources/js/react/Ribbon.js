import React from 'react';

class Ribbon extends React.Component {

  constructor(props) {
      super(props)
  }

  render() {

      return (
          <a href={this.props.repo}>
            <img className="ribbon" src="/public/assets/ribbon.png" alt="Fork me on Github" />
          </a>
      );
  }
}

export default Ribbon
