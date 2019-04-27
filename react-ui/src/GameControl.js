import React from 'react';

class GameControl extends React.Component {
    constructor(props) {
        super(props);
        this.action = props.action;
    }

    render() {
        return (<a onClick={this.props.action}>{this.props.icon}</a>);
    }
}

export default GameControl;