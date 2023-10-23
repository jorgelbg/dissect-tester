const API_URL = window.location.href.split("?")[0] + "api/";

document.querySelector("#samples").addEventListener("paste", function(e) {
  e.preventDefault(); // Prevent the default paste behavior

  // Get the pasted text as plain text
  let pastedText = (e.clipboardData || window.clipboardData).getData("text");

  // Insert the plain text into the contenteditable div
  const selection = window.getSelection();
  if (!selection.rangeCount) return;
  selection.deleteFromDocument();
  selection.getRangeAt(0).insertNode(document.createTextNode(pastedText));
  selection.collapseToEnd();
});

function testSamples() {
  let pattern = document.querySelector("#pattern").value;
  let samples = document.querySelector("#samples").innerText;
  let resultTextArea = document.querySelector("#results");

  let url = new URL(API_URL);

  let body = new URLSearchParams();
  body.set("tokenizer", pattern);
  body.set("str", samples.split('\n').filter(function(line) { return line.length > 0 }).join('\n'));

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
        textarea.className = "bg-gray-200 appearance-none border-2 border-gray-200 rounded py-2 px-4 text-gray-700 leading-tight focus:outline-none focus:bg-white focus:border-purple-500 font-mono w-full hover:bg-purple-400 hover:border-purple-400 hover:text-white focus:text-black";
        textarea.id = "right" + pos;
        textarea.setAttribute("z-index", "20");
        textarea.value = s;
        let rows = (s.match(/\n/g) || '').length + 1;
        textarea.rows = rows > 2 ? rows : 2;
        textarea.onmouseover = function(e) {
          document.querySelector("#samples").blur();
          document.querySelector("#samples").classList.remove('py-2')
          selectContentEditableLine(document.querySelector("#samples"), pos);
          let el = document.getElementById('left' + pos);
          el.scrollIntoView({
            block: 'nearest',
            inline: 'start'
          });

          textarea.scrollIntoView({
            block: 'nearest',
            inline: 'start'
          });

          const intersectionObserver = new IntersectionObserver((entries) => {
            let [entry] = entries;
            if (entry.isIntersecting) {
              setTimeout(function() {
                connectDivs("left" + pos, "right" + pos, "rgba(167, 139, 250, 1)", 1);
              }, 100)
            }
          });

          intersectionObserver.observe(el)
          // connectDivs("left"+pos, "right"+pos, "rgba(167, 139, 250, 1)", 0.8);
        }
        textarea.onmouseout = function(e) {
          document.querySelector("#samples").classList.add('py-2')
          document.getSelection().removeAllRanges()
          clearContentEditableLine(document.querySelector("#samples"));
          document.getElementById('svg-canvas').innerHTML = "";
        }

        resultTextArea.appendChild(textarea);
      })

      return;
    });
}

function selectContentEditableLine(el, pos) {
  lines = el.innerText.split('\n').filter(function(line) { return line.length > 0 });
  // console.log(lines);
  let mark = document.createElement('mark');
  mark.id = "left" + pos;
  mark.className = "block bg-purple-400 p-2 pl-4 rounded text-white py-4 -ml-4";
  mark.style = "z-index:20;";
  mark.innerText = lines[pos].trim();
  // console.log(mark.outerHTML);
  lines[pos] = mark.outerHTML;

  el.innerHTML = lines.join('\n');
}

function clearContentEditableLine(el) {
  // remove all mark elements
  let marks = el.querySelectorAll('mark');
  marks.forEach((mark) => {
    let text = mark.innerText;
    let textNode = document.createTextNode(text);

    console.log({ mark, textNode });

    mark.parentNode.replaceChild(textNode, mark);
  });
}

function createSVG() {
  let svg = document.getElementById("svg-canvas");
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

function drawPoint(x, y, radius, color) {
  let svg = createSVG();
  let shape = document.createElementNS("http://www.w3.org/2000/svg", "circle");
  shape.setAttributeNS(null, "cx", x);
  shape.setAttributeNS(null, "cy", y);
  shape.setAttributeNS(null, "r", radius);
  shape.setAttributeNS(null, "fill", color);
  svg.appendChild(shape);
}

function connectDivs(leftId, rightId, color, tension) {
  let left = document.getElementById(leftId);
  let right = document.getElementById(rightId);

  let leftOff = getOffset(left);
  let x1 = leftOff.left;
  let y1 = leftOff.top;
  x1 += left.offsetWidth;
  y1 += (left.offsetHeight / 2);

  let rightOff = getOffset(right);
  let x2 = rightOff.left;
  let y2 = rightOff.top;
  y2 += (right.offsetHeight / 2);

  drawPoint(x1, y1, 6, color);
  drawPoint(x2, y2, 6, color);
  drawCurvedLine(x1, y1, x2, y2, color, tension);
}

function drawCurvedLine(x1, y1, x2, y2, color, tension) {
  let svg = createSVG();
  let shape = document.createElementNS("http://www.w3.org/2000/svg",
    "path");
  let delta = (x2 - x1) * tension;
  let hx1 = x1 + delta;
  let hy1 = y1;
  let hx2 = x2 - delta;
  let hy2 = y2;
  let path = "M " + x1 + " " + y1 +
    " C " + hx1 + " " + hy1
    + " " + hx2 + " " + hy2
    + " " + x2 + " " + y2;
  shape.setAttributeNS(null, "d", path);
  shape.setAttributeNS(null, "fill", "none");
  shape.setAttributeNS(null, "stroke", color);
  shape.setAttributeNS(null, "stroke-width", 5);
  svg.appendChild(shape);
}

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
