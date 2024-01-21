import Link from "next/link"


const Navbar = () => {
    return <div className="w-full h-16 backdrop-filter backdrop-blur-xl bg-opacity-20 border-b flex items-center justify-center">
        <div className="max-w-7xl w-full flex items-center justify-between p-4">
            <h6 className="font-bold text-white">Tracking Detector</h6>
            <ul className="flex gap-8">
                <li><Link href="#home" className="hover:text-fuchsia-500 transition-colors text-xs text-white sm:text-base">Home</Link></li>
                <li><Link href="#about" className="hover:text-fuchsia-500 transition-colors text-xs text-white sm:text-base">About</Link></li>
                <li><Link href="#install" className="hover:text-fuchsia-500 transition-colors text-xs text-white sm:text-base">Install</Link></li>
                <li><Link href="#models" className="hover:text-fuchsia-500 transition-colors text-xs text-white sm:text-base">Models</Link></li>
                <li><Link href="#contact" className="hover:text-fuchsia-500 transition-colors text-xs text-white sm:text-base">Contact</Link></li>
            </ul>
        </div>
    </div>
}

export default Navbar