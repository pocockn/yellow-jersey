import React, {useEffect, useState} from 'react';
import axios from "axios";

const Register = () => {
    const [authURL, setAuthURL] = useState("")

    useEffect(() => {
        axios.get(`http://localhost:8080/auth`)
            .then(res => {
                console.log(res);
                const url = res.data;
                setAuthURL(url);
            })
    })

    return (
        <div className="container">
            <div className="row">
                <div className="col-lg-6 offset-lg-3">
                    <div className="p-5">
                        <div className="text-center">
                            <h1 className="h4 text-gray-900 mb-4">Connect with Strava!</h1>
                        </div>
                        <form className="user">
                            <a href={authURL} className="btn btn-google btn-user btn-block">
                                <i className="fab fa-google fa-fw"></i> Register or login with Strava
                            </a>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    )
}

export default Register;