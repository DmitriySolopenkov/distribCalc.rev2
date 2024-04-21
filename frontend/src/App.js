import './App.css';
import { Routes, Route, Link } from "react-router-dom";
import Tasks from "./components/tasks/tasks";
import Servers from "./components/servers/servers";
import {Component} from "react";
import apiService from "./api";
import Login from "./components/login";
import {withRouter} from "./with-router";
import Register from "./components/register";

class App extends Component{
    constructor(props) {
        super(props);
        this.state = {
            currentUser: undefined,
        }
    }

    componentDidMount() {
        const user = apiService.getCurrentUser();
        if (user) {
            this.setState({
                currentUser: user,
            })
        }
    }

    render() {
        if (this.state.currentUser === undefined) {
            return (
                <div>
                    <Routes>
                        <Route path="/" element={<Login />} />
                        <Route path="/register" element={<Register />} />
                    </Routes>
                </div>
            )
        }
        return (
            <div className="h-screen w-screen overflow-x-hidden">
                <div className="text-white rounded-lg lg:w-[900px] m-auto my-5">
                    <header className="text-black font-semibold text-xl p-5 flex justify-between items-center bg-white mb-5 rounded-md block-shadow">
                        <div className="w-1/2 flex gap-3">
                            <Link to={"/"} className="w-32 text-center rounded-md py-2 text-sm">ЗАДАЧИ</Link>
                            <Link to={"/servers"} className="w-32 text-center rounded-md py-2 text-sm">СЕРВЕРА</Link>
                            <span onClick={(e) => {
                                localStorage.removeItem("user");
                                this.setState({
                                    currentUser: undefined,
                                });
                            }} className="w-32 text-center rounded-md py-2 text-sm text-gray">ВЫЙТИ</span>
                        </div>
                        <a href="https://github.com/DmitriySolopenkov/distribCalc.rev2" className="text-gray text-sm font-normal">by Dmitriy Solopenkov</a>
                    </header>
                    <div className="p-5 bg-white rounded-md block-shadow">
                        <Routes>
                            <Route path="/" element={<Tasks />} />
                            <Route path="/servers" element={<Servers />} />
                        </Routes>
                    </div>
                </div>
            </div>
        );
    }
}

export default withRouter(App);