'use client';

import { useEffect, useState } from 'react';

export function ApiStatus() {
  const [status, setStatus] = useState<'checking' | 'online' | 'offline'>('checking');

  useEffect(() => {
    const checkApi = async () => {
      try {
        const response = await fetch('/healthz', { method: 'GET' });
        if (response.ok) {
          setStatus('online');
        } else {
          setStatus('offline');
        }
      } catch {
        setStatus('offline');
      }
    };

    checkApi();
    const interval = setInterval(checkApi, 30000); // Check every 30s

    return () => clearInterval(interval);
  }, []);

  if (status === 'checking') {
    return (
      <div className="fixed bottom-4 right-4 bg-surface border border-outline rounded-lg px-4 py-2 text-sm">
        <span className="text-on-surface-variant">Verificando API...</span>
      </div>
    );
  }

  if (status === 'offline') {
    return (
      <div className="fixed bottom-4 right-4 bg-error border border-error rounded-lg px-4 py-3 max-w-md shadow-lg">
        <div className="flex items-start gap-3">
          <span className="text-2xl">⚠️</span>
          <div>
            <div className="font-bold text-on-primary mb-1">API não está disponível</div>
            <div className="text-sm text-on-primary/90">
              A API Go não está respondendo em <code className="bg-black/20 px-1 rounded">localhost:8080</code>
            </div>
            <div className="text-xs text-on-primary/80 mt-2">
              Execute: <code className="bg-black/20 px-1 rounded">cd api && go run main.go</code>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="fixed bottom-4 right-4 bg-success border border-success rounded-lg px-4 py-2 text-sm flex items-center gap-2">
      <span className="w-2 h-2 bg-on-primary rounded-full animate-pulse"></span>
      <span className="text-on-primary font-medium">API Online</span>
    </div>
  );
}
