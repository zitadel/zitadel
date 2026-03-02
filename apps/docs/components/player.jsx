"use client";
import React from "react";

export function ResponsivePlayer({ url }) {
  const getYouTubeId = (url) => {
    if (!url) return null;
    const regExp = /^.*(youtu.be\/|v\/|u\/\w\/|embed\/|watch\?v=|&v=)([^#&?]*).*/;
    const match = url.match(regExp);
    return (match && match[2].length === 11) ? match[2] : null;
  };

  const videoId = getYouTubeId(url);

  if (videoId) {
    return (
      /* 1. Removed 'player-wrapper' class to avoid global.css conflict (double height).
         2. Reduced max-width to 'max-w-2xl' so it fits better in the text flow.
      */
      <div className="relative w-full max-w-2xl mx-auto h-0 pb-[56.25%] mb-8 rounded-lg overflow-hidden border border-zinc-200 dark:border-zinc-800 bg-black shadow-md">
        <iframe
          src={`https://www.youtube.com/embed/${videoId}`}
          title="YouTube video player"
          className="absolute top-0 left-0 w-full h-full"
          frameBorder="0"
          allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
          allowFullScreen
        />
      </div>
    );
  }

  return (
    <div className="p-4 mb-8 border border-red-200 bg-red-50 text-red-600 rounded-lg text-sm">
      Video source not supported: {url}
    </div>
  );
}