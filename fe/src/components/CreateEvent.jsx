import React, {useEffect, useState} from 'react';
import Sidebar from "./sidebar";
import Header from "./header";
import AuthenticationManager from "../services/authManager";
import {useNavigate} from "react-router-dom";
import {Alert, Snackbar} from "@mui/material";

const CreateEvent = () => {
    const [eventName, setEventName] = useState("")
    const [startDate, setStartDate] = useState("")
    const [finishDate, setFinishDate] = useState("")
    const [open, setOpen] = useState(false);
    const authManager = new AuthenticationManager();
    const navigate = useNavigate();

    useEffect(() => {
        if (authManager.getAccessToken() === "") {
            navigate("/")
        }
    })

    const handleClose = () => {
        setOpen(false);
    };

    let handleSubmit = async (e) => {
        e.preventDefault();
        try {
            let res = await fetch("http://localhost:8080/user/create-event", {
                headers: {Authorization: `Bearer ${authManager.getAccessToken()}`},
                method: "POST",
                body: JSON.stringify({
                    name: eventName,
                    start_date: new Date(startDate),
                    finish_date: new Date(finishDate)
                }),
            });
            await res
            if (res.status === 200) {
                setEventName("");
                setStartDate("");
                setFinishDate("");
                setOpen(true);
            } else {
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
                                        <label>Start Date</label>
                                        <input
                                            type="date"
                                            className="form-control date"
                                            name="start"
                                            onChange={(e) => setStartDate(e.target.value)}
                                        />
                                        <label>Finish Date</label>
                                        <input
                                            type="date"
                                            className="form-control date"
                                            name="finish"
                                            onChange={(e) => setFinishDate(e.target.value)}
                                        />
                                    </div>
                                    <button type="submit" className="btn btn-success">Create</button>
                                </form>
                            </div>
                        </div>
                        <Snackbar
                            anchorOrigin={{
                                horizontal: "left",
                                vertical: "bottom",
                            }}
                            open={open}
                            autoHideDuration={3000}
                        >
                            <Alert onClose={handleClose} severity="success" sx={{width: '100%'}}>
                                Event successfully created
                            </Alert>
                        </Snackbar>
                    </div>
                </div>
            </div>
        </div>
    )

}

export default CreateEvent;