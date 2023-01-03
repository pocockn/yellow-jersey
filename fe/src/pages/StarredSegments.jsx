import AuthenticationManager from "../services/authManager";
import React, {useEffect, useState} from "react";
import axios from "axios";
import Sidebar from "../components/sidebar";
import Header from "../components/header";
import {useNavigate, useParams} from "react-router-dom";
import SegmentMapPolyline from "../components/SegmentMapPolyline";


const Segments = () => {
    const authManager = new AuthenticationManager();
    const navigate = useNavigate();
    let {id} = useParams();
    const [segments, setSegments] = useState([])

    useEffect(() => {
        if (authManager.getAccessToken() === "") {
            navigate("/")
        }

        fetchSegments();
    }, [])

    const fetchSegments = () => {
        axios.get(`http://localhost:8080/user/segments`, {
            headers: {Authorization: `Bearer ${authManager.getAccessToken()}`}
        }).then(res => {
            setSegments(res.data.segments);
        })
    }

    const addSegment = (segment_id) => {
        const requestOptions = {
            method: 'PUT',
            headers: {'Content-Type': 'application/json', Authorization: `Bearer ${authManager.getAccessToken()}`},
        };
        fetch(`http://localhost:8080/user/event/` + id + "/segment/" + segment_id, requestOptions)
            .then(response => console.log(response))
    }

    return (
        <div id="wrapper">
            {<Sidebar/>}
            <div id="content-wrapper" className="d-flex flex-column">
                <div id="content">
                    {<Header/>}
                    <div className="container-fluid">
                        <div className="d-sm-flex align-items-center justify-content-between mb-4">
                            <h1 className="h3 mb-0 text-gray-800">Add segments to event</h1>
                        </div>
                        <div className="d-sm-4 align-items-center justify-content-between mb-4">
                            <p>Below are a list of your starred segments on Strava. Add the segment to your event to
                                include it in your race. The user with the lowest time on all included segments is the
                                yellow jersey!</p>
                        </div>
                        <div className="row">
                            {Array.isArray(segments)
                                ? segments.map((segment) => (
                                    <div className="col-md-3">
                                        <SegmentMapPolyline segment={segment} addSegment={addSegment}/>
                                    </div>
                                )) : null}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    )

}

export default Segments;