import React from 'react';
import GameControl from './GameControl'

class GameControls extends React.Component {

    render() {
        return (
            <div>
                <header>controls:</header>
                <GameControl action={this.handleClick} icon={'>>'}></GameControl>
            </div>
        );
    }

    handleClick(event) {
        alert(event.toString());
    }
}

export default GameControls;