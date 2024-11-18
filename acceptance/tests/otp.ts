import * as http from "node:http";

let messages = new Map<string, any>();

export function startSink() {
  const hostname = "127.0.0.1";
  const port = 3030;

  const server = http.createServer((req, res) => {
    console.log("Sink received message: ");
    let body = "";
    req.on("data", (chunk) => {
      body += chunk;
    });

    req.on("end", () => {
      console.log(body);
      const data = JSON.parse(body);
      messages.set(data.contextInfo.recipientEmailAddress, data.args.code);
      res.statusCode = 200;
      res.setHeader("Content-Type", "text/plain");
      res.write("OK");
      res.end();
    });
  });

  server.listen(port, hostname, () => {
    console.log(`Sink running at http://${hostname}:${port}/`);
  });
  return server;
}
