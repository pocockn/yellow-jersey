import React, {useEffect, useState} from 'react';
import Sidebar from "../components/sidebar";
import Header from "../components/header";
import AuthenticationManager from "../services/authManager";
import {useNavigate} from "react-router-dom";
import EventsTable from "../components/EventsTable";
import axios from "axios";

const Events = () => {
    const authManager = new AuthenticationManager();
    const navigate = useNavigate();
    const [events, setEvents] = useState([])

    useEffect(() => {
        if (authManager.getAccessToken() === "") {
            navigate("/")
        }

        fetchEvents();
    }, [])

    const handleEventRoute = (path) => {
        navigate("/user/event/" + path);
    }

    const fetchEvents = () => {
        axios.get(`http://localhost:8080/user/events`, {
            headers: {Authorization: `Bearer ${authManager.getAccessToken()}`}
        }).then(res => {
            console.log(res);
            setEvents(res.data.events)
        })
    }

    return (
        <div id="wrapper">
            {<Sidebar/>}
            <div id="content-wrapper" className="d-flex flex-column">
                <div id="content">
                    {<Header/>}
                    <div className="container-fluid">
                        <div className="d-sm-flex align-items-center justify-content-between mb-4">
                            <h1 className="h3 mb-0 text-gray-800">All events</h1>
                        </div>
                        <hr></hr>
                        <EventsTable events={events} handleRoute={handleEventRoute}/>
                    </div>
                </div>
            </div>
        </div>
    )

}

export default Events;