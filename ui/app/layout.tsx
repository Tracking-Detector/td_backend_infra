import type { Metadata } from "next";
import "@/styles/globals.css"


export const metadata: Metadata = {
  title: "Tracking Detector",
  description: "Tracking Detector is a chrome extension for blocking webtracker with A.I.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
