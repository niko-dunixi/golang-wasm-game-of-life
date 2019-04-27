import React from 'react';
import logo from './logo.svg';
import './App.css';
import GameCanvas from "./GameCanvas";
import GameControls from "./GameControls";

function App() {
    return (
        <div className="App">
            <header className="Game-header">
                <h1>Conway's Game of Life</h1>
                <h2>A <a href="https://golang.org/" target="_blank" className="GoLang-Formatting">GoLang</a> and <a
                    href="https://webassembly.org/" target="_blank"
                    className="WebAsm-Formatting">WebAssembly</a> implementation</h2>
            </header>
            <GameControls/>
            <GameCanvas/>
        </div>
    );
}

export default App;
