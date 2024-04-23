import React from 'react'

export default function Server({ agent }) {
    const lastPing = new Date(agent.last_ping);

    const addLeadingZero = (number) => {
        return number < 10 ? '0' + number : number;
    };
    const formattedLastPing = `${addLeadingZero(lastPing.getDate())}.${addLeadingZero(lastPing.getMonth() + 1)}.${lastPing.getFullYear()} ${addLeadingZero(lastPing.getHours())}:${addLeadingZero(lastPing.getMinutes())}:${addLeadingZero(lastPing.getSeconds())}`;

    const statusServer = {
        "connected": "#08c708",
        "disconnected": "#c70808",
    };

    return (
        <div className="flex gap-2">
            <span className="inline-block w-12 h-12 rounded-full" style={{ backgroundColor: statusServer[agent.status] }}></span>
            <div className="flex flex-col justify-between">
                <span className="text-black text-xl font-semibold">{agent.agent_id} <span className="text-gray text-sm font-normal">{formattedLastPing}</span></span>
                <span style={{ color: statusServer[agent.status] }}>{agent.status}</span>
            </div>
        </div>
    )
}
