var fs = require('fs');
var path = require('path')
var https = require('https');

var defaultEnvironmentJsonURL = 'http://localhost:8080/ui/console/assets/environment.json'
var devEnvFile = path.join(__dirname, "src", "assets", "environment.json")
var url = process.env["ENVIRONMENT_JSON_URL"] || defaultEnvironmentJsonURL;

https.get(url, function (res) {
    var body = '';

    res.on('data', function (chunk) {
        body += chunk;
    });

    res.on('end', function () {
        fs.writeFileSync(devEnvFile, body);
        console.log("Developing against the following environment")
        console.log(JSON.stringify(JSON.parse(body), null, 4))
    });
}).on('error', function (e) {
    console.error("Got an error: ", e);
});
