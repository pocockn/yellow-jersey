import AuthenticationManager from "../services/authManager";
import {useNavigate, useParams} from "react-router-dom";
import React, {useEffect, useState} from "react";
import Sidebar from "../components/sidebar";
import Header from "../components/header";
import {Button} from "@mui/material";
import AddIcon from '@mui/icons-material/Add';
import SegmentMapPolyline from "../components/SegmentMapPolyline";

const Event = () => {
    const authManager = new AuthenticationManager();
    const navigate = useNavigate();
    const [event, setEvent] = useState({
        owner: "",
        name: "",
        users: [],
        segment_ids: []
    })
    const [segments, setSegments] = useState({
        segments: {},
        segment_ids: [],
    })

    let {id} = useParams();

    useEffect(() => {
        if (authManager.getAccessToken() === "") {
            navigate("/")
        }

        fetchEvent().catch(err => console.log(err));
        fetchSegments().catch(err => console.log(err));
    }, [])

    function handleClick(path) {
        navigate("/user/event/" + id + "/" + path);
    }

    async function fetchEvent() {
        const resp = await fetch(`http://localhost:8080/user/event/` + id, {
            method: 'GET',
            mode: 'cors',
            headers: {
                Authorization: `Bearer ${authManager.getAccessToken()}`,
                'Content-Type': 'application/json'
            },
        })
        const event = await resp.json()
        setEvent(event.event)
    }

    async function fetchSegments() {
        const resp = await fetch(`http://localhost:8080/user/event/` + id + `/segments`, {
            method: 'POST',
            mode: 'cors',
            headers: {
                Authorization: `Bearer ${authManager.getAccessToken()}`,
                'Content-Type': 'application/json'
            },
        })
        const segments = await resp.json()
        setSegments(segments.segments)
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
                        <p>{event.users}</p>
                        <Button size="small" variant="outlined" endIcon={<AddIcon/>}> Add Users </Button>
                        <hr></hr>
                        {Array.isArray(segments)
                            ? segments.map((segment) => (
                                <div className="col-md-3">
                                    <SegmentMapPolyline segment={segment} addSegment={null}/>
                                </div>
                            )) : null}
                        <Button
                            size="small"
                            variant="outlined"
                            onClick={() => handleClick("add-segments")}
                            endIcon={<AddIcon/>}> Add Segments
                        </Button>
                    </div>
                </div>
            </div>
        </div>
    )
}

export default Event;