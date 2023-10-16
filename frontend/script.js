document.addEventListener("DOMContentLoaded", function () {
  const apiURL =
    "https://open.spotify.com/oembed?url=https%3A%2F%2Fopen.spotify.com%2Ftrack%2F3gShs30zG9OzLXHtKQaUR2";

  fetch(apiURL)
    .then((response) => response.json())
    .then((data) => {
      const apiResponseDiv = document.getElementById("api-response");
      apiResponseDiv.innerHTML = data.html;
    })
    .catch((error) => {
      console.error("Error:", error);
    });
});
