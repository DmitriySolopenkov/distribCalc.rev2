import React, { Component } from "react";
import Form from "react-validation/build/form";
import CheckButton from "react-validation/build/button";
import apiService from "../api";
import {Link, useNavigate} from "react-router-dom";
import {withRouter} from "../with-router";

class Register extends Component {
    constructor(props) {
        super(props);
        this.handleRegister = this.handleRegister.bind(this);
        this.onChangeLogin = this.onChangeLogin.bind(this);
        this.onChangePassword = this.onChangePassword.bind(this);

        this.state = {
            login: "",
            password: "",
            loading: false,
            message: ""
        };
    }

    onChangeLogin(e) {
        this.setState({
            login: e.target.value
        });
    }

    onChangePassword(e) {
        this.setState({
            password: e.target.value
        });
    }

    handleRegister(e) {
        e.preventDefault();

        this.setState({
            message: "",
            loading: true
        });

        this.form.validateAll();

        if (this.checkBtn.context._errors.length === 0) {
            apiService.register(this.state.login, this.state.password).then(
                (data) => {
                    console.log(data)
                    if (data.code !== 200) {
                        this.setState({
                            loading: false,
                            message: "Неверный логин или пароль"
                        })
                    } else {
                        if (data.data.token) {
                            localStorage.setItem("user", JSON.stringify(data.data));
                        }

                        this.props.router.navigate("/");
                        window.location.reload();
                    }
                },
                error => {
                    const resMessage =
                        (error.response &&
                            error.response.data &&
                            error.response.data.message) ||
                        error.message ||
                        error.toString();

                    this.setState({
                        loading: false,
                        message: resMessage
                    });
                }
            );
        } else {
            this.setState({
                loading: false
            });
        }
    }

    render() {
        return (
            <div className="absolute h-screen w-screen">
                <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 ">
                    <div className="bg-white p-5 rounded-sm">
                        <Form
                            onSubmit={this.handleRegister}
                            ref={c => {
                                this.form = c;
                            }}
                        >
                            <input
                                type="text"
                                className="h-12 w-full bg-dark-white focus:outline-none rounded-lg px-5 mb-2"
                                name="username"
                                placeholder="Логин"
                                value={this.state.login}
                                onChange={this.onChangeLogin}
                                required={true}></input>

                            <input
                                type="password"
                                name="password"
                                className="h-12 w-full bg-dark-white focus:outline-none rounded-lg px-5 mb-2"
                                placeholder="Пароль"
                                value={this.state.password}
                                onChange={this.onChangePassword}
                                required={true}></input>

                            {this.state.message && (
                                <div className="border border-red-600 rounded-lg p-2 mb-2 text-red-600 text-center" role="alert">
                                    {this.state.message}
                                </div>
                            )}

                            <button
                                className="w-full h-10 block text-white font-medium rounded-lg"
                                style={{backgroundColor: "#000"}}
                                disabled={this.state.loading}
                            >
                                {this.state.loading && (
                                    <span className="spinner-border spinner-border-sm"></span>
                                )}
                                <span>Регистрация</span>
                            </button>
                            <Link to={"/"} className={"block w-100 text-center"}>Уже есть аккаунт?</Link>
                            <CheckButton
                                style={{ display: "none" }}
                                ref={c => {
                                    this.checkBtn = c;
                                }}
                            />
                        </Form>
                    </div>
                </div>
            </div>
        );
    }
}

export default withRouter(Register);