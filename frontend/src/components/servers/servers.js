import React, { Component } from "react";
import apiService from "../../api";
import Server from "./server";

export default class Servers extends Component {
    constructor(props) {
        super(props);

        this.state = {
            agents: []
        };

        this.socket = new WebSocket(`ws://${process.env.REACT_APP_API_SERVER}/api/v1/agent/ws`);
    }

    componentDidMount() {
        this.getAgentList()

        this.socket.onmessage = (event) => {
            const response = event.data;
            let json = JSON.parse(response)
            let currentAgents = this.state.agents
            console.log(json)
            if (json.action === "update_agent") {
                let flag = false;

                if (json.data.status === "deleted") {
                    flag = true;
                    let deletedAgents = currentAgents.filter((currentAgent) => currentAgent.agent_id != json.data.agent_id)
                    currentAgents = deletedAgents
                } else {
                    for (let i = 0; i < currentAgents.length; i++) {
                        if (currentAgents[i].agent_id === json.id) {
                            flag = true;
                            currentAgents[i] = json.data
                            break;
                        }
                    }
                }

                if (!flag) {
                    // is new agent
                    currentAgents.push(json.data)
                }

                this.setState({ agents: currentAgents })
            }
        }
    }

    getAgentList() {
        apiService.getAgents().then(
            response => {
                this.setState({ agents: response.data.data })
            }
        )
    }

    render() {
        return (
            <div>
                <div className="w-1/2 mb-5">
                    <h1 className="text-black text-xl mb-1">Список серверов</h1>

                    <div className="flex flex-col gap-2">
                        {this.state.agents.map(agent => {
                            return (
                                <Server agent={agent} />
                            )
                        })}
                    </div>
                </div>
            </div>
        );
    }
}
