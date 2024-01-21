import Link from "next/link"
import { FaCheckCircle } from "react-icons/fa";
import StatusBadge from "./StatusBadge";
const ServerStatus = () => {
    return <section id="serverStatus" className="px-4 pt-4">
        <div className="w-full flex gap-3 justify-start">
            <StatusBadge endpoint="/api/requests/health" serviceName="Requests"></StatusBadge>
            <StatusBadge endpoint="/api/users/health" serviceName="Users"></StatusBadge>
            <StatusBadge endpoint="/api/dispatch/health" serviceName="Dispatch"></StatusBadge>
        </div>
    </section>
}
export default ServerStatus