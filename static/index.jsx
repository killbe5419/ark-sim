import React from "react";
import ReactDOM from "react-dom";
import axios from "axios";

class Case_One extends React.Component {
    render() {
        if(this.props.data.loaded) {
            const  imgUrl = "http://localhost:8966/static/images/character/ark-icon/" + this.props.data.content.Codename + ".png";
            return (
                <div className="case">
                    <p> Name: { this.props.data.content.Codename } </p>
                    <p> Rare: { this.props.data.content.Rare } </p>
                    <p> Class: { this.props.data.content.Class } </p>
                    <img src={ imgUrl } alt={ this.props.data.content.name }/>
                </div>
            );
        }
        return null;
    }
}

class App extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            data:{
                loaded: false,
            }
        }
    }

    handlePickOne = () => {
        const data = {
            params: {
                method: "pickOne",
                fromClient: true
            }
        }
        axios.get("/pickOne",data)
            .then(res => {
                console.log(res.data);
                this.setState({
                    data:{
                        loaded: true,
                        content: res.data
                    }
                })
            })
    }

    render() {
        return (
            <div>
                <h1>Hello, Golang!</h1>
                <button onClick={ () => console.log("Hello, Golang!") }>Click here to say Hello!</button>
                <button onClick={ this.handlePickOne }>Pick One Card</button>
                <Case_One data={ this.state.data }/>
            </div>
        );
    }
}

ReactDOM.render(
    <App />,
    document.getElementById("root")
)