import type { Metadata } from "next";
import { Roboto } from "next/font/google";
import { ClerkProvider } from "@clerk/nextjs";
import { ApiStatus } from "@/components/ui/ApiStatus";
import { ThemeProvider } from "@/components/providers/ThemeProvider";
import "./globals.css";

const roboto = Roboto({
  weight: ['400', '500', '700'],
  subsets: ["latin"],
  display: 'swap',
  variable: '--font-roboto',
});

// ClerkProvider no layout impede SSG em todas as páginas (exige key em build time).
// force-dynamic garante que todas as rotas sejam renderizadas sob demanda (SSR/edge),
// o que é o comportamento correto para um app autenticado.
export const dynamic = 'force-dynamic';

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
    <ClerkProvider dynamic>
      <html lang="pt-BR">
        <body className={`${roboto.variable} antialiased`}>
          <ThemeProvider>
            {children}
            <ApiStatus />
          </ThemeProvider>
        </body>
      </html>
    </ClerkProvider>
  );
}
