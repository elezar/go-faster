var http = require('http');

var blessed = require('blessed');
var contrib = require('blessed-contrib');
var screen = blessed.screen();
screen.key(['escape', 'q', 'C-c'], function(ch, key) {
    return process.exit(0);
});

http.get('http://go-faster.devfest.com:8080/series/dev_fest/data/conversions?field=time:unix&field=download_speed:mbit_s', function(res) {
  res.on('data', function (chunk) {
    var parsedData = JSON.parse(chunk);

    var line = contrib.line({
          style: {
              line: "yellow",
              text: "green",
              baseline: "black"
          },
          xLabelPadding: 1,
          xPadding: 1,
          label: 'download_speed'
      });
    var data = {
        x: [],
        y: []
    };


    for (var e of parsedData.data) {
      data.x.push(e[0]);
      data.y.push(e[1]);
    }
    screen.append(line); //must append before setting data
    line.setData([data]);

    screen.render()

  });
}).on('error', function(e) {
  console.log("Got error: " + e.message);
});


