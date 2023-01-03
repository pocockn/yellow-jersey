import React from 'react';
import {MapContainer, Polyline, TileLayer} from "react-leaflet";
import 'leaflet/dist/leaflet.css';
import CardContent from "@mui/material/CardContent";
import Typography from "@mui/material/Typography";
import CardActions from "@mui/material/CardActions";
import Button from "@mui/material/Button";
import Card from "@mui/material/Card";
import polyUtil from "polyline-encoded"

const SegmentMapPolyline = (props) => {
    const center = [props.segment.start_latlng[0], props.segment.end_latlng[1]];
    const polyline = polyUtil.decode(props.segment.map.polyline);
    const redOptions = {color: 'red'}
    const kilometres = Math.round(props.segment.distance / 1000)

    return (
        <Card sx={{m: 2}}>
            <CardContent>
                <MapContainer style={{height: '250px', width: '100wh'}} center={center} zoom={11.5}
                              scrollWheelZoom={false}>
                    <TileLayer
                        attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                    />
                    <Polyline pathOptions={redOptions} positions={polyline}/>
                </MapContainer>
                <hr></hr>
                <Typography gutterBottom variant="h7" component="div">
                    {props.segment.name}
                </Typography>
                <Typography gutterBottom variant="body2" color="text.secondary">
                    Length: {kilometres}km<br/>
                    Average Gradient: {props.segment.average_grade}<br/>
                    Country: {props.segment.country}
                </Typography>
            </CardContent>
            <CardActions>
                <Button size="small">Add</Button>
                {props.addSegment !== null ?
                    <Button onClick={() => props.addSegment(props.segment.id)} size="small">Add</Button> : null
                }
            </CardActions>
        </Card>
    )
}

export default SegmentMapPolyline;