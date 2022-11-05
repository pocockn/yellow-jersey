import AuthenticationManager from "../services/authManager";
import CardContent from '@mui/material/CardContent';
import CardMedia from '@mui/material/CardMedia';
import Button from '@mui/material/Button';
import Typography from '@mui/material/Typography';
import React, {useEffect, useState} from "react";
import axios from "axios";
import Card from '@mui/material/Card';
import Sidebar from "../components/sidebar";
import Header from "../components/header";
import {useNavigate, useParams} from "react-router-dom";
import CardActions from '@mui/material/CardActions';
import sacalobra from '../assets/img/sacalobra.jpg' // relative path to image


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
                                include it in your race. The user with the lowest team on all included segments is the
                                yellow jersey!</p>
                        </div>
                        <div className="row">
                            {segments.map((segment) => (
                                <div className="col-md-3">
                                    <Card sx={{ m: 2 }}>
                                        <CardMedia
                                            component="img"
                                            height="140"
                                            image={sacalobra}
                                            alt="green iguana"
                                        />
                                        <CardContent>
                                            <Typography gutterBottom variant="h5" component="div">
                                                {segment.name}
                                            </Typography>
                                            <Typography variant="body2" color="text.secondary">
                                                <p>Length: {segment.distance}km</p>
                                                <p>Average Gradient: {segment.average_grade}</p>
                                                <p>Country: {segment.country}</p>
                                            </Typography>
                                        </CardContent>
                                        <CardActions>
                                            <Button size="medium">Add</Button>
                                        </CardActions>
                                    </Card>
                                </div>
                            ))}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    )

}

export default Segments;