const API_URL = window.location.href.split("?")[0] + "api/";

function testSamples() {
  let pattern = document.querySelector("#pattern").value;
  let samples = document.querySelector("#samples").innerText;
  let resultTextArea = document.querySelector("#results");

  let url = new URL(API_URL);

  let body = new URLSearchParams();
  body.set("tokenizer", pattern);
  body.set("str", samples);

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
        textarea.className = "bg-gray-200 appearance-none border-2 border-gray-200 rounded py-2 px-4 text-gray-700 leading-tight focus:outline-none focus:bg-white focus:border-purple-500 font-mono w-full hover:bg-purple-200";
        textarea.value = s;
        let rows = (s.match(/\n/g) || '').length + 1;
        textarea.rows = rows > 2 ? rows : 2;
        textarea.onmouseover = function(e) {
          // selectTextareaLine(document.querySelector("#samples"), pos)
          selectContentEditableLine(document.querySelector("#samples"), pos);
        }
        textarea.onmouseout = function(e) {
          // document.getSelection().removeAllRanges()
          clearContentEditableLine(document.querySelector("#samples"));
        }
        
        resultTextArea.appendChild(textarea);
      })

      return;
    });
}

function selectContentEditableLine(el, pos) {
  lines = el.innerText.split('\n');
  console.log(lines);
  let mark = document.createElement('mark');
  mark.className="block w-full bg-purple-200 p-2";
  mark.innerText = lines[pos].trim();
  console.log(mark.outerHTML);
  lines[pos] = mark.outerHTML;

  el.innerHTML = lines.join('\n');
}

function clearContentEditableLine(el) {
  // remove all mark elements
  let lines = el.innerHTML.replace(/<\/?[^>]+>/gi, '').trim();
  el.innerText = lines;
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