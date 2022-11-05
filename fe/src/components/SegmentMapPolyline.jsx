import React from 'react';
import {GoogleMap, LoadScript, Polyline} from '@react-google-maps/api';

const SegmentMapPolyline = () => {
    const pathCoordinates = [
        {lat: 36.05298765935, lng: -112.083756616339},
        {lat: 36.2169884797185, lng: -112.056727493181}
    ];

    const mapContainerStyle = {
        height: "400px",
        width: "800px"
    };

    const center = {
        lat: 0,
        lng: -180
    };


    return (
        <LoadScript
            googleMapsApiKey="AIzaSyBnYuLtIYFw1UNzGJ1RO7JOVhzY9SuGE0o"
        >
        <GoogleMap
            id="marker-example"
            mapContainerStyle={mapContainerStyle}
            zoom={2}
            center={center}
        >
            {/*for creating path with the updated coordinates*/}
            <Polyline
                path={pathCoordinates}
                geodesic={true}
                options={{
                    strokeColor: "#ff2527",
                    strokeOpacity: 0.75,
                    strokeWeight: 2,
                }}
            />
        </GoogleMap>
        </LoadScript>
    )
}

export default SegmentMapPolyline;