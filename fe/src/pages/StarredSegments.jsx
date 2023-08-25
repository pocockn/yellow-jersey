import AuthenticationManager from "../services/authManager";
import React, {useEffect, useState} from "react";
import axios from "axios";
import Sidebar from "../components/sidebar";
import Header from "../components/header";
import {useNavigate, useParams} from "react-router-dom";
import SegmentMapPolyline from "../components/SegmentMapPolyline";
import {Alert, Snackbar} from "@mui/material";


const Segments = () => {
    const authManager = new AuthenticationManager();
    const navigate = useNavigate();
    let {id} = useParams();
    const [openSuccess, setOpenSuccess] = useState(false);
    const [openExists, setOpenExists] = useState(false);

    const [segments, setSegments] = useState({
        segments: {},
        segment_ids: [],
    })

    const handleClose = () => {
        setOpenSuccess(false);
        setOpenExists(false);
    };

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
            const newState = res.data.segments.map(obj => {
                return {...obj, segments: res.data.segments};
            });
            setSegments(newState);
        })
    }

    const addSegment = (segment_id) => {
        const requestOptions = {
            method: 'PUT',
            headers: {'Content-Type': 'application/json', Authorization: `Bearer ${authManager.getAccessToken()}`},
        };
        fetch(`http://localhost:8080/user/event/` + id + "/segment/" + segment_id, requestOptions)
            .then(response => {
                if (response.ok) {
                    setOpenSuccess(true)
                } else if (response.status === 409) {
                    setOpenExists(true)
                }
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
                            <h1 className="h3 mb-0 text-gray-800">Add segments to event</h1>
                        </div>
                        <div className="d-sm-4 align-items-center justify-content-between mb-4">
                            <p>Below are a list of your starred segments on Strava. Add the segment to your event to
                                include it in your race. The user with the lowest time on all included segments is the
                                yellow jersey!</p>
                            <hr></hr>
                        </div>
                        <div className="row">
                            {Array.isArray(segments)
                                ? segments.map((segment) => (
                                <div className="col-md-3" key={segment.id}>
                                    <SegmentMapPolyline segment={segment} addSegment={addSegment}/>
                                </div>
                            )) : null }
                        </div>
                        <Snackbar
                            anchorOrigin={{
                                horizontal: "left",
                                vertical: "bottom",
                            }}
                            open={openSuccess}
                            autoHideDuration={3000}
                        >
                            <Alert onClose={handleClose} severity="success" sx={{ width: '100%' }}>
                                Segment successfully added
                            </Alert>
                        </Snackbar>
                        <Snackbar
                            anchorOrigin={{
                                horizontal: "left",
                                vertical: "bottom",
                            }}
                            open={openExists}
                            autoHideDuration={3000}
                        >
                            <Alert onClose={handleClose} severity="warning" sx={{ width: '100%' }}>
                                Segment has already been added to event
                            </Alert>
                        </Snackbar>
                    </div>
                </div>
            </div>
        </div>
    )

}

export default Segments;