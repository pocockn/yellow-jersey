import AuthenticationManager from "../services/authManager";
import React, {useEffect, useState} from "react";
import Sidebar from "../components/sidebar";
import Header from "../components/header";
import {useNavigate, useParams} from "react-router-dom";
import {Alert, Snackbar} from "@mui/material";
import CardContent from "@mui/material/CardContent";
import Typography from "@mui/material/Typography";
import CardActions from "@mui/material/CardActions";
import Button from "@mui/material/Button";
import Card from "@mui/material/Card";


const AddUsers = () => {
    const authManager = new AuthenticationManager();
    const navigate = useNavigate();
    let {id} = useParams();
    const [openSuccess, setOpenSuccess] = useState(false);
    const [openExists, setOpenExists] = useState(false);

    const [users, setUsers] = useState({
        users: {},
    })

    const handleClose = () => {
        setOpenSuccess(false);
        setOpenExists(false);
    };

    useEffect(() => {
        if (authManager.getAccessToken() === "") {
            navigate("/")
        }

        fetchUsers().catch(err => console.log(err));
    }, [])

    async function fetchUsers() {
        const resp = await fetch(`http://localhost:8080/user/users`, {
            method: 'GET',
            mode: 'cors',
            headers: {
                Authorization: `Bearer ${authManager.getAccessToken()}`,
            },
        })
        const respJSON = await resp.json()
        console.log(respJSON)
        setUsers(respJSON.users)
    }

    const addUser = (user_id) => {
        const requestOptions = {
            method: 'PUT',
            headers: {'Content-Type': 'application/json', Authorization: `Bearer ${authManager.getAccessToken()}`},
        };
        fetch(`http://localhost:8080/user/event/` + id + "/users/" + user_id, requestOptions)
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
                            <h1 className="h3 mb-0 text-gray-800">Add Users to event</h1>
                        </div>
                        <div className="d-sm-4 align-items-center justify-content-between mb-4">
                            <p>Below are a list of users on Yellow-Jersey. Add the user to your event to allow them
                                to compete.</p>
                            <hr></hr>
                        </div>
                        <div className="row">
                            {Array.isArray(users)
                                ? users.map((user) => (
                                    <div className="col-md-3" key={user.id}>
                                        <Card sx={{m: 2}}>
                                            <CardContent>
                                                <hr></hr>
                                                <Typography gutterBottom variant="h7" component="div">
                                                    {user.id}
                                                </Typography>
                                            </CardContent>
                                            <CardActions>
                                                <Button onClick={() => addUser(user.id)} size="small">Add</Button>
                                            </CardActions>
                                        </Card>
                                    </div>
                                )) : null}
                        </div>
                        <Snackbar
                            anchorOrigin={{
                                horizontal: "left",
                                vertical: "bottom",
                            }}
                            open={openSuccess}
                            autoHideDuration={3000}
                        >
                            <Alert onClose={handleClose} severity="success" sx={{width: '100%'}}>
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
                            <Alert onClose={handleClose} severity="warning" sx={{width: '100%'}}>
                                Segment has already been added to event
                            </Alert>
                        </Snackbar>
                    </div>
                </div>
            </div>
        </div>
    )

}

export default AddUsers;