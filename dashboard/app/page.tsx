import Header from "@/components/Header";
import ServerStatus from "@/components/ServerStatus";
import Image from "next/image";

export default function Home() {
  return (
    <main className="bg-gray-200 min-h-screen">
        <Header/>
        <ServerStatus/>
    </main>
  );
}
