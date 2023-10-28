const _supabase = supabase.createClient(
  "https://kkfqicodrymzypjruroz.supabase.co",
  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYW5vbiIsImlhdCI6MTY0MTE2MDAyNSwiZXhwIjoxOTU2NzM2MDI1fQ.uSADYXf2KKJ7faHBaM152Gz2m0djlOErD5rb6_jiYYI",
);

document.addEventListener("DOMContentLoaded", getLastTrack);
setInterval(getLastTrack, 5000);

const selectGenre = document.querySelector("#select-genre");
document.addEventListener("DOMContentLoaded", initSelectGenre);
selectGenre.addEventListener("change", initSelectGenre);
selectGenre.addEventListener("change", getLastTrack);

async function getLastTrack() {
  let { data, error } = await _supabase
    .from(selectGenre.value)
    .select("*")
    .limit(1)
    .order("id", { ascending: false });

  if (error) {
    console.log("ERROR", error);
    throw error;
  }

  let currentTrack = data[0];
  document.getElementById("fip-image").src =
    currentTrack.raw.now.cardVisual.src;
  document.getElementById("fip-firstLine").innerText =
    currentTrack.raw.now.firstLine;
  document.getElementById("fip-secondLine").innerText =
    currentTrack.raw.now.secondLine;
  if (currentTrack.raw.now.song.year !== 0) {
    document.getElementById("fip-year").innerText =
      currentTrack.raw.now.song.year;
  } else {
    document.getElementById("fip-year").innerText = "";
  }
  if (currentTrack.raw.now.song.release.label !== "") {
    document.getElementById("fip-label").innerText =
      currentTrack.raw.now.song.release.label;
  } else {
    document.getElementById("fip-label").innerText = "";
  }
  if (currentTrack.raw.now.song.release.title !== "") {
    document.getElementById("fip-title").innerText =
      currentTrack.raw.now.song.release.title;
  } else {
    document.getElementById("fip-title").innerText = "";
  }
  const spotifyBlock = document.getElementById("div-spotify");
  if (currentTrack.spotify_id !== "") {
    console.log("track ID", currentTrack.spotify_id);
    let oldHTML = spotifyBlock.innerHTML;
    let newHTML = `<iframe src="https://open.spotify.com/embed?uri=spotify:track:${currentTrack.spotify_id}" width="100%" height="400px" frameborder="0"></iframe>`;
    if (oldHTML !== newHTML) {
      spotifyBlock.innerHTML = newHTML;
    }
  } else {
    spotifyBlock.innerHTML = "";
  }
}

function initSelectGenre(event) {
  const helperElement = document.querySelector("#helper-element");
  helperElement.innerHTML =
    event.target.querySelector("option:checked").innerText;
  const val = event.target.querySelector("option:checked").value;
  if (val === "fip") {
    document.getElementById("link").href =
      "https://www.radiofrance.fr/fip/titres-diffuses";
  } else {
    document.getElementById("link").href =
      "https://www.radiofrance.fr/fip/radio-" +
      val.replace("fip_", "").replace("_", "-");
  }
  resize(helperElement.offsetWidth + 30);
}

function resize(width) {
  document.documentElement.style.setProperty("--dynamic-size", `${width}px`);
}
