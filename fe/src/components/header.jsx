import React from "react";
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome'
import {faCogs, faList, faSignOut, faUser} from "@fortawesome/free-solid-svg-icons";
import Search from "./search";
import profileImage from '../assets/img/undraw_profile.svg';

const Header = () => {

    return (
        <nav className="navbar navbar-expand navbar-light bg-white topbar mb-4 static-top shadow">
            <button id="sidebarToggleTop" className="btn btn-link d-md-none rounded-circle mr-3">
                <i className="fa fa-bars"></i>
            </button>

            <Search/>

            <ul className="navbar-nav ml-auto">

                <div className="topbar-divider d-none d-sm-block"></div>

                <li className="nav-item dropdown no-arrow" key="nav-dropdown">
                    <a className="nav-link dropdown-toggle" href="#" id="userDropdown" role="button"
                       data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                        <span className="mr-2 d-none d-lg-inline text-gray-600 small">Nick Pocock</span>
                        <img className="img-profile rounded-circle"
                             src={profileImage}/>
                    </a>
                    <div className="dropdown-menu dropdown-menu-right shadow animated--grow-in"
                         aria-labelledby="userDropdown">
                        <a className="dropdown-item" href="#">
                            <i><FontAwesomeIcon icon={faUser}/></i>
                            Profile
                        </a>
                        <a className="dropdown-item" href="#">
                            <i><FontAwesomeIcon icon={faCogs}/></i>
                            Settings
                        </a>
                        <a className="dropdown-item" href="#">
                            <i><FontAwesomeIcon icon={faList}/></i>
                            Activity Log
                        </a>
                        <div className="dropdown-divider"></div>
                        <a className="dropdown-item" href="#" data-toggle="modal" data-target="#logoutModal">
                            <i><FontAwesomeIcon icon={faSignOut}/></i>
                            Logout
                        </a>
                    </div>
                </li>
            </ul>

        </nav>
    )
}

export default Header;