"use client";
import { useEffect, useState } from 'react';

export default function BrowserOnly({ children, fallback }: any) {
  const [mounted, setMounted] = useState(false);
  
  useEffect(() => setMounted(true), []);
  
  if (!mounted) return fallback || null;
  
  return children();
}
