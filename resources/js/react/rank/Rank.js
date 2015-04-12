import React from 'react';

class Rank extends React.Component {

  constructor(props) {
      super(props)
  }

  render() {

      var nodes = this.props.items.map(function(item) {
          return (
              <tr key={item.target.id}>
                  <td>{item.rank}</td>
                  <td>{item.target.name}({item.value})</td>
              </tr>
          );
      });

      return (
          <div className="table-responsive">
              <table className="table table-bordered">
                <thead>
                    <tr>
                        <th colSpan="2">{this.props.title}</th>
                    </tr>
                </thead>
                <tbody>
                    {nodes}
                </tbody>
              </table>
          </div>
      );
  }
}

export default Rank
