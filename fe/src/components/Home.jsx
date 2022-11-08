import React from 'react';
import axios from "axios";
import Sidebar from "./sidebar";
import Header from "./header";

export default class Home extends React.Component {
    state = {
        authURL: ""
    }

    componentDidMount() {
        axios.get(`http://localhost:8080/auth`)
            .then(res => {
                console.log(res);
                const url = res.data;
                this.setState({authURL: url});
            })
    }

    render() {
        return (
            <div id="wrapper">
                {<Sidebar/>}
                <div id="content-wrapper" className="d-flex flex-column">
                    <div id="content">
                        {<Header/>}
                        <div className="container-fluid">
                            <div className="d-sm-flex align-items-center justify-content-between mb-4">
                                <h1 className="h3 mb-0 text-gray-800">Home</h1>
                            </div>
                            <div className="row">
                                <a href={this.state.authURL}>
                                    <button>Connect with Strava</button>
                                </a>
                            </div>
                            <hr></hr>
                        </div>
                    </div>
                </div>
            </div>
        );
    };
}