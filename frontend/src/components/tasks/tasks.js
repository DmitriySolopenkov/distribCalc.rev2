import React, { Component } from "react";
import Task from "./task";
import apiService from "./../../api";

export default class Tasks extends Component {
    constructor(props) {
        super(props);

        this.state = {
            tasks: [],
            expression: "",
        };

        this.handleAddTask = this.handleAddTask.bind(this);
        this.getTaskList = this.getTaskList.bind(this);

        this.socket = new WebSocket(`ws://${process.env.REACT_APP_API_SERVER}/api/v1/agent/ws`);
    }

    handleExpressionChange = (event) => {
        this.setState({ expression: event.target.value });
    };

    handleAddTask() {
        this.setState({ expression: "" })

        apiService.addTask(this.state.expression)
    }

    getTaskList() {
        apiService.getTasks().then(
            response => {
                this.setState({ tasks: response.data.data })
            }
        )
    }

    componentDidMount() {
        this.getTaskList()

        this.socket.onmessage = (event) => {
            const response = event.data;
            let json = JSON.parse(response)
            let currentTasks = this.state.tasks
            switch (json.action) {
                case "new_task":
                    // add new task to list
                    this.setState({ tasks: [...[json.data].concat(currentTasks)] })
                    break;
                case "update_task":
                    // update task in list

                    for (let i = 0; i < currentTasks.length; i++) {
                        if (currentTasks[i].task_id === json.id) {
                            currentTasks[i] = json.data
                            break;
                        }
                    }

                    this.setState({ tasks: currentTasks })
                    break;
                default:
                    break;
            }
        }
    }

    render() {
        return (
            <div>
                <div className="lg:w-1/2 mb-5">
                    <h1 className="text-black text-xl mb-1">Добавить новую задачу</h1>
                    <input className="border border-gray w-1/2 rounded-md p-1 mr-1 text-black" placeholder="Введите пример" onChange={this.handleExpressionChange} value={this.state.expression}></input>
                    <button className="bg-blue p-1 rounded-md px-5" onClick={this.handleAddTask}>Решить</button>
                </div>
                <div>
                    <h1 className="text-black text-xl mb-1">Список всех задач</h1>
                    <div className="flex flex-col gap-2">
                        {this.state.tasks.map(task => {
                            return (
                                <Task task={task} />
                            )
                        })}
                    </div>
                </div>
            </div>
        );
    }
}
