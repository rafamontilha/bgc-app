import type { Metadata } from "next";
import { Roboto } from "next/font/google";
import { ApiStatus } from "@/components/ui/ApiStatus";
import "./globals.css";

const roboto = Roboto({
  weight: ['400', '500', '700'],
  subsets: ["latin"],
  display: 'swap',
  variable: '--font-roboto',
});

export const metadata: Metadata = {
  title: "BGC - Brasil Global Connect",
  description: "Dashboard TAM / SAM / SOM - Sistema de analytics para dados de exportação brasileira",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="pt-BR">
      <body className={`${roboto.variable} antialiased`}>
        {children}
        <ApiStatus />
      </body>
    </html>
  );
}
