import React, {useEffect} from 'react';
import {useNavigate} from 'react-router-dom'
import AuthenticationManager from "../services/authManager"
import axios from "axios";

const StravaCallback = () => {
    const navigate = useNavigate();
    const authManager = new AuthenticationManager();

    useEffect(() => {
        const queryParams = new URLSearchParams(window.location.search);
        const state = queryParams.get('state');
        const scope = queryParams.get('scope');
        const code = queryParams.get('code');

        if (code !== "") {
            axios.post(`http://localhost:8080/exchange_token`, {
                code: code,
                state: state,
                scope: scope,
            })
                .then(res => {
                    authManager.updateToken(res.data.token);
                    console.log(res);
                }).catch(err => {
                console.log(err)
            })
        }

        navigate("/")
    })
}

export default StravaCallback;