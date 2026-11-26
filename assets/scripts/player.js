const playerWraper = document.getElementsByClassName("videoPlayerWraper");

const vid = document.querySelector(`.syncVideo`);
// const videoSoruce = document.querySelector(".syncVideo source");

const playPauseBtn = document.querySelector(".playPauseBtn");
const seekBar = document.querySelector(".seekBar");
const timeDisplay = document.querySelector(".timeDisplay");
const subtitlesToggle = document.querySelector(".subToggleBtn");

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

