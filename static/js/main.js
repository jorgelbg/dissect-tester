const API_URL = window.location.href.split("?")[0] + "api/";

function testSamples() {
  let pattern = document.querySelector("#pattern").value;
  let samples = document.querySelector("#samples").innerText;
  let resultTextArea = document.querySelector("#results");

  let url = new URL(API_URL);

  let body = new URLSearchParams();
  body.set("tokenizer", pattern);
  body.set("str", samples.split('\n').filter(function(line) {return line.length > 0}).join('\n'));

  fetch(url, {
    method: "POST",
    cache: "no-cache",
    body: body,
  })
    .then((res) => {
      if (res.ok) {
        return res.json();
      }

      return res.text();
    })
    .then((payload) => {
      resultTextArea.replaceChildren();
      if (!Array.isArray(payload)) {
        payload = [payload.replace("\n", "")];
      }

      payload.forEach((s, pos) => {
        s = JSON.stringify(s, null, 2)
        let textarea = document.createElement("textarea");
        textarea.className = "bg-gray-200 appearance-none border-2 border-gray-200 rounded py-2 px-4 text-gray-700 leading-tight focus:outline-none focus:bg-white focus:border-purple-500 font-mono w-full hover:bg-purple-400 hover:border-purple-400 hover:text-white";
        textarea.id="right"+pos;
        textarea.setAttribute("z-index", "20");
        textarea.value = s;
        let rows = (s.match(/\n/g) || '').length + 1;
        textarea.rows = rows > 2 ? rows : 2;
        textarea.onmouseover = function(e) {
          // selectTextareaLine(document.querySelector("#samples"), pos)
          selectContentEditableLine(document.querySelector("#samples"), pos);
          connectDivs("left"+pos, "right"+pos, "rgba(167, 139, 250, 1)", 0.8);
        }
        textarea.onmouseout = function(e) {
          // document.getSelection().removeAllRanges()
          clearContentEditableLine(document.querySelector("#samples"));
          document.getElementById('svg-canvas').innerHTML = "";
        }
        
        resultTextArea.appendChild(textarea);
      })

      return;
    });
}

function selectContentEditableLine(el, pos) {
  lines = el.innerText.split('\n').filter(function(line) {return line.length > 0});
  // console.log(lines);
  let mark = document.createElement('mark');
  mark.id="left"+pos;
  mark.className="block w-full bg-purple-400 p-2 rounded text-white -ml-2";
  mark.style="z-index:20;";
  mark.innerText = lines[pos].trim();
  // console.log(mark.outerHTML);
  lines[pos] = mark.outerHTML;

  el.innerHTML = lines.join('\n');
}

function clearContentEditableLine(el) {
  // remove all mark elements
  let lines = el.innerHTML.replace(/<\/?[^>]+>/gi, '').trim();
  el.innerText = lines;
}

function createSVG() {
  var svg = document.getElementById("svg-canvas");
  if (null == svg) {
    svg = document.createElementNS("http://www.w3.org/2000/svg", 
                                   "svg");
    svg.setAttribute('id', 'svg-canvas');
    svg.setAttribute('style', 'position:absolute;top:0px;left:0px');
    svg.setAttribute('width', document.body.clientWidth);
    svg.setAttribute('height', document.body.clientHeight);
    svg.setAttribute('pointer-events', 'none');
    svg.setAttribute("z-index", "10");
    svg.setAttributeNS("http://www.w3.org/2000/xmlns/", 
                       "xmlns:xlink", 
                       "http://www.w3.org/1999/xlink");
    document.body.appendChild(svg);
  }
  return svg;
}

function drawCircle(x, y, radius, color) {
  var svg = createSVG();
    var shape = document.createElementNS("http://www.w3.org/2000/svg", "circle");
  shape.setAttributeNS(null, "cx", x);
  shape.setAttributeNS(null, "cy", y);
  shape.setAttributeNS(null, "r",  radius);
  shape.setAttributeNS(null, "fill", color);
  svg.appendChild(shape);
}

function findAbsolutePosition(htmlElement) {
  var x = htmlElement.offsetLeft;
  var y = htmlElement.offsetTop;
  for (var x=0, y=0, el=htmlElement; 
       el != null; 
       el = el.offsetParent) {
         x += el.offsetLeft;
         y += el.offsetTop;
  }

  return {
      "x": x,
      "y": y
  };
}

function connectDivs(leftId, rightId, color, tension) {
  var left = document.getElementById(leftId);
  var right = document.getElementById(rightId);
	
  // var leftPos = findAbsolutePosition(left);
  // var x1 = leftPos.x;
  // var y1 = leftPos.y;
  let leftOff = getOffset(left);
  var x1 = leftOff.left;
  var y1 = leftOff.top;
  x1 += left.offsetWidth;
  y1 += (left.offsetHeight / 2);

  let rightOff = getOffset(right);
  var x2 = rightOff.left;
  var y2 = rightOff.top;
  y2 += (right.offsetHeight / 2);

  var width=x2-x1;
  var height = y2-y1;

  drawCircle(x1, y1, 6, color);
  drawCircle(x2, y2, 6, color);
  drawCurvedLine(x1, y1, x2, y2, color, tension);
}

function drawCurvedLine(x1, y1, x2, y2, color, tension) {
  var svg = createSVG();
  var shape = document.createElementNS("http://www.w3.org/2000/svg", 
                                       "path");
  var delta = (x2-x1)*tension;
  var hx1=x1+delta;
  var hy1=y1;
  var hx2=x2-delta;
  var hy2=y2;
  var path = "M "  + x1 + " " + y1 + 
             " C " + hx1 + " " + hy1 
                   + " "  + hx2 + " " + hy2 
             + " " + x2 + " " + y2;
  shape.setAttributeNS(null, "d", path);
  shape.setAttributeNS(null, "fill", "none");
  shape.setAttributeNS(null, "stroke", color);
  shape.setAttributeNS(null, "stroke-width", 5);
  svg.appendChild(shape);
}

function selectTextareaLine(tarea,lineNum) {
  // lineNum--; // array starts at 0
  var lines = tarea.value.split("\n");

  // calculate start/end
  var startPos = 0, endPos = tarea.value.length;
  for(var x = 0; x < lines.length; x++) {
      if(x == lineNum) {
          break;
      }
      startPos += (lines[x].length+1);

  }

  var endPos = lines[lineNum].length+startPos;

  // do selection
  // Chrome / Firefox

  if(typeof(tarea.selectionStart) != "undefined") {
      tarea.focus();
      tarea.selectionStart = startPos;
      tarea.selectionEnd = endPos;
      return true;
  }

  // IE
   if (document.selection && document.selection.createRange) {
      tarea.focus();
      tarea.select();
      var range = document.selection.createRange();
      range.collapse(true);
      range.moveEnd("character", endPos);
      range.moveStart("character", startPos);
      range.select();
      return true;
  }

  return false;
}

const editor = document.querySelector('pre')
editor.addEventListener("paste", function(e) {
  e.preventDefault();
  const text = e.clipboardData.getData('text/plain');
  document.execCommand("insertHTML", false, text);
});

function getOffset(el) {
  const rect = el.getBoundingClientRect();
  return {
    left: rect.left + window.scrollX,
    top: rect.top + window.scrollY
  };
}

document.querySelector("#results").addEventListener('scroll', function(e) {
  document.getElementById('svg-canvas').innerHTML = "";
});