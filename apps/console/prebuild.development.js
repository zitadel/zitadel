var fs = require('fs');
var path = require('path')
var http = require('http');
var https = require('https');
var urlModule = require('url');

var defaultEnvironmentJsonURL = 'http://localhost:8080/ui/console/assets/environment.json'
var devEnvFile = path.join(__dirname, "src", "assets", "environment.json")
var url = process.env["ENVIRONMENT_JSON_URL"] || defaultEnvironmentJsonURL;

var protocol = urlModule.parse(url).protocol;
var getter = protocol === 'https:' ? https.get : http.get;

getter(url, function (res) {
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
