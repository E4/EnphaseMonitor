'use strict'

var Chart = function() {

  connectToServer()


  function receiveMessage(e) {
    var data = JSON.parse(e.data)
    addDataPoint(data);
  }

  var width = window.innerWidth;
  var height = window.innerHeight*0.75;
  var margin = { top: 20, right: 20, bottom: 30, left: 110 };
  var maxPoints = 900;

  var data = [[], [], []];

  var svg = d3.select(document.body).append("svg").attr("width", width).attr("height", height).style("font-family", "Roboto");

  var plotArea = svg.append("g")
    .attr("transform", `translate(${margin.left},${margin.top})`);

  var innerWidth = width - margin.left - margin.right;
  var innerHeight = height - margin.top - margin.bottom;

  var xScale = d3.scaleTime().range([0, innerWidth]);
  var yScale = d3.scaleLinear().range([innerHeight, 0]);

  var xAxis = plotArea.append("g").attr("transform", `translate(0,${innerHeight})`);
  var yAxis = plotArea.append("g");

  var lines = data.map(() => plotArea.append("path").attr("fill", "none").attr("stroke-width", 1.5));

  var colors = ["#23c35a", "orange", "red"];
  var labels = [
    "Production Phase Average",
    "Consumption Phase A",
    "Consumption Phase B",
  ];

  var movementDirection = 0;
  var movementPosition = 0;
  var movementLastFrameTime = performance.now();

  lines.forEach((line, i) => line.attr("stroke", colors[i]));


  var elements = {
    production:createDiv("production"),
    consumption:createDiv("consumption"),
    production_a:createDiv("production_a"),
    production_b:createDiv("production_b"),
    consumption_a:createDiv("consumption_a"),
    consumption_b:createDiv("consumption_b"),
    moving_bar:createDiv("moving_bar_producing"),
  }
  for(var i in elements) {
    document.body.appendChild(elements[i]);
  }

  updateChart();

  animate();

  this.addDataPoint = addDataPoint;

  function connectToServer() {
    var protocol = (location.protocol === 'https:') ? 'wss://' : 'ws://';
    var host = location.hostname;
    var port = location.port ? ':' + location.port : '';
    var ws = new WebSocket(protocol + host + port + '/stream');
    ws.onmessage = receiveMessage;
    ws.onclose = ()=>{window.setTimeout(()=>connectToServer(),5000)};
    ws.onerror = ()=>{window.setTimeout(()=>connectToServer(),5000)};
  }

  function addDataPoint(din) {
    elements.production_a.innerText = Math.round(din['production']['ph-a']['p']);
    elements.production_b.innerText = Math.round(din['production']['ph-b']['p']);
    elements.consumption_a.innerText = Math.round(din['total-consumption']['ph-a']['p']);
    elements.consumption_b.innerText = Math.round(din['total-consumption']['ph-b']['p']);
    elements.production.innerText = Math.round(din['production']['ph-a']['p'] + din['production']['ph-b']['p']);
    elements.consumption.innerText = Math.round(din['total-consumption']['ph-a']['p'] + din['total-consumption']['ph-b']['p']);
    const now = new Date();
    data[0].push({ time: now, value: din['production']['ph-a']['p']+din['production']['ph-b']['p'] });
    data[1].push({ time: now, value: din['total-consumption']['ph-a']['p'] });
    data[2].push({ time: now, value: din['total-consumption']['ph-b']['p'] });
    for(let i=0;i<3;i++) if(data[i].length>maxPoints) data[i].shift();
    let oldMovementDirection = movementDirection;
    movementDirection = (din['production']['ph-a']['p'] + din['production']['ph-b']['p']) - (din['total-consumption']['ph-a']['p'] + din['total-consumption']['ph-b']['p'])
    if(Math.sign(oldMovementDirection)!=Math.sign(movementDirection)) {
      if(Math.sign(movementDirection)==1) elements.moving_bar.className = "moving_bar_producing";
      else elements.moving_bar.className = "moving_bar_consuming";
    }
    updateChart();
  }

  function updateChart() {
    const now = new Date();
    const allData = data.flat();
    if (allData.length === 0) return;

    const firstTime = d3.min(allData, d => d.time);
    const xStart = firstTime;
    const xEnd = now;
    xScale.domain([xStart, xEnd]);

    const allValues = allData.map(d => d.value);
    const yMax = d3.max(allValues) ?? 1;
    yScale.domain([0, yMax]);

    const lineGen = d3.line().x(d => xScale(d.time)).y(d => yScale(d.value));

    plotArea.selectAll(".grid-y").remove();

    plotArea.selectAll(".grid-y")
      .data(yScale.ticks(5))
      .enter()
      .append("line")
      .attr("class", "grid-y")
      .attr("x1", 0)
      .attr("x2", innerWidth)
      .attr("y1", d => yScale(d))
      .attr("y2", d => yScale(d))
      .attr("stroke", "#444")
      .attr("stroke-width", 0.5);

    plotArea.selectAll(".grid-x").remove();

    plotArea.selectAll(".grid-x")
      .data(xScale.ticks(5))
      .enter()
      .append("line")
      .attr("class", "grid-x")
      .attr("y1", 0)
      .attr("y2", innerHeight)
      .attr("x1", d => xScale(d))
      .attr("x2", d => xScale(d))
      .attr("stroke", "#444")
      .attr("stroke-width", 0.5);

    lines.forEach((line, i) => {
      line.datum(data[i]).attr("d", lineGen);
    });

    xAxis.call(d3.axisBottom(xScale)).selectAll("text").attr("fill", "silver").attr("font-size", "20px");
    yAxis.call(d3.axisLeft(yScale).tickFormat(d => `${d} W`)).selectAll("text").attr("fill", "silver").attr("font-size", "20px");;
    xAxis.selectAll("path,line").attr("stroke", "silver");
    yAxis.selectAll("path,line").attr("stroke", "silver");
  }

  function createDiv(className) {
    var rv = document.createElement("div");
    rv.className = className;
    return rv;
  }

  function animate() {
    var currentFrameTime = performance.now();
    var timeElapsed = currentFrameTime - movementLastFrameTime;
    movementLastFrameTime = currentFrameTime;
    movementPosition += (movementDirection / 300000) * timeElapsed;
    if(movementPosition<-2) movementPosition = 100;
    if(movementPosition>102) movementPosition = -2;
    elements.moving_bar.style.backgroundPositionX = `${movementPosition}%`
    window.requestAnimationFrame(animate);
  }

};
