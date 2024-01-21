import React from "react";
import Link from "next/link";
import { RxSketchLogo, RxDashboard } from "react-icons/rx";
import { FaDocker } from "react-icons/fa";
import { DiMongodb } from "react-icons/di";
import { SiMinio } from "react-icons/si";
import { GiArtificialIntelligence } from "react-icons/gi";
import { IoIosStats } from "react-icons/io";

interface ISideBar {
    children: React.ReactNode
}

const Sidebar = ({children}: ISideBar) => {
    return <div className="flex ">
        <div className="fixed w-20 h-screen p-4 bg-white border-r-[1px] flex flex-col justify-between">
            <div className="flex flex-col items-center">
                <Link href="/dashboard">
                    <div className="bg-purple-800 text-white p-3 rounded-lg inline-block">
                        <RxSketchLogo size={20}/>
                    </div>
                </Link>
                <span className="border-b-[1px] border-gray-200 w-full p-1"></span>
                <Link href="/dashboard">
                    <div className="bg-gray-100 hover:bg-gray-200 transition-colors my-3 text-black p-3 rounded-lg inline-block">
                        <RxDashboard size={20}/>
                    </div>
                </Link>
                <Link href="/dashboard/statistic">
                    <div className="bg-gray-100 hover:bg-gray-200 transition-colors my-3 text-black p-3 rounded-lg inline-block">
                        <IoIosStats size={20}/>
                    </div>
                </Link>
                <Link href="/dashboard/models">
                    <div className="bg-gray-100 hover:bg-gray-200 transition-colors my-3 text-black p-3 rounded-lg inline-block">
                        <GiArtificialIntelligence size={20}/>
                    </div>
                </Link>
                <span className="border-b-[1px] border-gray-200 w-full p-1"></span>
                <Link href="/dashboard">
                    <div className="bg-gray-100 hover:bg-gray-200 transition-colors my-3 text-black p-3 rounded-lg inline-block">
                        <SiMinio size={20}/>
                    </div>
                </Link>
                <Link href="/mongo">
                    <div className="bg-gray-100 hover:bg-gray-200 transition-colors my-3 text-black p-3 rounded-lg inline-block">
                        <DiMongodb size={20}/>
                    </div>
                </Link>
                <Link href="/dashboard">
                    <div className="bg-gray-100 hover:bg-gray-200 transition-colors my-3 text-black p-3 rounded-lg inline-block">
                        <FaDocker size={20}/>
                    </div>
                </Link>
            </div>
        </div>
        <main className="ml-20 w-full">
            {children}
        </main>
    </div>
}

export default Sidebar;