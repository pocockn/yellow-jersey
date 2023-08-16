import AuthenticationManager from "../services/authManager";
import {useNavigate, useParams} from "react-router-dom";
import React, {useEffect, useState} from "react";
import Sidebar from "../components/sidebar";
import Header from "../components/header";
import {Button} from "@mui/material";
import AddIcon from '@mui/icons-material/Add';
import SegmentMapPolyline from "../components/SegmentMapPolyline";
import TableContainer from "@mui/material/TableContainer";
import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";
import TableCell from "@mui/material/TableCell";
import TableBody from "@mui/material/TableBody";

const Event = () => {
    const authManager = new AuthenticationManager();
    const navigate = useNavigate();
    const [event, setEvent] = useState({
        owner: "",
        name: "",
        users: [],
        segment_ids: [],
        segment_efforts: {}
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

    function getSegmentCompleteCount(userID) {
        return event.segment_efforts[userID].length;
    }

    function getTotalSegmentTimes(userID) {
        let totalTime = 0;
        event.segment_efforts[userID].map((segment) =>
            totalTime += segment.elapsed_time
        )
        return totalTime;
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
        console.log(event)
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

    function fancyTimeFormat(duration) {
        // Hours, minutes and seconds
        const hrs = ~~(duration / 3600);
        const mins = ~~((duration % 3600) / 60);
        const secs = ~~duration % 60;

        // Output like "1:01" or "4:03:59" or "123:03:59"
        let ret = "";

        if (hrs > 0) {
            ret += "" + hrs + ":" + (mins < 10 ? "0" : "");
        }

        ret += "" + mins + ":" + (secs < 10 ? "0" : "");
        ret += "" + secs;

        return ret;
    }

    const populateSegmentWithPBs = () => {
        const requestOptions = {
            method: 'POST',
            headers: {'Content-Type': 'application/json', Authorization: `Bearer ${authManager.getAccessToken()}`},
        };
        fetch(`http://localhost:8080/user/event/` + id + "/segments/populate", requestOptions)
            .then(response => {
                if (response.ok) {
                    console.log(response)
                } else if (response.status === 409) {
                    console.log(response)
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
                            <h1 className="h3 mb-0 text-gray-800">{event.name}</h1>
                        </div>
                        <h3>User Statistics</h3>
                        <TableContainer component={Paper} style={{width: 800}}>
                            <Table sx={{minWidth: 650}} size="small" aria-label="a dense table">
                                <TableHead>
                                    <TableRow>
                                        <TableCell>Name</TableCell>
                                        <TableCell>Segments Completed</TableCell>
                                        <TableCell>Total Time</TableCell>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    {Array.isArray(event.users)
                                        ? event.users.map((user) => (
                                            <TableRow key={user.id}>
                                                <TableCell component="th" scope="row">
                                                    {user.athlete_detailed.firstname} {user.athlete_detailed.lastname}
                                                </TableCell>

                                                <TableCell component="th" scope="row">
                                                    {getSegmentCompleteCount(user.id)}
                                                </TableCell>

                                                <TableCell component="th" scope="row">
                                                    {fancyTimeFormat(getTotalSegmentTimes(user.id))}
                                                </TableCell>
                                            </TableRow>
                                        )) : null}
                                </TableBody>

                            </Table>
                        </TableContainer>
                        <br/>
                        <h3>Event Segments</h3>

                        {Array.isArray(segments)
                            ? segments.map((segment) => (
                                <div className="col-md-3" key={segment.id}>
                                    <SegmentMapPolyline segment={segment} addSegment={null}/>
                                </div>
                            )) : null}

                        <Button
                            size="small"
                            variant="outlined"
                            onClick={() => handleClick("add-segments")}
                            endIcon={<AddIcon/>}> Add Segments
                        </Button>

                        <Button
                            size="small"
                            variant="outlined"
                            onClick={() => handleClick("add-users")}
                            endIcon={<AddIcon/>}> Add Users
                        </Button>

                        <Button
                            size="small"
                            variant="outlined"
                            onClick={() => populateSegmentWithPBs()}
                            endIcon={<AddIcon/>}> Add Segment PBs to event
                        </Button>
                    </div>
                </div>
            </div>
        </div>
    )
}

export default Event;