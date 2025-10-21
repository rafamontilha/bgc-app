import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Configuração de proxy para API Go
  // Em dev local: usa localhost:8080
  // Em K8s: usa variável de ambiente API_URL (bgc-api:8080)
  async rewrites() {
    // API_URL é definida em build time via ARG no Dockerfile
    // Em produção (K8s): http://bgc-api:8080
    // Em dev local: http://localhost:8080
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || process.env.API_URL || 'http://bgc-api:8080';

    return [
      {
        source: '/market/:path*',
        destination: `${apiUrl}/market/:path*`,
      },
      {
        source: '/routes/:path*',
        destination: `${apiUrl}/routes/:path*`,
      },
      {
        source: '/chapters/:path*',
        destination: `${apiUrl}/chapters/:path*`,
      },
      {
        source: '/healthz',
        destination: `${apiUrl}/healthz`,
      },
      {
        source: '/docs/:path*',
        destination: `${apiUrl}/docs/:path*`,
      },
      {
        source: '/swagger/:path*',
        destination: `${apiUrl}/swagger/:path*`,
      },
    ];
  },

  // Otimizações para SSG
  reactStrictMode: true,

  // Output standalone para container Docker otimizado
  output: 'standalone',

  // Configurações de imagem
  images: {
    unoptimized: true,
  },

  // Permitir CORS em dev
  ...(process.env.NODE_ENV === 'development' && {
    async headers() {
      return [
        {
          source: '/api/:path*',
          headers: [
            { key: 'Access-Control-Allow-Origin', value: '*' },
            { key: 'Access-Control-Allow-Methods', value: 'GET,POST,PUT,DELETE,OPTIONS' },
            { key: 'Access-Control-Allow-Headers', value: 'Content-Type' },
          ],
        },
      ];
    },
  }),
};

export default nextConfig;
