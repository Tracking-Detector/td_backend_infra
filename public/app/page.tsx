"use client"
import Navbar from "@/components/Navbar";
import Spline from "@splinetool/react-spline";
import Image from "next/image";
import Link from "next/link";

export default function Home() {
  return (
    <main className="flex min-h-screen h-fit flex-col items-center justify-center relative bg-black">
      <Navbar />
      <header id="home" className="flex flex-col-reverse md:flex-row w-full h-screen max-w-7xl items-center justify-center p-8 relative overflow-x-hidden">
        <div className="w-full h-2/4 md:h-full md:w-2/5 flex flex-col justify-center items-center md:items-start gap-8">
          <div className="flex flex-col gap-2">
            <h1 className="text-4xl text-white md:text-8xl">Tracking Detector</h1>
            <h2 className="text-md md:text-2xl text-white">Block web tracker with the help of A.I.</h2>
          </div>
          <p className="max-w-md text-sm md:text-base text-zinc-500">Tracking Detector provides web extensions that block web tracker with the help of A.I. Additionally, it has a backend which can handle model training and comparison.</p>
          <div className="w-full flex items-center justify-center md:justify-start gap-4">
            <button className="w-48 h-12 text-sm sm:text-base rounded bg-white text-black hover:bg-fuchsia-700 hover:text-white transition-colors">Install Now!</button>
            <button className="w-48 h-12 text-white text-sm sm:text-base rounded hover:bg-white hover:text-white hover:bg-opacity-5 transition-colors">Contact</button>
          </div>
        </div>
        <div className="w-fullh-2/4 md:h-full md:w-3/5 flex items-center justify-center relative">
          <Spline className="w-full flex scale-[.25] sm:scale-[.35] lg:scale-[.5] items-center justify-center md:justify-start" draggable={false} scene="https://prod.spline.design/BTQDRWxiHGxWxAGl/scene.splinecode" />
        </div>
      </header>
      <section id="about" className="h-fit min-h-screen w-full flex relative items-center justify-center p-8">
        <div className="w-full h-full flex items-start justify-start flex-col gap-8 max-w-7xl bg-opacity-5 bg-white rounded p-3">
          <h3 className="text-3xl md:text-5xl text-white">About</h3>
          <p className="text-1xl md:text-1xl text-white">Welcome to <span className="text-bold text-fuchsia-500">Tracking Detector</span>, an innovative open-source project dedicated to enhancing your online privacy by harnessing the power of machine learning. Our mission is to empower users with effective tools to thwart web trackers and take control of their digital footprint.</p>
          <p className="text-1xl md:text-1xl text-white">Tracking Detector is a project that manifests in the form of two distinct Chrome extensions, each designed to provide a tailored solution to your online privacy needs.</p>
          <ul className="text-start">
            <li className="text-1xl md:text-1xl text-white">
              <h4 className="my-3 font-bold">Tracking Detector Extension</h4>
              <p>The Tracking Detector extension leverages a statically trained machine learning model to identify and block web trackers seamlessly. By utilizing state-of-the-art technology, we ensure that your browsing experience remains secure and free from intrusive tracking mechanisms.</p>
            </li>
            <li className="text-1xl md:text-1xl text-white">
              <h4 className="my-3 font-bold">Train on Surf Extension</h4>
              <p>For those who crave a more personalized approach to privacy, our Train on Surf extension is the perfect fit. This extension allows users to actively train the machine learning model while surfing the internet. You have the autonomy to decide which model to employ for blocking trackers, giving you the flexibility to customize your privacy preferences.</p>
            </li>
            <li className="text-1xl md:text-1xl text-white">
              <h4 className="my-3 font-bold">Backend</h4>
              <p>Behind the scenes, Tracking Detector boasts a sophisticated backend infrastructure. Our backend serves a crucial role in comparing models and continuously training them on new data. This iterative process ensures that our machine learning models stay up-to-date and effective against evolving tracking techniques.</p>
            </li>
          </ul>
          <p className="text-1xl md:text-1xl text-white">At Tracking Detector, we believe in transparency, user empowerment, and staying ahead of the curve when it comes to online privacy. Join us in our mission to make the internet a safer space for everyone. Install Tracking Detector extensions today and take control of your digital footprint.</p>
        </div>
      </section>
      <section id="install" className="h-fit min-h-screen w-full flex relative items-center justify-center p-8">
        <div className="w-full h-full flex items-center justify-start gap-8 max-w-7xl ">
          <div className="bg-opacity-5 bg-white rounded p-3 w-2/4 flex flex-col justify-center align-middle items-center">
            <Image src="/td_icon.png" alt="" width={128} height={128}></Image>
            <h4 className="text-2xl md:text-1xl text-white mt-4">Tracking Detector Extension</h4>
            <ul className="my-3 flex flex-col items-center">
              <li className="text-1xl md:text-1xl text-white my-1">Utilizes statically trained model to block trackers.</li>
              <li className="text-1xl md:text-1xl text-white my-1">Download latest release from Github.</li>
              <li className="text-1xl md:text-1xl text-white my-1">Install extension through developer options.</li>
            </ul>
            <button className="w-28 h-10 text-sm sm:text-base rounded bg-fuchsia-700 text-white hover:bg-fuchsia-500 hover:text-white transition-colors mt-2">
              <Link href="https://github.com/Tracking-Detector/tracking_detector/releases">Install Now!</Link></button>
          </div>
          <div className="bg-opacity-5 bg-white rounded p-3 w-2/4 flex flex-col justify-center align-middle items-center">
            <Image src="/tos_icon.png" alt="" width={128} height={128}></Image>
            <h4 className="text-2xl md:text-1xl text-white mt-4">Train on Surf Extension</h4>
            <ul className="my-3 flex flex-col items-center">
              <li className="text-1xl md:text-1xl text-white my-1">Extension trains model while you surf.</li>
              <li className="text-1xl md:text-1xl text-white my-1">Download latest release from Github.</li>
              <li className="text-1xl md:text-1xl text-white my-1">Install extension through developer options.</li>
            </ul>
            <button className="w-28 h-10 text-sm sm:text-base rounded bg-fuchsia-700 text-white hover:bg-fuchsia-500 hover:text-white transition-colors mt-2">
              <Link href="https://github.com/Tracking-Detector/td_train_on_surf/releases">Install Now!</Link></button>
          </div>
        </div>
      </section>
      <section id="models" className="h-fit min-h-screen w-full flex relative items-center justify-center p-8">
        <div className="w-full h-full flex items-start justify-start flex-col gap-8 max-w-7xl bg-opacity-5 bg-white rounded p-3">
          <h3 className="text-3xl md:text-5xl text-white">Models</h3>
          <p className="text-1xl md:text-1xl text-white">Welcome to <span className="text-bold text-fuchsia-500">Tracking Detector</span>, an innovative open-source project dedicated to enhancing your online privacy by harnessing the power of machine learning. Our mission is to empower users with effective tools to thwart web trackers and take control of their digital footprint.</p>
          <p className="text-1xl md:text-1xl text-white">Tracking Detector is a project that manifests in the form of two distinct Chrome extensions, each designed to provide a tailored solution to your online privacy needs.</p>
          <ul className="text-start">
            <li className="text-1xl md:text-1xl text-white">
              <h4 className="my-3 font-bold">Tracking Detector Extension</h4>
              <p>The Tracking Detector extension leverages a statically trained machine learning model to identify and block web trackers seamlessly. By utilizing state-of-the-art technology, we ensure that your browsing experience remains secure and free from intrusive tracking mechanisms.</p>
            </li>
            <li className="text-1xl md:text-1xl text-white">
              <h4 className="my-3 font-bold">Train on Surf Extension</h4>
              <p>For those who crave a more personalized approach to privacy, our Train on Surf extension is the perfect fit. This extension allows users to actively train the machine learning model while surfing the internet. You have the autonomy to decide which model to employ for blocking trackers, giving you the flexibility to customize your privacy preferences.</p>
            </li>
            <li className="text-1xl md:text-1xl text-white">
              <h4 className="my-3 font-bold">Backend</h4>
              <p>Behind the scenes, Tracking Detector boasts a sophisticated backend infrastructure. Our backend serves a crucial role in comparing models and continuously training them on new data. This iterative process ensures that our machine learning models stay up-to-date and effective against evolving tracking techniques.</p>
            </li>
          </ul>
          <p className="text-1xl md:text-1xl text-white">At Tracking Detector, we believe in transparency, user empowerment, and staying ahead of the curve when it comes to online privacy. Join us in our mission to make the internet a safer space for everyone. Install Tracking Detector extensions today and take control of your digital footprint.</p>
        </div>
      </section>
      <section id="contact" className="h-fit min-h-screen w-full flex relative items-center justify-center p-8">
        <div className="w-full h-full flex items-center justify-center flex-col gap-8 max-w-7xl  rounded p-3">
          <h3 className="text-3xl md:text-5xl text-white">Contact</h3>
          <div className="md:w-3/5 md:max-w-full w-full mx-auto">
            <div className="sm:rounded-md p-6 border border-fuchsia-300">
              <form method="POST" action="/api/message">
                <label className="block mb-6">
                  <span className="text-white">Your name</span>
                  <input
                    type="text"
                    name="name"
                    className=" focus:border-fuchsia-300 focus:ring text-white h-10 px-2 bg-opacity-5 bg-white focus:ring-fuchsia-300 focus:ring-opacity-50 block w-full mt-1 rounded"
                  />
                </label>
                <label className="block mb-6">
                  <span className="text-white">Email address</span>
                  <input
                    name="email"
                    type="email"
                    className=" focus:border-fuchsia-300 bg-opacity-5 h-10 px-2 text-white bg-white focus:ring focus:ring-fuchsia-300 focus:ring-opacity-50 block w-full mt-1 rounded"
                    required
                  />
                </label>
                <label className="block mb-6">
                  <span className="text-white">Message</span>
                  <textarea
                    name="message"
                    className=" focus:border-fuchsia-300 bg-opacity-5 text-white px-2 bg-white focus:ring focus:ring-fuchsia-300 focus:ring-opacity-50 block w-full mt-1 rounded"
                    rows={6}
                    required
                  ></textarea>
                </label>
                <div className="mb-2 flex flex-row-reverse">
                  <button
                    type="submit"
                    className="focus:shadow-outline hover:bg-fuchsia-800 h-10 px-5 text-indigo-100 transition-colors duration-150 bg-fuchsia-700 rounded-lg"
                  >
                    Contact Us
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      </section>
    </main>
  );
}
