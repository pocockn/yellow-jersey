import React from "react";
import {faCog, faTachometer, faWrench} from "@fortawesome/free-solid-svg-icons";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import logo from '../assets/img/logo.png' // relative path to image

const Sidebar = () => {

    return (
        <ul className="navbar-nav bg-gradient-primary sidebar sidebar-dark accordion" id="accordionSidebar">

            <a className="sidebar-brand d-flex align-items-center justify-content-center" href="/" key="logo">
                <img src={logo} className="img-fluid"></img>
            </a>

            <hr className="sidebar-divider my-0" key="divider"/>

            <li className="nav-item active" key="dashboard">
                <a className="nav-link" href="/">
                    <i><FontAwesomeIcon icon={faTachometer}/></i>
                    <span>Dashboard</span></a>
            </li>

            <hr className="sidebar-divider" key="sidebar-divider"/>

            <div className="sidebar-heading" key="sidebar-heading">
                Events
            </div>

            <li className="nav-item" key="create-event">
                <a className="nav-link" href="/user/create-event"
                   aria-expanded="true" aria-controls="collapseTwo">
                    <i><FontAwesomeIcon icon={faCog}/></i>
                    <span>Create Event</span>
                </a>
            </li>

            <li className="nav-item" key="events">
                <a className="nav-link" href="/user/events">
                    <i><FontAwesomeIcon icon={faWrench}/></i>
                    <span>Events</span>
                </a>
            </li>

            <hr className="sidebar-divider" key="another-divider"/>

        </ul>
    )
}

export default Sidebar;