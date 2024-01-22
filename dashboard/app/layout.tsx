import type { Metadata } from "next";
import Sidebar from "@/components/Sidebar";
import "@/styles/globals.css"


export const metadata: Metadata = {
  title: "Dashboard",
  description: "Admin panel for tracking detector backend.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>
        <Sidebar>
        {children}
        </Sidebar>
      </body>
    </html>
  );
}
