function testSamples() {
  let pattern = document.querySelector("#pattern").value;
  let samples = document.querySelector("#samples").value;
  let resultTextArea = document.querySelector("#results");

  let url = new URL("http://localhost:8080/api/");
  let params = {
    tokenizer: pattern,
    str: samples
  };

  Object.keys(params).forEach(k => url.searchParams.append(k, params[k]));

  fetch(url)
    .then(res => {
      return res.json();
    })
    .then(payload => {
      console.log(payload);
      resultTextArea.value = JSON.stringify(payload);
    });
}
