
import Link from "next/link";
import { BiMessage } from "react-icons/bi";
const Header = () => {
    return <div className="flex justify-between px-4 pt-4">
        <h2>Dashboard</h2>
        <Link href="/messages" className="relative inline-flex items-center p-3 text-sm font-medium text-center text-black bg-white  hover:bg-gray-200 transition-colors rounded-lg border-gray-200 border-solid">
            <BiMessage size={20}/>
            <span className="sr-only">Notifications</span>
            <div className="absolute inline-flex items-center justify-center w-6 h-6 text-xs font-bold text-white bg-purple-800 border-2 border-white rounded-full -top-2 -end-2 dark:border-gray-900">20</div>
        </Link>
    </div>
}

export default Header;