"use client";
import Image from 'next/image';
import { useTheme } from 'next-themes';
import { useEffect, useState } from 'react';

export default function ThemedImage({ sources, alt, ...props }: any) {
  const { resolvedTheme } = useTheme();
  const [mounted, setMounted] = useState(false);
  
  useEffect(() => setMounted(true), []);
  
  if (!mounted) return <Image src={sources.light} alt={alt} {...props} />;
  
  const src = resolvedTheme === 'dark' ? sources.dark : sources.light;
  
  return <Image src={src} alt={alt} {...props} />;
}
