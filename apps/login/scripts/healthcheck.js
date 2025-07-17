const url = process.argv[2];

if (!url) {
  console.error("❌ No URL provided as command line argument.");
  process.exit(1);
}

try {
  const res = await fetch(url);
  if (!res.ok) process.exit(1);
  process.exit(0);
} catch (e) {
  process.exit(1);
}
