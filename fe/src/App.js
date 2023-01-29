import './App.css';
import React from 'react';
import Home from "./components/Home";
import StravaCallback from "./components/stravaCallback";
import {Route, Routes} from "react-router-dom";
import Profile from "./components/profile";
import Register from "./components/register";
import CreateEvent from "./components/CreateEvent";
import Events from "./pages/Events";
import Event from "./pages/Event";
import Segments from "./pages/StarredSegments";
import AddUsers from "./pages/AddUsers";

function App() {
    return (
        <Routes>
            <Route path="/" element={<Home/>}/>

            <Route path="/user/create-event" element={<CreateEvent/>}/>
            <Route path="/user/events" element={<Events/>}/>
            <Route path="/user/event/:id" element={<Event/>}/>
            <Route path="/user/event/:id/add-segments" element={<Segments/>}/>
            <Route path="/user/event/:id/add-users" element={<AddUsers/>}/>

            <Route path="/callback" element={<StravaCallback/>}/>
            <Route path="/profile" element={<Profile/>}/>
            <Route path='/register' element={<Register/>}/>
        </Routes>
    );
}

export default App;