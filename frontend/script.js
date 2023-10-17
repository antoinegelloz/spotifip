const _supabase = supabase.createClient(
  "https://kkfqicodrymzypjruroz.supabase.co",
  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYW5vbiIsImlhdCI6MTY0MTE2MDAyNSwiZXhwIjoxOTU2NzM2MDI1fQ.uSADYXf2KKJ7faHBaM152Gz2m0djlOErD5rb6_jiYYI",
);
document.addEventListener("DOMContentLoaded", function () {
  getLastSpotifyID();
});

setInterval(getLastSpotifyID, 5000);

async function getLastSpotifyID() {
  let { data, error } = await _supabase
    .from("fip_electro")
    .select("spotify_id")
    .neq("spotify_id", "")
    .limit(1)
    .order("id", { ascending: false });

  if (error) {
    console.log("ERROR", error);
    throw error;
  }

  console.log("track ID", data[0].spotify_id);
  const spotifyBlock = document.getElementById("spotify-block");
  oldHTML = spotifyBlock.innerHTML;
  newHTML = `<iframe src="https://open.spotify.com/embed?uri=spotify:track:${data[0].spotify_id}" width="100%" height="400px" frameborder="0"></iframe>`;
  if (oldHTML !== newHTML) {
    spotifyBlock.innerHTML = newHTML;
  }
}
