const elts = {
  text1: document.getElementById("text1"),
  text2: document.getElementById("text2"),
  textContainer: document.getElementById("container"),
  copyrightYear: document.getElementById("copyright-year"),
};

const texts = [
  { word: "Visualize", colour: "#945DEC", width: "388px" },
  { word: "Format", colour: "#FF8911", width: "309px" },
  { word: "Parse", colour: "#FC2947", width: "250px" },
];

const morphTime = 1;
const cooldownTime = 1.5;

let textIndex = texts.length - 1;
let time = new Date();
let morph = 0;
let cooldown = cooldownTime;

elts.copyrightYear.textContent = new Date().getFullYear();

elts.text1.textContent = texts[textIndex % texts.length].word;
elts.text2.textContent = texts[(textIndex + 1) % texts.length].word;
elts.text1.style.color = texts[textIndex % texts.length].colour;
elts.text2.style.color = texts[(textIndex + 1) % texts.length].colour;
elts.textContainer.style.width = texts[(textIndex + 1) % texts.length].width;

function doMorph() {
  morph -= cooldown;
  cooldown = 0;

  let fraction = morph / morphTime;

  if (fraction > 1) {
    cooldown = cooldownTime;
    fraction = 1;
  }

  setMorph(fraction);
}

function setMorph(fraction) {
  elts.text2.style.filter = `blur(${Math.min(8 / fraction - 8, 100)}px)`;
  elts.text2.style.opacity = `${Math.pow(fraction, 0.4) * 100}%`;

  fraction = 1 - fraction;
  elts.text1.style.filter = `blur(${Math.min(8 / fraction - 8, 100)}px)`;
  elts.text1.style.opacity = `${Math.pow(fraction, 0.4) * 100}%`;

  elts.text1.textContent = texts[textIndex % texts.length].word;
  elts.text2.textContent = texts[(textIndex + 1) % texts.length].word;
  elts.text1.style.color = texts[textIndex % texts.length].colour;
  elts.text2.style.color = texts[(textIndex + 1) % texts.length].colour;
  elts.textContainer.style.width = texts[(textIndex + 1) % texts.length].width;
}

function doCooldown() {
  morph = 0;

  elts.text2.style.filter = "";
  elts.text2.style.opacity = "100%";

  elts.text1.style.filter = "";
  elts.text1.style.opacity = "0%";
}

function animate() {
  requestAnimationFrame(animate);

  let newTime = new Date();
  let shouldIncrementIndex = cooldown > 0;
  let dt = (newTime - time) / 1000;
  time = newTime;

  cooldown -= dt;

  if (cooldown <= 0) {
    if (shouldIncrementIndex) {
      textIndex++;
    }

    doMorph();
  } else {
    doCooldown();
  }
}

animate();
