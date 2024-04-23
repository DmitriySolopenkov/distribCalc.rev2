import React from 'react'

export default function Task({ task }) {
    const date = new Date(task.created_at);

    const addLeadingZero = (number) => {
        return number < 10 ? '0' + number : number;
    };
    const formattedDate = `${addLeadingZero(date.getDate())}.${addLeadingZero(date.getMonth() + 1)}.${date.getFullYear()} ${addLeadingZero(date.getHours())}:${addLeadingZero(date.getMinutes())}:${addLeadingZero(date.getSeconds())}`;

    const statusColor = {
        "completed": "bg-[#08c708]",
        "created": "bg-[#fbff03]",
        "processed": "bg-[#ffae00]",
        "fail": "bg-[#c70808]"
    };

    const taskColor = statusColor[task.status];

    return (
        <div className="border border-gray flex items-center rounded-md p-3 justify-between">
            <div className="flex items-center gap-5">
                <span className={`w-8 h-8 inline-block rounded-full ${taskColor}`}></span>
                <span className="text-black font-semibold text-xl">{task.expression}{task.status === "completed" ? " = " + task.answer : ""}</span>
                <span className={`text-gray`}>{prepareStatus(task)}</span>
                {task.status === "fail" ? (
                    <span className="text-[#c70808] font-semibold text-xl">
                        {task.answer}
                    </span>
                ) : (<span></span>)}
            </div>
            <div>
                <span className="text-completed"></span>
                <span className="text-gray text-sm">{formattedDate}</span>
            </div>
        </div>
    )
}

function prepareStatus(task) {
    switch (task.status) {
        case "created":
            return "Создан"
        case "processed":
            return `В процессе (Агент ${task.agent_id})`
        case "completed":
            let createdAt = new Date(task.created_at)
            let updatedAt = new Date(task.updated_at)
            let completeTime = (updatedAt.getTime() - createdAt.getTime()) / 1000
            return `Завершен за ${completeTime} секунд`
        case "fail":
            return "Ошибка"
        default:
            return "Unknown status"
    }
}