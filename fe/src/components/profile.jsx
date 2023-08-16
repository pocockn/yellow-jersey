import {useNavigate} from "react-router-dom";
import AuthenticationManager from "../services/authManager";
import React, {useEffect, useState} from "react";
import axios from "axios";
import Sidebar from "./sidebar";
import Header from "./header";

const Profile = () => {
    const navigate = useNavigate();
    const authManager = new AuthenticationManager();
    const [routes, setRoutes] = useState([])

    useEffect(() => {
        if (authManager.getAccessToken() === "") {
            navigate("/")
        }

        console.log(authManager.getAccessToken())

        axios.get(`http://localhost:8080/user/routes`, {
            headers: {Authorization: `Bearer ${authManager.getAccessToken()}`}
        }).then(res => {
            console.log(res);
            setRoutes(res.data.routes)
        })
    })

    return routes.map((route) => {
        return (
            <div id="wrapper">
                {<Sidebar/>}
                <div id="content-wrapper" className="d-flex flex-column">
                    <div id="content">
                        {<Header/>}
                        <p>{route[0]}</p>
                    </div>
                </div>
            </div>
        );
    });
}

export default Profile;