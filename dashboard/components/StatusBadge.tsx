"use client"
import Link from "next/link"
import { useEffect, useState } from "react";
import { FaCheckCircle } from "react-icons/fa";
import { MdOutlineError } from "react-icons/md";
interface IStatusBade {
    endpoint: string
    serviceName: string
}

const StatusBadge = ({endpoint, serviceName}: IStatusBade) => {
    const [error, setError] = useState(false);
    useEffect(() => {
        fetch(endpoint).then(response => {
            if (response.status != 200) {
                setError(true)
            }
        })
    }, [endpoint, serviceName])
    return <Link href={endpoint} className="bg-white hover:bg-purple-300 transition-colors rounded-lg w-40 h-full flex align-baseline gap-4 px-4 py-4">
     {error ? <MdOutlineError size={30} className="text-red-700"/> : <FaCheckCircle size={30} className="text-green-700"/>}
     {serviceName}
 </Link>
}

export default StatusBadge