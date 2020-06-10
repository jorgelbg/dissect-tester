const API_URL = window.location.href.split("?")[0] + "api/";

function testSamples() {
  let pattern = document.querySelector("#pattern").value;
  let samples = document.querySelector("#samples").value;
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
      if (Array.isArray(payload)) {
        let str = payload.map((s) => JSON.stringify(s, null, 2));
        resultTextArea.value = str.join("\n");
        resultTextArea.focus();

        return;
      }

      // handle the error message
      resultTextArea.value = payload;
      resultTextArea.focus();
    });
}
