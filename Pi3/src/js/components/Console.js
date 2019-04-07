import React from "react";

class Console extends React.Component {
    render() {
        return (
            <div>
                <p>
                    <strong>Status:</strong> Connecting...
                </p>
                <div className="console__container">
                    <video className="console" />
                    <div className="console__overlay">
                        <h2>Connecting to Remote Machine</h2>
                        <h2 className="console__loading_container">
                            <span className="console__loading">.</span>&nbsp;
                            <span className="console__loading">.</span>&nbsp;
                            <span className="console__loading">.</span>&nbsp;
                        </h2>
                    </div>
                </div>
            </div>
        )
    }
}

export default Console;