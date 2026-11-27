const playerWraper = document.querySelector(".videoPlayer");

const vid = document.querySelector(`.syncVideo`);
// const videoSoruce = document.querySelector(".syncVideo source");

const playPauseBtn = document.querySelector(".playPauseBtn");
const seekBar = document.querySelector(".seekBar");
const timeDisplay = document.querySelector(".timeDisplay");
const subtitlesToggle = document.querySelector(".subToggleBtn");

let isSubtitlesEnabled = true
const subOverlay = document.querySelector(".subOverlay");
const subText = document.querySelector(".subText");
const fullscreenBtn = document.querySelector(".fullscreenBtn");

// handle play pause 

function togglePlay() {
  if (vid.paused) {
    vid.play();
  } else {
    vid.pause();
  }
}
playPauseBtn.addEventListener("click",togglePlay)
vid.addEventListener("click", togglePlay);

vid.addEventListener("play", () => {
    playPauseBtn.textContent = "⏸ Pause";
});

vid.addEventListener("pause", () => {
    playPauseBtn.textContent = "▶ Play";
});

// handle seek  and procesbar
 function formatTime(sec) {
   if (isNaN(sec)) return "00:00";
   const m = Math.floor(sec / 60);
   const s = Math.floor(sec % 60);
   return `${String(m).padStart(2, "0")}:${String(s).padStart(2, "0")}`;
 }
// load video duration 
vid.addEventListener("loadedmetadata",()=>{
    seekBar.max = vid.duration
    timeDisplay.textContent = `00:00 / ${formatTime(video.duration)}`;
});

// update time  and proess bar when plying video 
vid.addEventListener("timeupdate",()=>{
    seekBar.value = vid.currentTime;
    timeDisplay.textContent = `${formatTime(video.currentTime)} / ${formatTime(
      video.duration
    )}`;
});
// seek handler
seekBar.addEventListener("input", () => {
  video.currentTime = seekBar.value;
});

vid.addEventListener("loadedmetadata",()=>{
    const track = vid.loadedmetadata;
    if(!track) return

    // hide htm render subtitle 
    track.mode = "hidden"
    // track subtitles ad set text in to sub overlat text 
    track.addEventListener("cuechange",()=>{
        const activeCue = track.activeCues[0];
        if(activeCue && isSubtitlesEnabled){
            subText.textContent = activeCue.text   
        }else{
            subText.textContent = ""
        }
    })
});

// hadle subtitle on off 
subtitlesToggle.addEventListener("click",()=>{
    isSubtitlesEnabled = !isSubtitlesEnabled

    subtitlesToggle.textContent = isSubtitlesEnabled ? "CC: ON" : "CC: OFF";
    subtitlesToggle.classList.toggle("active", isSubtitlesEnabled);

    if (!isSubtitlesEnabled) {
      subText.textContent = "";
    }

})

// handler full screen 
fullscreenBtn.addEventListener("click", () => {
  if (!playerWraper) {
    console.error("playerWraper element not found in HTML!");
    return;
  }

  if (!document.fullscreenElement) {
    playerWraper.requestFullscreen().catch((err) => {
      console.error("Fullscreen failed:", err);
    });
  } else {
    document.exitFullscreen();
  }
});