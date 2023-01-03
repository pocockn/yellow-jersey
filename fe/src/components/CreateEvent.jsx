import React, {useEffect, useState} from 'react';
import Sidebar from "./sidebar";
import Header from "./header";
import AuthenticationManager from "../services/authManager";
import {useNavigate} from "react-router-dom";

const CreateEvent = () => {
    const [eventName, setEventName] = useState("")
    const [message, setMessage] = useState("");
    const authManager = new AuthenticationManager();
    const navigate = useNavigate();

    useEffect(() => {
        if (authManager.getAccessToken() === "") {
            navigate("/")
        }
    })

    let handleSubmit = async (e) => {
        e.preventDefault();
        try {
            let res = await fetch("http://localhost:8080/user/create-event", {
                headers: {Authorization: `Bearer ${authManager.getAccessToken()}`},
                method: "POST",
                body: JSON.stringify({
                    name: eventName,
                }),
            });
            await res
            if (res.status === 200) {
                setEventName("");
                setMessage("Event created successfully");
            } else {
                setMessage("Error occurred creating event");
            }
        } catch (err) {
            console.log(err)
        }
    }


    return (
        <div id="wrapper">
            {<Sidebar/>}
            <div id="content-wrapper" className="d-flex flex-column">
                <div id="content">
                    {<Header/>}
                    <div className="container-fluid">
                        <div className="d-sm-flex align-items-center justify-content-between mb-4">
                            <h1 className="h3 mb-0 text-gray-800">Create Event</h1>
                        </div>
                        <hr></hr>
                        <div className="row">
                            <div className="col-lg-6">
                                <form onSubmit={handleSubmit}>
                                    <div className="form-group">
                                        <label>Event Name</label>
                                        <input
                                            className="form-control"
                                            type="text"
                                            onChange={(e) => setEventName(e.target.value)}
                                        />
                                    </div>
                                    <button type="submit" className="btn btn-success">Create</button>
                                    <div className="message">{message ? <p>{message}</p> : null}</div>
                                </form>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    )

}

export default CreateEvent;