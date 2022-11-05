import AuthenticationManager from "../services/authManager";
import {useNavigate, useParams} from "react-router-dom";
import React, {useEffect, useState} from "react";
import axios from "axios";
import Sidebar from "../components/sidebar";
import Header from "../components/header";
import {Button} from "@mui/material";
import AddIcon from '@mui/icons-material/Add';

const Event = () => {
    const authManager = new AuthenticationManager();
    const navigate = useNavigate();
    const [event, setEvent] = useState({})
    let {id} = useParams();

    useEffect(() => {
        if (authManager.getAccessToken() === "") {
            navigate("/")
        }

        fetchEvent();
    }, [])

    function handleClick(path) {
        navigate("/user/event/" + id + "/" + path);
    }

    const fetchEvent = () => {
        axios.get(`http://localhost:8080/user/event/` + id, {
            headers: {Authorization: `Bearer ${authManager.getAccessToken()}`}
        }).then(res => {
            console.log(res);
            setEvent(res.data.event)
        })
    }

    let addUsersButton;
    if (Array.isArray(event.users)) {
        addUsersButton = <p>{event.users}</p>;
    } else {
        addUsersButton = <Button size="small" variant="outlined" endIcon={<AddIcon/>}> Add Users </Button>;
    }

    let addSegmentsButton;
    if (Array.isArray(event.segment_ids)) {
        addSegmentsButton = <p>{event.segment_ids}</p>;
    } else {
        addSegmentsButton = <Button
            size="small"
            variant="outlined"
            onClick={() => handleClick("add-segments")}
            endIcon={<AddIcon/>}> Add Segments
        </Button>;
    }

    return (
        <div id="wrapper">
            {<Sidebar/>}
            <div id="content-wrapper" className="d-flex flex-column">
                <div id="content">
                    {<Header/>}
                    <div className="container-fluid">
                        <div className="d-sm-flex align-items-center justify-content-between mb-4">
                            <h1 className="h3 mb-0 text-gray-800">{event.name}</h1>
                        </div>
                        <hr></hr>
                        {addUsersButton}
                        <hr></hr>
                        {addSegmentsButton}
                    </div>
                </div>
            </div>
        </div>
    )

}

export default Event;