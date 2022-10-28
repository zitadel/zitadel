import { useEffect } from 'react';

export default function HideLink({ links }) {
  useEffect(() => {
    setTimeout(() => {
      links.forEach((link) => {
        const links = document.querySelectorAll(`a[href="${link}"]`);
        console.log(document.querySelectorAll(`a[href*="github.com"]`));
        console.log(`[href="${link}"]`, links);
      });
    }, 1000);
  }, []);

  return null;
}
